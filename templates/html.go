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
	"github.com/DataDrake/static-cling/file"
	log "github.com/DataDrake/waterlog"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// HTML is a standard html/template
type HTML struct {
	file.File
	tmpl *template.Template
}

// NewHTML creates a new HTML template from the file at the specified path
func NewHTML(path string) (h *HTML, err error) {
	h = &HTML{}
	h.Init(filepath.Split(path))
	log.Debugf("Creating template from %q\n", path)
	err = h.Update()
	return
}

// Execute generates an HTML document from the template and the provided data
func (h *HTML) Execute(out io.Writer, data interface{}) error {
	return h.tmpl.Execute(out, data)
}

// Update re-reads the template from disk if it has changed
func (h *HTML) Update() error {
	changed, err := h.Stat()
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	if err := h.Open(os.O_RDONLY); err != nil {
		return err
	}
	defer h.Close()
	raw, err := h.ReadString()
	if err != nil {
		return err
	}
	h.tmpl, err = template.New(h.Name).Funcs(template.FuncMap(functions())).Parse(raw)
	return err
}
