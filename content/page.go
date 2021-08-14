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

package content

import (
	"errors"
	"github.com/DataDrake/static-cling/config"
	"github.com/DataDrake/static-cling/file"
	// "github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

// Page represents a single page to be rendered to the build tree
type Page struct {
	Title    string           `yaml:"title"`
	Author   string           `yaml:"author"`
	Date     time.Time        `yaml:"date"`
	Category string           `yaml:"category"`
	Vars     config.Variables `yaml:"vars"`
	Content  string           `yaml:"-"`
	file     *file.File
	meta     *file.File
}

// NewPage creates a Page record from a known File
func NewPage(content *file.File) (p *Page, err error) {
	p = &Page{
		file: content,
	}
	p.meta = file.NewFile(content.Name+".yaml", content.Dir)
	err = p.Update()
	return
}

// OutputPath provides the name and subdirectory for this Page in the build directory
func (p *Page) OutputPath() (name, dir string, err error) {
	name = p.file.Name + ".html"
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	dir, err = filepath.Rel(wd, p.file.Dir)
	return
}

// IsNewer checks if either the metadata or content have been modified after a certain time
func (p *Page) IsNewer(other time.Time) bool {
	return p.file.Modified.After(other) || p.meta.Modified.After(other)
}

// Update re-reads the content and metadata for this PAge
func (p *Page) Update() (err error) {
	if err = p.updateContent(); err != nil {
		return err
	}
	if err = p.meta.Open(os.O_RDONLY); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
		return
	}
	defer p.meta.Close()
	dec := yaml.NewDecoder(p.meta)
	err = dec.Decode(p)
	return
}

// ErrUnsupportedContent indicates that a file in the content tree cannot be converted to HTML
var ErrUnsupportedContent = errors.New("file contains content which cannot be rendered to HTML")

func (p *Page) updateContent() (err error) {
	if err = p.file.Open(os.O_RDONLY); err != nil {
		return
	}
	defer p.file.Close()
	if p.Content, err = p.file.ReadString(); err != nil {
		return
	}
	switch p.file.Ext {
	case ".html":
		// already in HTML format
	case ".md":
		// out := blackfriday.Run([]byte(p.Content))
		// p.Content = string(out)
		fallthrough
	case ".haml", ".timber":
		fallthrough
	default:
		err = ErrUnsupportedContent
	}
	return
}
