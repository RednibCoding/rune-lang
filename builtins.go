package runevm

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Function to print elements
func builtin_Print(args ...interface{}) interface{} {
	for _, arg := range args {
		switch v := arg.(type) {
		case []interface{}:
			fmt.Print(formatArray(v))
		case map[string]interface{}:
			fmt.Print(formatMap(v))
		default:
			fmt.Print(v)
		}
	}
	return nil
}

// Function to print elements with a newline
func builtin_Println(args ...interface{}) interface{} {
	for _, arg := range args {
		switch v := arg.(type) {
		case []interface{}:
			fmt.Print(formatArray(v))
		case map[string]interface{}:
			fmt.Print(formatMap(v))
		default:
			fmt.Print(v)
		}
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
	case map[string]interface{}:
		return "map"
	default:
		return "unknown"
	}
}

func builtin_append(args ...interface{}) interface{} {
	if len(args) < 2 {
		return fmt.Errorf("append requires exactly 2 arguments for array/string or 3 arguments for map")
	}

	// First argument should be the array, string, or map
	switch arg := args[0].(type) {
	case []interface{}:
		return append(arg, args[1])
	case string:
		return arg + fmt.Sprint(args[1])
	case map[string]interface{}:
		if len(args) != 3 {
			return fmt.Errorf("append requires 3 arguments for map: map, key, value")
		}
		key, ok := args[1].(string)
		if !ok {
			return fmt.Errorf("second argument must be a string key for a map")
		}
		arg[key] = args[2]
		return arg
	default:
		return fmt.Errorf("first argument must be an array, string, or map, got %T", args[0])
	}
}

func builtin_remove(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("remove requires exactly 2 arguments")
	}

	// First argument should be the array, string, or map
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
			return fmt.Errorf("index %d out of bounds for string[%d]", index, len(arg))
		}
		return arg[:index] + arg[index+1:]
	case map[string]interface{}:
		key, ok := args[1].(string)
		if !ok {
			return fmt.Errorf("second argument must be a string key for a map")
		}
		if _, exists := arg[key]; !exists {
			return fmt.Errorf("key '%s' does not exist in map", key)
		}
		delete(arg, key)
		return arg
	default:
		return fmt.Errorf("first argument must be an array, string, or map, got %T", args[0])
	}
}

func builtin_hasKey(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("hasKey requires exactly 2 arguments")
	}

	// First argument should be the map
	argMap, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first argument must be a map, got %T", args[0])
	}

	// Second argument should be the key
	key, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("second argument must be a string key, got %T", args[1])
	}

	_, exists := argMap[key]
	return exists
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
	case map[string]interface{}:
		return len(arg)
	default:
		return fmt.Errorf("argument must be an array, string, or map, got %T", args[0])
	}
}

// //////////////////////////////////////////////////////////////////////////////
// Helper Functions
// //////////////////////////////////////////////////////////////////////////////
// Helper function to format arrays
func formatArray(arr []interface{}) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, elem := range arr {
		sb.WriteString(fmt.Sprintf("%v", elem))
		if i < len(arr)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// Helper function to format maps
func formatMap(m map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString("{")
	i := 0
	for key, value := range m {
		sb.WriteString(fmt.Sprintf("%v: %v", key, value))
		if i < len(m)-1 {
			sb.WriteString(", ")
		}
		i++
	}
	sb.WriteString("}")
	return sb.String()
}
