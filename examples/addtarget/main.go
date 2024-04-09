package main

import (
	"fmt"
	"os"

	"github.com/tiagomelo/go-makefile-gen/mfile"
)

func main() {
	const makeFilePath = "."
	targetName := "my-target"
	if err := mfile.AddTargetToMakefile(makeFilePath, targetName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
