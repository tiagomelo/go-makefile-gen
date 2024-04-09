package main

import (
	"fmt"
	"os"

	"github.com/tiagomelo/go-makefile-gen/mfile"
)

func main() {
	const makeFilePath = "."
	// Pass false if you don't want to overwrite the existing Makefile.
	if err := mfile.GenerateMakefile(makeFilePath, true); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
