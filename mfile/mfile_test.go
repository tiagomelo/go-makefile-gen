// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package mfile

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateMakefile(t *testing.T) {
	testCases := []struct {
		name          string
		mockClosure   func(m *mockFileSystem)
		overwrite     bool
		expectedError error
	}{
		{
			name: "happy path",
			mockClosure: func(m *mockFileSystem) {
			},
		},
		{
			name: "happy path, overwrite",
			mockClosure: func(m *mockFileSystem) {
			},
			overwrite: true,
		},
		{
			name: "happy path, is directory",
			mockClosure: func(m *mockFileSystem) {
				m.isDirOutput = true
			},
		},
		{
			name: "happy path, file does not exist",
			mockClosure: func(m *mockFileSystem) {
				m.readFileErr = os.ErrNotExist
				m.isNotExistOutput = true
			},
		},
		{
			name: "stat returned error",
			mockClosure: func(m *mockFileSystem) {
				m.statErr = errors.New("stat error")
			},
		},
		{
			name: "error when reading file",
			mockClosure: func(m *mockFileSystem) {
				m.readFileErr = errors.New("read error")
			},
			expectedError: errors.New("reading Makefile at some/path: read error"),
		},
		{
			name: "error when writing file",
			mockClosure: func(m *mockFileSystem) {
				m.writeFileErr = errors.New("write error")
			},
			expectedError: errors.New("writing MakeFile at some/path: write error"),
		},
	}
	for _, tc := range testCases {
		m := new(mockFileSystem)
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(m)
			fsProvider = m
			err := GenerateMakefile("some/path", tc.overwrite)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAddTargetToMakefile(t *testing.T) {
	testCases := []struct {
		name          string
		targetName    string
		mockClosure   func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor)
		expectedError error
	}{
		{
			name:        "happy path",
			targetName:  "test-target",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {},
		},
		{
			name:       "happy path, is directory",
			targetName: "test-target",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mfs.isDirOutput = true
			},
		},
		{
			name:          "target name has space",
			targetName:    "test target",
			mockClosure:   func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {},
			expectedError: errors.New("target name cannot contain space"),
		},
		{
			name:       "error when opening file",
			targetName: "test-target",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mfs.openErr = errors.New("open error")
			},
			expectedError: errors.New("opening path/to/Makefile: open error"),
		},
		{
			name:       "error when parsing template",
			targetName: "test-target",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mtp.err = errors.New("parse error")
			},
			expectedError: errors.New("parsing template: parse error"),
		},
		{
			name:       "error when executing template",
			targetName: "test-target",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mte.err = errors.New("execute error")
			},
			expectedError: errors.New("executing template: execute error"),
		},
	}
	for _, tc := range testCases {
		mfs := new(mockFileSystem)
		mtp := new(mockTemplateProcessor)
		mte := new(mockTemplateExecutor)
		mtp.te = mte
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfs, mtp, mte)
			fsProvider = mfs
			templateProcessorProvider = mtp
			err := AddTargetToMakefile("path/to/Makefile", tc.targetName)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

func TestAddTargetWithContentToMakefile(t *testing.T) {
	testCases := []struct {
		name          string
		targetName    string
		targetContent string
		mockClosure   func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor)
		expectedError error
	}{
		{
			name:          "happy path",
			targetName:    "test-target",
			targetContent: "@ do something",
			mockClosure:   func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {},
		},
		{
			name:          "happy path, is directory",
			targetName:    "test-target",
			targetContent: "@ do something",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mfs.isDirOutput = true
			},
		},
		{
			name:          "target name has space",
			targetName:    "test target",
			targetContent: "@ do something",
			mockClosure:   func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {},
			expectedError: errors.New("target name cannot contain space"),
		},
		{
			name:          "error when opening file",
			targetName:    "test-target",
			targetContent: "@ do something",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mfs.openErr = errors.New("open error")
			},
			expectedError: errors.New("opening path/to/Makefile: open error"),
		},
		{
			name:          "error when parsing template",
			targetName:    "test-target",
			targetContent: "@ do something",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mtp.err = errors.New("parse error")
			},
			expectedError: errors.New("parsing template: parse error"),
		},
		{
			name:          "error when executing template",
			targetName:    "test-target",
			targetContent: "@ do something",
			mockClosure: func(mfs *mockFileSystem, mtp *mockTemplateProcessor, mte *mockTemplateExecutor) {
				mte.err = errors.New("execute error")
			},
			expectedError: errors.New("executing template: execute error"),
		},
	}
	for _, tc := range testCases {
		mfs := new(mockFileSystem)
		mtp := new(mockTemplateProcessor)
		mte := new(mockTemplateExecutor)
		mtp.te = mte
		t.Run(tc.name, func(t *testing.T) {
			tc.mockClosure(mfs, mtp, mte)
			fsProvider = mfs
			templateProcessorProvider = mtp
			err := AddTargetWithContentToMakefile("path/to/Makefile", tc.targetName, tc.targetContent)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf("expected no error, got %v", err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf("expected error to be %v, got nil", tc.expectedError)
				}
			}
		})
	}
}

type mockFileSystem struct {
	openFile         *os.File
	fileInfo         os.FileInfo
	statErr          error
	file             []byte
	openErr          error
	readFileErr      error
	writeFileErr     error
	isNotExistOutput bool
	isDirOutput      bool
}

func (m *mockFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return m.openFile, m.openErr
}

func (m *mockFileSystem) Stat(name string) (os.FileInfo, error) {
	return m.fileInfo, m.statErr
}

func (m *mockFileSystem) ReadFile(name string) ([]byte, error) {
	return m.file, m.readFileErr
}

func (m *mockFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return m.writeFileErr
}

func (m *mockFileSystem) IsNotExist(err error) bool {
	return m.isNotExistOutput
}

func (m *mockFileSystem) IsDir(fi fs.FileInfo) bool {
	return m.isDirOutput
}

type mockTemplateExecutor struct {
	err error
}

func (m *mockTemplateExecutor) Execute(wr io.Writer, data interface{}) error {
	return m.err
}

type mockTemplateProcessor struct {
	te  templateExecutor
	err error
}

func (m *mockTemplateProcessor) Parse(name, text string) (templateExecutor, error) {
	return m.te, m.err
}
