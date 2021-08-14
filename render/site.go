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
	log "github.com/DataDrake/waterlog"
)

// Site is the content of the Site we are rendering
type Site struct {
	Config *config.Site
	layout templates.Template
	tmpls  *templates.Tree
}

// NewSite creates a Site from a configuration and template tree
func NewSite(conf *config.Site, tmpls *templates.Tree) (site *Site, err error) {
	layout, err := tmpls.Root.Get("layout")
	if err != nil {
		return
	}
	site = &Site{
		Config: conf,
		layout: layout,
		tmpls:  tmpls,
	}
	return
}

// Render each of the sections and the root pages of the site
func (s *Site) Render(src, dst *content.Tree, force bool) error {
	if err := s.sections(src, dst, force); err != nil {
		return err
	}
	return s.index(src.Root, dst.Root, force)
}

// sections updates each section of the site as needed
func (s *Site) sections(src *content.Tree, dst *content.Tree, force bool) error {
	log.Infoln("Checking for sections that no longer exist")
	for name := range dst.Root.Dirs {
		if _, ok := src.Root.Dirs[name]; !ok {
			if err := dst.Root.RemoveAll(name); err != nil {
				return nil
			}
		}
	}
	log.Goodln("DONE")
	log.Infoln("Updating sections")
	for name := range src.Root.Dirs {
		config, ok := s.Config.Sections[name]
		if !ok {
			log.Warnf("Missing config for section %q, skipping", name)
			continue
		}
		section, err := NewSection(s, name, config)
		if err != nil {
			return err
		}
		if err := section.Render(src.Root, dst.Root, force); err != nil {
			return err
		}
	}
	log.Goodln("DONE")
	return nil
}

func (s *Site) index(src, dst *content.Dir, force bool) error {
	return nil
}
