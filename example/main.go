package main

import (
	"fmt"
	"os"

	"github.com/RednibCoding/runevm"
)

func main() {

	args := []string{"example.exe", "test.rune"}
	// args := os.Args
	if len(args) < 2 {
		fmt.Println("USAGE: rune <sourcefile>")
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

	printerFunc, err := vm.GetFun("printer")
	if err != nil {
		fmt.Println(err)
		return
	}

	toPrint, err := vm.GetString("toPrint")
	if err != nil {
		fmt.Println(err)
	}

	printerFunc(toPrint)
}
