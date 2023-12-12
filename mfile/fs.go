// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package mfile

import (
	"io/fs"
	"os"
)

// fileSystem interface abstracts the file system operations. This allows
// for easier testing by mocking file system interactions. It includes
// methods for opening, reading, writing files, and checking their status.
type fileSystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	IsNotExist(err error) bool
	IsDir(fi fs.FileInfo) bool
}

// osFileSystem struct implements the fileSystem interface using
// the standard library's os package. This is the real implementation
// that interacts with the actual file system.
type osFileSystem struct{}

func (osFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (osFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (osFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (osFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (osFileSystem) IsDir(fi fs.FileInfo) bool {
	return fi.IsDir()
}
