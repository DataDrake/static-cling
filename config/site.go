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

package config

import (
	log "github.com/DataDrake/waterlog"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// SiteFile is the relative path of the Site configuration
var SiteFile = "_site.yaml"

// Dir is the relative directory for configuration files
const Dir = "config"

// Path of the configuration directory
func Path(src string) string {
	return filepath.Join(src, Dir)
}

// Site is the full configuration for the site
type Site struct {
	Name       string    `yaml:"name"`
	Deployment string    `yaml:"deploy"`
	Vars       Variables `yaml:"vars"`
	Sections   Sections  `yaml:"-"`
	Dir        string    `yaml:"-"`
	modified   time.Time
}

// Load parses all of the config directories
func Load(dir string) (conf Site, err error) {
	if conf, err = loadSite(dir); err != nil {
		log.Errorf("Failed to load site config, reason: %q\n", err)
		return
	}
	var modified time.Time
	if conf.Sections, modified, err = loadSections(dir); err != nil {
		log.Errorf("Failed to load section configs, reason: %q\n", err)
	}
	if modified.After(conf.modified) {
		conf.modified = modified
	}
	return
}

// NewSite creates an empty Site configuration
func NewSite() Site {
	return Site{
		Vars: NewVariables(),
	}
}

func loadSite(dir string) (conf Site, err error) {
	file := filepath.Join(dir, SiteFile)
	log.Debugf("Loading site configuration file: %q\n", file)
	f, err := os.Open(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
		log.Warnf("Site configuration %q not found\n", file)
		err = nil
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return
	}
	conf = NewSite()
	conf.modified = info.ModTime()
	conf.Dir = dir
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&conf)
	return
}

// Equal checks if there are any differences between this Site and another
func (s *Site) Equal(other *Site) bool {
	return reflect.DeepEqual(s, other)
}

// IsNewer checks if this Site was modified after a certain time
func (s *Site) IsNewer(other time.Time) bool {
	return s.modified.After(other)
}

// Update rereads the configuration from disk and checks if it has changed, updating if it has
func (s *Site) Update() (changed bool, err error) {
	next, err := Load(s.Dir)
	if err != nil {
		return
	}
	if changed = !s.Equal(&next); changed {
		*s = next
	}
	return
}
