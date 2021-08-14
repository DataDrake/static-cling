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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// ErrAlreadyClosed is returned when a file is closed twice
	ErrAlreadyClosed = errors.New("file was already closed")
	// ErrFileOpen is returned when a file is opened twice
	ErrFileOpen = errors.New("file is already opened")
	// ErrNotOpen is returned when a file is accessed while closed
	ErrNotOpen = errors.New("fie has not been opened")
)

// File is a generic representation of a File on disk
type File struct {
	Name     string
	Ext      string
	Dir      string
	Mode     fs.FileMode
	Modified time.Time
	f        *os.File
}

// NewFile creates a new File from the specified path
func NewFile(dir, name string) (f *File) {
	f = &File{}
	f.Init(dir, name)
	return f
}

// Init handles the setup for a new File
func (f *File) Init(dir, name string) {
	pieces := strings.SplitN(name, ".", 2)
	filename := pieces[0]
	var ext string
	if len(pieces) > 1 {
		ext = "." + pieces[1]
	}
	f.Name = filename
	f.Ext = ext
	f.Dir = dir
}

// Create opens the file, encuring it is creaded on disk if missing
func (f *File) Create(mode os.FileMode) (err error) {
	return f.openFile(os.O_CREATE|os.O_RDWR, mode)
}

// Open this File, with the specified access flags
func (f *File) Open(flag int) error {
	return f.openFile(flag, f.Mode)
}

// openFile with the specified flags and mode
func (f *File) openFile(flag int, mode os.FileMode) (err error) {
	_, err = f.Stat()
	if err != nil {
		return
	}
	if f.f != nil {
		err = ErrFileOpen
		return
	}
	f.f, err = os.OpenFile(f.Path(), flag, mode)
	return
}

// Close the file if it is open
func (f *File) Close() (err error) {
	if f.f == nil {
		err = ErrAlreadyClosed
		return
	}
	if err = f.f.Close(); err == nil {
		return
	}
	f.f = nil
	_, err = f.Stat()
	return
}

// IsNewer checks if the file was modified after a certain time
func (f *File) IsNewer(old *File) bool {
	return f.Modified.After(old.Modified)
}

// Duplicate clones a source File to this File's destination
func (f *File) Duplicate(src *File) error {
	if err := src.Open(os.O_RDONLY); err != nil {
		return err
	}
	if err := f.Open(os.O_CREATE | os.O_RDWR | os.O_TRUNC); err != nil {
		return err
	}
	if _, err := io.Copy(f, src); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return src.Close()
}

// Read retrieves data from this file (satisfies io.Reader)
func (f *File) Read(p []byte) (n int, err error) {
	if f.f == nil {
		err = ErrNotOpen
		return
	}
	return f.f.Read(p)
}

// ReadString reads the entire contents of this file as a string
func (f *File) ReadString() (raw string, err error) {
	if f.f == nil {
		err = ErrNotOpen
		return
	}
	data, err := io.ReadAll(f.f)
	if err != nil {
		return
	}
	raw = string(data)
	return
}

// Write stores data into this file (satisfies io.Writer)
func (f *File) Write(p []byte) (n int, err error) {
	if f.f == nil {
		err = ErrNotOpen
		return
	}
	return f.f.Write(p)
}

// Path gets the full filepath of the underlying file
func (f *File) Path() string {
	return filepath.Join(f.Dir, f.Name+f.Ext)
}

// Stat updates the metadata for this File in memory
func (f *File) Stat() (changed bool, err error) {
	var info fs.FileInfo
	if f.f == nil {
		info, err = os.Stat(f.Path())
	} else {
		info, err = f.f.Stat()
	}
	if err != nil {
		return
	}
	if mode := info.Mode(); mode != f.Mode {
		f.Mode = mode
		changed = true
	}
	if modified := info.ModTime(); modified != f.Modified {
		f.Modified = modified
		changed = true
	}
	return
}
