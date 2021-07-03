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

package util

import (
	log "github.com/DataDrake/waterlog"
	"io"
	"os"
	"path/filepath"
)

// CreateDir creates a new directory recursively
func CreateDir(path string) {
	log.Infof("Creating directory '%s'\n", path)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Fatalf("Failed to create directory '%s', reason: %s\n", path, err)
	}
	log.Goodln("Done.")
}

// ReparentPath converts a relative path from one directory to another
func ReparentPath(src, dst, path string) string {
	relPath, _ := filepath.Rel(src, path)
	return filepath.Join(dst, relPath)
}

// CopyFile copies the contents of a file from one location to another
func CopyFile(src, dst string) {
	log.Infof("Copying '%s' to '%s'\n", src, dst)
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatalf("Failed to open source file '%s', reason: %s\n", src, err)
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		log.Fatalf("Failed to open destination file '%s', reason: %s\n", dst, err)
	}
	defer srcFile.Close()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		log.Infof("Failed to copy from '%s' to '%s', reason: %s\n", src, dst, err)
	}
	log.Goodln("Done.")
}

// CopyDir copies the full contents of one directory to another
func CopyDir(srcDir, dstDir string) {
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == srcDir {
			return nil
		}
		dstPath := ReparentPath(srcDir, dstDir, path)
		if info.IsDir() {
			CreateDir(dstPath)
			return nil
		}
		CopyFile(path, dstPath)
		return nil
	})
}
