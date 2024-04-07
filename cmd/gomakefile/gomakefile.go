// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/tiagomelo/go-makefile-gen/mfile"
)

// GenerateCommand is used to generate a Makefile
type GenerateCommand struct {
	MakefilePath              string `short:"p" long:"path" description:"Path to the Makefile" default:"."`
	OverwriteExistingMakefile bool   `short:"o" long:"overwrite" description:"Overwrite existing Makefile"`
}

// Execute is the method invoked for the generate command
func (g *GenerateCommand) Execute(args []string) error {
	if err := mfile.GenerateMakefile(g.MakefilePath, g.OverwriteExistingMakefile); err != nil {
		return err
	}
	absPath, err := absPath(g.MakefilePath)
	if err != nil {
		return err
	}
	fmt.Printf("Makefile was generated successfully at %s\n", absPath)
	return nil
}

// AddTargetCommand is used to add a target to the Makefile
type AddTargetCommand struct {
	TargetName    string `short:"t" long:"target" description:"Name of the target" required:"true"`
	TargetContent string `short:"c" long:"targetContent" description:"Content of the target"`
	MakefilePath  string `short:"p" long:"path" description:"Path to the Makefile" default:"."`
}

// Execute is the method invoked for the addtarget command
func (a *AddTargetCommand) Execute(args []string) error {
	if a.TargetContent != "" {
		if err := mfile.AddTargetWithContentToMakefile(a.MakefilePath, a.TargetName, a.TargetContent); err != nil {
			return err
		}
		return nil
	}
	if err := mfile.AddTargetToMakefile(a.MakefilePath, a.TargetName); err != nil {
		return err
	}
	absPath, err := absPath(a.MakefilePath)
	if err != nil {
		return err
	}
	makeFilePath := fmt.Sprintf("%s/%s", absPath, "Makefile")
	fmt.Printf("Target %s was generated successfully added to %s\n", a.TargetName, makeFilePath)
	return nil
}

// Options holds the command-line options
type Options struct {
	Generate  GenerateCommand  `command:"generate" description:"Generate a basic Makefile"`
	AddTarget AddTargetCommand `command:"addtarget" description:"Add a target to the Makefile"`
}

// absPath converts a relative file path to an absolute path.
func absPath(path string) (string, error) {
	return filepath.Abs(path)
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Println(err)
			os.Exit(1)
		default:
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
