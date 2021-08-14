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
	"strings"
	"time"
)

// Category is all of the data necessary to render an Category index page
type Category struct {
	Site     *config.Site
	Section  *config.Section
	Category *config.Category
	Page     *content.Page
	Pages    content.Pages
	name     string
	layout   templates.Template
	template templates.Template
}

// NewCategory creates a Category
func NewCategory(section *Section, conf *config.Category) (category *Category, err error) {
	template, err := section.tmpls.Get(section.Config.Templates.Category)
	if err != nil {
		return
	}
	page := &content.Page{
		Title:    conf.Name,
		Date:     time.Now(),
		Category: conf.Name,
	}
	category = &Category{
		Site:     section.Site,
		Section:  section.Config,
		Category: conf,
		Page:     page,
		name:     strings.ToLower(conf.Name),
		layout:   section.layout,
		template: template,
	}
	return
}

// Render generates HTML for this Category, using the specified templates
// TODO: figure out how to not do this all the time
func (c *Category) Render(src, dst *content.Dir, force bool) (err error) {
	sub, ok := dst.Subs[c.name]
	if !ok {
		if _, err = dst.Mkdir(c.name); err != nil {
			return
		}
		if sub, err = content.NewDir(filepath.Join(dst.Path, c.name)); err != nil {
			return err
		}
		dst.Subs[c.name] = sub
	}
	c.setPages(src)
	out, err := c.applyTemplates()
	if err != nil {
		return err
	}
	return writeHTML(sub, "index.html", out)
}

// setPages recurses the source directory for any and all pages in this Category
func (c *Category) setPages(src *content.Dir) {
	for _, dir := range src.Subs {
		c.setPages(dir)
	}
	for _, page := range src.Pages {
		if page.Category == c.Category.Name {
			c.Pages = append(c.Pages, page)
		}
	}
}

// applyTemplates generates HTML for this Index, using the specified templates
func (c *Category) applyTemplates() (out string, err error) {
	var content strings.Builder
	if err = c.template.Execute(&content, c); err != nil {
		return
	}
	c.Page.Content = content.String()
	content.Reset()
	if err = c.layout.Execute(&content, c); err != nil {
		return
	}
	out = content.String()
	return
}
