package runevm

import "fmt"

func builtin_Print(args ...interface{}) interface{} {
	for _, arg := range args {
		fmt.Print(arg)
	}
	return nil
}

func builtin_Println(args ...interface{}) interface{} {
	for _, arg := range args {
		fmt.Print(arg)
	}
	fmt.Print("\n")
	return nil
}
