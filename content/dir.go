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
	"github.com/DataDrake/static-cling/file"
)

// Dir is a directory of the content Tree
type Dir struct {
	*file.Dir
	Subs  map[string]*Dir
	Pages Pages
}

// NewDir creates a new directory for the specified path
func NewDir(path string) (d *Dir, err error) {
	dir, err := file.NewDir(path)
	if err != nil {
		return
	}
	d = &Dir{
		Dir: dir,
	}
	return
}

// Update rereads the underlying directory and updates the pages as needed
func (d *Dir) Update(force bool) error {
	if err := d.Read(); err != nil {
		return err
	}
	// TODO update the Subs map
	// TODO update the Pages list
	return nil
}
