// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package mfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// For ease of unit testing.
var (
	// fsProvider is a variable of interface type fileSystem. It abstracts
	// file system operations and allows the use of different file system
	// implementations (like mocks for testing).
	fsProvider fileSystem = osFileSystem{}

	// templateProcessorProvider is a variable of interface type templateProcessor.
	// It abstracts template parsing and execution and allows different implementations.
	templateProcessorProvider templateProcessor = htmlTemplateProcessor{}
)

// Templates for the content to be added to the Makefile.
const (
	generateTemplate = `.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]\n"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: test
## test: run unit tests
test:
	@ go test -v ./... -count=1

.PHONY: coverage
## coverage: run unit tests and generate coverage report in html format
coverage:
	@ go test -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out
`
	addTargetTemplate = `.PHONY: {{ .TargetName }}
## {{ .TargetName }}: explain what {{ .TargetName }} does
{{ .TargetName }}:

`
	makefileName = "Makefile" // Default name for the Makefile.
)

// GenerateMakefile creates or updates a Makefile at the specified path.
// If `overwrite`, the existing Makefile will be overwritten.
func GenerateMakefile(path string, overwrite bool) error {
	makeFilePath := mkFilePath(path)
	content := generateTemplate
	if !overwrite {
		existingContent, err := fsProvider.ReadFile(makeFilePath)
		if err != nil && !fsProvider.IsNotExist(err) {
			return errors.Wrapf(err, "reading Makefile at %s", makeFilePath)
		}
		content = generateTemplate + string(existingContent)
	}
	if err := fsProvider.WriteFile(makeFilePath, []byte(content), 0644); err != nil {
		return errors.Wrapf(err, "writing MakeFile at %s", makeFilePath)
	}
	return nil
}

// AddTargetToMakefile appends a custom target to a Makefile.
// It ensures that the target name does not contain spaces and uses
// template processing to format the target addition.
func AddTargetToMakefile(path, targetName string) error {
	if containsSpace(targetName) {
		return errors.New("target name cannot contain space")
	}
	file, err := fsProvider.OpenFile(mkFilePath(path), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "opening %s", path)
	}
	defer file.Close()
	tmplExecutor, err := templateProcessorProvider.Parse("target", addTargetTemplate)
	if err != nil {
		return errors.Wrap(err, "parsing template")
	}
	err = tmplExecutor.Execute(file, map[string]string{"TargetName": targetName})
	if err != nil {
		return errors.Wrap(err, "executing template")
	}
	return nil
}

// mkFilePath calculates the full path to the Makefile.
// It checks if the provided path is a directory and appends the Makefile name to it.
func mkFilePath(path string) string {
	path = filepath.Clean(path)
	makeFilePath := path
	if fileInfo, err := fsProvider.Stat(path); err == nil && fsProvider.IsDir(fileInfo) {
		makeFilePath = filepath.Join(path, makefileName)
	}
	return makeFilePath
}

// containsSpace checks if the given string contains any spaces.
func containsSpace(s string) bool {
	return strings.Contains(s, " ")
}
