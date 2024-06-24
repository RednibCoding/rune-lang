package main

import (
	"fmt"
	"os"

	"github.com/RednibCoding/runevm"
)

/*********************************************************
*
* To build the rune binary: use the build.bat on windows or build.sh on linux/macos
*
**********************************************************/

const version = "v0.1.40"

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Rune interpreter %s\n", version)
		fmt.Println("  USAGE: rune <sourcefile>")
		os.Exit(1)
	}
	source, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Printf("ERROR: Can't find source file '%s'.\n", args[1])
		os.Exit(1)
	}

	filepath := args[1]
	vm := runevm.NewRuneVM()
	vm.Run(string(source), filepath)
}
