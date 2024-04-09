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
	targetContent := `@ do it\n\t@ do that\n\t@ echo "ok"`
	if err := mfile.AddTargetWithContentAndDependenciesToMakefile(makeFilePath, targetName, targetContent, targetDependencies); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
