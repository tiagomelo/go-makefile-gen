package main

import (
	"fmt"
	"os"

	"github.com/tiagomelo/go-makefile-gen/mfile"
)

func main() {
	const makeFilePath = "."
	targetName := "my-target"
	targetDependencies := []string{"target-one", "target-two"}
	if err := mfile.AddTargetWithDependenciesToMakefile(makeFilePath, targetName, targetDependencies); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
