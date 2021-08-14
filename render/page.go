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
)

// Page contains all of the data necessary to render a single content Page
type Page struct {
	Site     *config.Site
	Section  *config.Section
	Page     *content.Page
	output   string
	layout   templates.Template
	template templates.Template
}

// NewPage creates a new Page
func NewPage(d *Dir, p *content.Page) (page *Page, err error) {
	name, _, err := p.OutputPath()
	if err != nil {
		return
	}
	page = &Page{
		Site:     d.Site,
		Section:  d.Section,
		Page:     p,
		output:   name,
		layout:   d.layout,
		template: d.content,
	}
	return
}

// Render generates the Page content as HTML, using the specified templates
func (p *Page) Render(src, dst *content.Dir, force bool) error {
	out, err := p.applyTemplates()
	if err != nil {
		return err
	}
	return writeHTML(dst, p.output, out)
}

// applyTemplates evaluates the page template and then uses the output as the content for the layout template
func (p *Page) applyTemplates() (out string, err error) {
	var content strings.Builder
	if err = p.template.Execute(&content, p); err != nil {
		return
	}
	p.Page.Content = content.String()
	content.Reset()
	if err = p.layout.Execute(&content, p); err != nil {
		return
	}
	out = content.String()
	return
}
