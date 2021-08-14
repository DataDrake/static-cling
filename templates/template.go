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
	"errors"
	"github.com/DataDrake/static-cling/file"
	"io"
	"path/filepath"
)

// Template represents any supported template type for content
type Template interface {
	// Execute generates an HTML document from the template
	Execute(out io.Writer, data interface{}) error
	// IsNewer checks if this template has been modified after a specific time
	IsNewer(other *file.File) bool
	// Update re-reads the template from disk
	Update() error
}

// ErrUnsupportedTemplate indicates that the file type of the specified template is not supported
var ErrUnsupportedTemplate = errors.New("template specified has unsupported extension")

// NewTemplate creates a new Template from a file on disk
func NewTemplate(path string) (Template, error) {
	switch filepath.Ext(path) {
	case ".html":
		return NewHTML(path)
	case ".haml":
		// return NewHAML(path)
		fallthrough
	default:
		return nil, ErrUnsupportedTemplate
	}
}
