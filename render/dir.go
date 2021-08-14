//
// Copyright 2021 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package render

import (
	"errors"
	"github.com/DataDrake/static-cling/config"
	"github.com/DataDrake/static-cling/content"
	"github.com/DataDrake/static-cling/templates"
)

var (
	// ErrNoListing is returned when there is no configured Listing template for the current directory
	ErrNoListing = errors.New("no listing template available at this depth of the tree")
)

// Dir contains all of the data necessary to configure rendering for a Directory
type Dir struct {
	Site     *config.Site
	Section  *config.Section
	Pages    content.Pages
	name     string
	listings []string
	layout   templates.Template
	content  templates.Template
	tmpls    *templates.Dir
	section  bool
}

// NewSectionDir creates a new Section
func NewSectionDir(section *Section) (dir *Dir, err error) {
	content, err := section.tmpls.Get(section.Config.Templates.Content)
	if err != nil {
		return
	}
	dir = &Dir{
		Site:     section.Site,
		Section:  section.Config,
		name:     section.name,
		listings: section.Config.Templates.Listings,
		layout:   section.layout,
		content:  content,
		tmpls:    section.tmpls,
		section:  true,
	}
	return
}

// Sub creates a subdirectory of this directory
func (d *Dir) Sub(name string) *Dir {
	var listings []string
	if len(d.listings) > 0 {
		listings = d.listings[1:]
	}
	return &Dir{
		Site:     d.Site,
		Section:  d.Section,
		name:     name,
		listings: listings,
		layout:   d.layout,
		content:  d.content,
		tmpls:    d.tmpls,
	}
}

// Render updates the contents of a destination directory from a source directory, for a given Dir config
func (d *Dir) Render(src, dst *content.Dir, force bool) error {
	for name := range dst.Dirs {
		if _, ok := src.Dirs[name]; ok {
			continue
		}
		if err := dst.RemoveAll(name); err != nil {
			return err
		}
	}
	d.Pages = src.Pages
	if err := d.renderIndex(src, dst, force); err != nil {
		return err
	}
	if err := d.renderPages(src, dst, force); err != nil {
		return err
	}
	return d.renderDirs(src, dst, force)
}

// renderIndex generates an index page if needed
func (d *Dir) renderIndex(src, dst *content.Dir, force bool) error {
	index, err := NewIndex(d)
	if err != nil {
		if err == ErrNoListing {
			return nil
		}
		return err
	}
	return index.Render(src, dst, force)
}

// renderPages generates a new page for each of the pages in this directory
func (d *Dir) renderPages(src, dst *content.Dir, force bool) error {
	for _, page := range src.Pages {
		p, err := NewPage(d, page)
		if err != nil {
			return err
		}
		if err = p.Render(src, dst, force); err != nil {
			return err
		}
	}
	return nil
}

// renderDirs update each subdirectory in this directory
func (d *Dir) renderDirs(src, dst *content.Dir, force bool) error {
	for name, dir := range src.Subs {
		sub := d.Sub(name)
		dstSub := dst.Subs[name]
		if err := sub.Render(dir, dstSub, force); err != nil {
			return err
		}
	}
	return nil
}
