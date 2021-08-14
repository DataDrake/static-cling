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
	"strings"
	"time"
)

// Index is all of the data necessary to render an index page
type Index struct {
	Site     *config.Site
	Section  *config.Section
	Page     *content.Page
	Pages    content.Pages
	layout   templates.Template
	template templates.Template
}

// NewIndex creates an Index
func NewIndex(d *Dir) (index *Index, err error) {
	if len(d.listings) == 0 {
		err = ErrNoListing
		return
	}
	tmpl, err := d.tmpls.Get(d.listings[0])
	if err != nil {
		return
	}
	page := &content.Page{
		Title: d.name,
		Date:  time.Now(),
	}
	index = &Index{
		Site:     d.Site,
		Section:  d.Section,
		Page:     page,
		Pages:    d.Pages,
		layout:   d.layout,
		template: tmpl,
	}
	return
}

// Render generates an index page using templates and page metadata
func (i *Index) Render(src, dst *content.Dir, force bool) error {
	out, err := i.applyTemplates()
	if err != nil {
		return err
	}
	return writeHTML(dst, "index.html", out)
}

// applyTemplates generates HTML for this Index, using the specified templates
func (i *Index) applyTemplates() (out string, err error) {
	var content strings.Builder
	if err = i.template.Execute(&content, i); err != nil {
		return
	}
	i.Page.Content = content.String()
	content.Reset()
	if err = i.layout.Execute(&content, i); err != nil {
		return
	}
	out = content.String()
	return
}
