package runevm

import (
	"fmt"
	"os"
	"time"
)

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

func builtin_Wait(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("wait requires exactly 1 argument")
	}

	// Using type assertions to check if x and y are of type int
	ms, ok := args[0].(int)

	if !ok {
		return fmt.Errorf("argument must be of type int, got: %T", args[0])
	}

	time.Sleep(time.Duration(ms) * time.Millisecond)

	return nil
}

func builtin_Exit(args ...interface{}) interface{} {
	if len(args) != 0 {
		return fmt.Errorf("exit requires no arguments")
	}
	os.Exit(0)
	return nil
}

func builtin_TypeOf(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("typeof requires exactly 1 argument")
	}

	switch args[0].(type) {
	case int:
		return "int"
	case float64:
		return "float"
	case bool:
		return "bool"
	case string:
		return "string"
	default:
		return "unknown"
	}
}
