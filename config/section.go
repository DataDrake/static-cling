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
	"strings"
	"time"
)

// Section for rendering content
type Section struct {
	Name       string      `yaml:"name"`
	Templates  Templates   `yaml:"templates"`
	Categories []*Category `yaml:"categories"`
	Vars       Variables   `yaml:"vars"`
}

// NewSection creates an empty Section configuration
func NewSection() *Section {
	return &Section{
		Vars: NewVariables(),
	}
}

// Sections is a map of Section configurations
type Sections map[string]*Section

func loadSections(dir string) (confs Sections, modified time.Time, err error) {
	log.Debugf("Loading configuration files from %q\n", dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	confs = make(Sections)
	var conf *Section
	var latest time.Time
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)
		if name == SiteFile {
			log.Debugf("Skipping site configuration file %q\n", path)
			continue
		}
		if entry.IsDir() {
			log.Debugf("Skipping directory %q\n", path)
			continue
		}
		if conf, latest, err = loadSection(path); err != nil {
			return
		}
		if latest.After(modified) {
			modified = latest
		}
		name = strings.TrimSuffix(name, filepath.Ext(name))
		confs[name] = conf
	}
	return
}

func loadSection(path string) (conf *Section, modified time.Time, err error) {
	log.Debugf("Loading section configuration %q\n", path)
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return
	}
	modified = info.ModTime()
	conf = NewSection()
	dec := yaml.NewDecoder(f)
	err = dec.Decode(conf)
	return
}

// HasCategory determines if a Category exists in this Section
func (s *Section) HasCategory(name string) bool {
	for _, category := range s.Categories {
		if strings.ToLower(category.Name) == name {
			return true
		}
	}
	return false
}
