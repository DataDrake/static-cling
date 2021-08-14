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

package file

import (
	"errors"
)

// ErrMissingDir occurs when requesting a subtree for a non-existent directory
var ErrMissingDir = errors.New("subdir not found")

// Tree represents a generic file tree on disk
type Tree struct {
	Root *Dir
}

// NewTree creates a Tree from the specified path, recursively
func NewTree(path string) (t *Tree, err error) {
	t = &Tree{}
	t.Root, err = NewDir(path)
	return
}

// Read updates the contents of this Tree from disk
func (t *Tree) Read() (err error) {
	err = t.Root.Read()
	return
}

// Sub provides a subtree of an immediate directory
func (t *Tree) Sub(name string) (sub *Tree, err error) {
	d, ok := t.Root.Dirs[name]
	if !ok {
		err = ErrMissingDir
		return
	}
	sub = &Tree{
		Root: d,
	}
	return
}

// Merge combines two trees as if they had the same root
func (t *Tree) Merge(other *Tree) (next *Tree, err error) {
	root, err := t.Root.Merge(other.Root)
	if err != nil {
		return
	}
	next = &Tree{
		Root: root,
	}
	return
}
