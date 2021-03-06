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

// Tree contains all of the content
type Tree struct {
	Root *Dir
}

// NewTree creates a new directory for the specified path
func NewTree(path string) (t *Tree, err error) {
	root, err := NewDir(path)
	if err != nil {
		return
	}
	t = &Tree{
		Root: root,
	}
	return
}
