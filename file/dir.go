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
	"os"
	"path/filepath"
)

var (
	// ErrDirectoryCollision occurs when merging directories with the name subdir names
	ErrDirectoryCollision = errors.New("tried to merge directories with identically named subdirs")
	// ErrFileCollision occurs when merging directories with the name file names
	ErrFileCollision = errors.New("tried to merge firectories with identically named files")
)

// Dir represents a generic directory on Disk
type Dir struct {
	Path  string
	Dirs  map[string]*Dir
	Files map[string]*File
}

// NewDir creates a Dir from the specified path, recursively
func NewDir(path string) (d *Dir, err error) {
	d = &Dir{
		Path:  path,
		Dirs:  make(map[string]*Dir),
		Files: make(map[string]*File),
	}
	err = d.Read()
	return
}

// Merge combines two directories as if they shared the same root
func (d *Dir) Merge(other *Dir) (next *Dir, err error) {
	next = &Dir{
		Dirs:  make(map[string]*Dir),
		Files: make(map[string]*File),
	}
	if err = next.merge(d); err != nil {
		return
	}
	err = next.merge(other)
	return
}

// merge adds these
func (d *Dir) merge(other *Dir) error {
	for name, entry := range other.Dirs {
		if _, ok := d.Dirs[name]; !ok {
			return ErrDirectoryCollision
		}
		d.Dirs[name] = entry
	}
	for name, entry := range other.Files {
		if _, ok := d.Files[name]; !ok {
			return ErrFileCollision
		}
		d.Files[name] = entry
	}
	return nil
}

// Mkdir creates a new directory immediately inside of this directory
func (d *Dir) Mkdir(name string) (dir *Dir, err error) {
	path := filepath.Join(d.Path, name)
	if err = os.Mkdir(path, 0755); err != nil {
		return
	}
	if dir, err = NewDir(path); err != nil {
		return
	}
	d.Dirs[name] = dir
	return
}

// Read updates the contents of this directory from disk
func (d *Dir) Read() (err error) {
	entries, err := os.ReadDir(d.Path)
	if err != nil {
		return
	}
	if err = d.cleanDirs(entries); err != nil {
		return
	}
	if err = d.cleanFiles(entries); err != nil {
		return
	}
	err = d.update(entries)
	return
}

// cleanDirs removes directories that no longer exist
func (d *Dir) cleanDirs(entries []os.DirEntry) (err error) {
	for name := range d.Dirs {
		found := false
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if entry.Name() == name {
				found = true
				break
			}
		}
		if !found {
			delete(d.Dirs, name)
		}
	}
	return
}

// cleanFiles removes files that no logner exist
func (d *Dir) cleanFiles(entries []os.DirEntry) (err error) {
	for name := range d.Files {
		found := false
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if entry.Name() == name {
				found = true
				break
			}
		}
		if !found {
			delete(d.Files, name)
		}
	}
	return
}

// update all of the entries in this directory
func (d *Dir) update(entries []os.DirEntry) (err error) {
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			err = d.updateDir(d.Path, name)
		} else {
			err = d.updateFile(d.Path, name)
		}
		if err != nil {
			return
		}
	}
	return
}

// updateDir creates a missing directory or updates an existing one
func (d *Dir) updateDir(dir, name string) (err error) {
	next, ok := d.Dirs[name]
	if ok {
		err = next.Read()
	} else {
		next, err = NewDir(filepath.Join(dir, name))
		d.Dirs[name] = next
	}
	return
}

// updateFile creates a missing file or updates an existing one
func (d *Dir) updateFile(dir, name string) (err error) {
	next, ok := d.Files[name]
	if !ok {
		next = NewFile(dir, name)
		d.Files[name] = next
	}
	_, err = next.Stat()
	return
}

// RemoveAll recursively removes a subdirectory from disk
func (d *Dir) RemoveAll(name string) error {
	if err := os.RemoveAll(filepath.Join(d.Path, name)); err != nil {
		return err
	}
	delete(d.Dirs, name)
	return nil
}
