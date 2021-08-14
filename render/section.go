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
	"github.com/DataDrake/static-cling/config"
	"github.com/DataDrake/static-cling/content"
	"github.com/DataDrake/static-cling/templates"
	"path/filepath"
)

// Section contains all of the data necessary to configure rendering for a section
type Section struct {
	Site   *config.Site
	Config *config.Section
	name   string
	layout templates.Template
	tmpls  *templates.Dir
}

// NewSection creates a new Section
func NewSection(site *Site, name string, conf *config.Section) (section *Section, err error) {
	tmpls, err := site.tmpls.Sub(name)
	if err != nil {
		return
	}
	section = &Section{
		Site:   site.Config,
		Config: conf,
		name:   name,
		layout: site.layout,
		tmpls:  tmpls,
	}
	return
}

// Render updates the contents of a destination tree from a source tree, for a given Section
func (s *Section) Render(src, dst *content.Dir, force bool) (err error) {
	srcDir := src.Subs[s.name]
	dstDir, ok := dst.Subs[s.name]
	if !ok {
		if _, err = dstDir.Mkdir(s.name); err != nil {
			return err
		}
		sub, err := content.NewDir(filepath.Join(dstDir.Path, s.name))
		if err != nil {
			return err
		}
		dstDir.Subs[s.name] = sub
	}
	for name := range dstDir.Dirs {
		if s.Config.HasCategory(name) {
			continue
		}
		if _, ok := srcDir.Dirs[name]; ok {
			continue
		}
		if err = dstDir.RemoveAll(name); err != nil {
			return
		}
	}
	if err = s.renderCategories(srcDir, dstDir, force); err != nil {
		return
	}
	dir, err := NewSectionDir(s)
	if err != nil {
		return err
	}
	return dir.Render(srcDir, dstDir, force)
}

// renderCategories iterates through each category and renders as needed
func (s *Section) renderCategories(src, dst *content.Dir, force bool) error {
	for _, config := range s.Config.Categories {
		category, err := NewCategory(s, config)
		if err != nil {
			return err
		}
		if err = category.Render(src, dst, force); err != nil {
			return err
		}
	}
	return nil
}
