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
	case []interface{}:
		return "array"
	default:
		return "unknown"
	}
}

func builtin_append(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("append requires exactly 2 arguments")
	}

	// First argument should be the array or string
	switch arg := args[0].(type) {
	case []interface{}:
		return append(arg, args[1])
	case string:
		return arg + fmt.Sprint(args[1])
	default:
		return fmt.Errorf("first argument must be an array or a string, got %T", args[0])
	}
}

func builtin_remove(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("remove requires exactly 2 arguments")
	}

	// First argument should be the array or string
	switch arg := args[0].(type) {
	case []interface{}:
		index, ok := args[1].(int)
		if !ok {
			return fmt.Errorf("second argument must be a valid index")
		}
		if index < 0 || index >= len(arg) {
			return fmt.Errorf("index %d out of bounds for array[%d]", index, len(arg))
		}
		return append(arg[:index], arg[index+1:]...)
	case string:
		index, ok := args[1].(int)
		if !ok {
			return fmt.Errorf("second argument must be a valid index")
		}
		if index < 0 || index >= len(arg) {
			return fmt.Errorf("index %d out of bounds for array[%d]", index, len(arg))
		}
		return arg[:index] + arg[index+1:]
	default:
		return fmt.Errorf("first argument must be an array or a string, got %T", args[0])
	}
}

func builtin_len(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("len requires exactly 1 argument")
	}

	switch arg := args[0].(type) {
	case []interface{}:
		return len(arg)
	case string:
		return len(arg)
	default:
		return fmt.Errorf("argument must be an array or a string, got %T", args[0])
	}
}
