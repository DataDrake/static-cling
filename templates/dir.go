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

package templates

import (
	"fmt"
	"github.com/DataDrake/static-cling/file"
	log "github.com/DataDrake/waterlog"
)

// Dir is a directory containing Template files
type Dir struct {
	*file.Dir
	Templates map[string]Template
}

// NewDir creates a Dir from the specified path, recursively
func NewDir(path string) (d *Dir, err error) {
	dir, err := file.NewDir(path)
	if err != nil {
		return
	}
	d = &Dir{
		Dir:       dir,
		Templates: make(map[string]Template),
	}
	err = d.Update(true)
	return
}

// Sub returns a subdirectory of the current directory
func (d *Dir) Sub(name string) (next *Dir, err error) {
	dir, ok := d.Dirs[name]
	if !ok {
		err = fmt.Errorf("failed to find subdir %s", name)
		return
	}
	next = &Dir{
		Dir:       dir,
		Templates: make(map[string]Template),
	}
	err = next.Update(true)
	return
}

// Get retrieves a specific template by name
func (d *Dir) Get(name string) (tmpl Template, err error) {
	tmpl, ok := d.Templates[name]
	if !ok {
		err = fmt.Errorf("failed to find template %q", name)
	}
	return
}

// Update rereads the entire Dir from disk
func (d *Dir) Update(force bool) (err error) {
	if force {
		if err = d.Read(); err != nil {
			return
		}
	}
	// remove deleted templates
	for name := range d.Templates {
		found := false
		for _, file := range d.Files {
			if file.Name == name {
				found = true
				break
			}
		}
		if !found {
			delete(d.Templates, name)
		}
	}
	for _, file := range d.Files {
		name := file.Name
		next, ok := d.Templates[name]
		if ok {
			err = next.Update()
		} else {
			println(file.Path())
			next, err = NewTemplate(file.Path())
		}
		if err != nil {
			if err != ErrUnsupportedTemplate {
				return
			}
			log.Warnf("Failed to update template %q, reason: %s\n", name, err)
			continue
		}
		d.Templates[name] = next
	}
	return
}
