package main

import (
	"fmt"
	"os"

	"github.com/tiagomelo/go-makefile-gen/mfile"
)

func main() {
	const makeFilePath = "."
	targetName := "my-target"
	targetContent := `@ do it\n\t@ do that\n\t@ echo "ok"`
	if err := mfile.AddTargetWithContentToMakefile(makeFilePath, targetName, targetContent); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
