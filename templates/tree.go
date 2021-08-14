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
	log "github.com/DataDrake/waterlog"
	"path/filepath"
)

// Tree contains the full tree of Templates for this site
type Tree struct {
	Root *Dir
}

// Load reads the template Tree from disk
func Load(path string) (t *Tree, err error) {
	log.Debugln("Loading initial template tree")
	t = &Tree{}
	t.Root, err = NewDir(filepath.Join(path, "templates"))
	return
}

// Sub retrieves a subdirectory by name
func (t *Tree) Sub(path string) (dir *Dir, err error) {
	return t.Root.Sub(path)
}

// Update rereads the entire Tree
func (t *Tree) Update(force bool) error {
	return t.Root.Update(force)
}
