package runevm

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
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

func builtin_Millisecs(args ...interface{}) interface{} {
	if len(args) != 0 {
		return fmt.Errorf("millisecs requires no arguments")
	}

	// Get the current time and return the milliseconds since the Unix epoch
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func builtin_Exit(args ...interface{}) interface{} {
	if len(args) != 0 {
		return fmt.Errorf("exit requires no arguments")
	}
	os.Exit(0)
	return nil
}

// Read the contents of a file and return them as a string
func builtin_ReadFileStr(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("readfile requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	filename, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	// Read the contents of the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	return string(content)
}

// Write a string to a file
func builtin_WriteFileStr(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("writefile requires exactly 2 arguments")
	}

	// Using type assertions to check if the arguments are of type string
	filename, ok1 := args[0].(string)
	content, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("arguments must be of type string, got: %T and %T", args[0], args[1])
	}

	// Write the contents to the file
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// Split a string into a slice of substrings based on a specified delimiter
func builtin_StrSplit(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("strsplit requires exactly 2 arguments")
	}

	// Using type assertions to check if the arguments are of type string
	str, ok1 := args[0].(string)
	delimiter, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("arguments must be of type string, got: %T and %T", args[0], args[1])
	}

	// Split the string based on the delimiter
	parts := strings.Split(str, delimiter)
	return parts
}

// Trim whitespace from both ends of a string
func builtin_StrTrim(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("trim requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	str, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	// Trim whitespace from both ends of the string
	return strings.TrimSpace(str)
}

// Trim whitespace from the left side of a string
func builtin_TrimLeft(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("trimleft requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	str, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	// Trim whitespace from the left side of the string
	return strings.TrimLeftFunc(str, unicode.IsSpace)
}

// Trim whitespace from the right side of a string
func builtin_TrimRight(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("trimright requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	str, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	// Trim whitespace from the right side of the string
	return strings.TrimRightFunc(str, unicode.IsSpace)
}

// Check if a given string character is a digit
func builtin_IsDigit(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("isdigit requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	char, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	if len(char) != 1 {
		return fmt.Errorf("argument must be a single character, got a string of length %d", len(char))
	}

	// Check if the character is a digit
	return unicode.IsDigit(rune(char[0]))
}

// Check if a given string character is an alphabetical character
func builtin_IsAlpha(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("isalpha requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	char, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	if len(char) != 1 {
		return fmt.Errorf("argument must be a single character, got a string of length %d", len(char))
	}

	// Check if the character is an alphabetical character
	return unicode.IsLetter(rune(char[0]))
}

// Check if a given string character is a whitespace character
func builtin_IsWhite(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("iswhite requires exactly 1 argument")
	}

	// Using type assertions to check if the argument is of type string
	char, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("argument must be of type string, got: %T", args[0])
	}

	if len(char) != 1 {
		return fmt.Errorf("argument must be a single character, got a string of length %d", len(char))
	}

	// Check if the character is a whitespace character
	return unicode.IsSpace(rune(char[0]))
}

// Replace occurrences of a substring within a string with another substring
func builtin_Replace(args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("replace requires exactly 3 arguments")
	}

	// Using type assertions to check if the arguments are of type string
	str, ok1 := args[0].(string)
	old, ok2 := args[1].(string)
	new, ok3 := args[2].(string)
	if !ok1 || !ok2 || !ok3 {
		return fmt.Errorf("arguments must be of type string, got: %T, %T, and %T", args[0], args[1], args[2])
	}

	// Replace occurrences of the old substring with the new substring
	return strings.ReplaceAll(str, old, new)
}

// Check if a string contains a specified substring
func builtin_Contains(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("contains requires exactly 2 arguments")
	}

	// Using type assertions to check if the arguments are of type string
	str, ok1 := args[0].(string)
	substr, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("arguments must be of type string, got: %T and %T", args[0], args[1])
	}

	// Check if the string contains the substring
	return strings.Contains(str, substr)
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

func builtin_slice(args ...interface{}) interface{} {
	if len(args) != 3 {
		return fmt.Errorf("slice requires exactly 3 arguments")
	}

	start, ok := args[1].(int)
	if !ok {
		return fmt.Errorf("second argument must be a valid start index")
	}

	end, ok := args[2].(int)
	if !ok {
		return fmt.Errorf("third argument must be a valid end index")
	}

	switch arg := args[0].(type) {
	case []interface{}:
		if start < 0 || end > len(arg) || start > end {
			return fmt.Errorf("index out of bounds for array slice")
		}
		return arg[start:end]
	case map[string]interface{}:
		keys := make([]string, 0, len(arg))
		for k := range arg {
			keys = append(keys, k)
		}
		if start < 0 || end > len(keys) || start > end {
			return fmt.Errorf("index out of bounds for map slice")
		}
		slicedMap := make(map[string]interface{})
		for _, k := range keys[start:end] {
			slicedMap[k] = arg[k]
		}
		return slicedMap
	case string:
		if start < 0 || end > len(arg) || start > end {
			return fmt.Errorf("index out of bounds for string slice")
		}
		return arg[start:end]
	default:
		return fmt.Errorf("first argument must be an array, map, or string, got %T", args[0])
	}
}

func builtin_sliceLeft(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("sliceFirst requires exactly 2 arguments")
	}

	end, ok := args[1].(int)
	if !ok {
		return fmt.Errorf("second argument must be a valid end index")
	}

	switch arg := args[0].(type) {
	case []interface{}:
		if end > len(arg) || end < 0 {
			return fmt.Errorf("index out of bounds for array slice")
		}
		return arg[:end]
	case map[string]interface{}:
		keys := make([]string, 0, len(arg))
		for k := range arg {
			keys = append(keys, k)
		}
		if end > len(keys) || end < 0 {
			return fmt.Errorf("index out of bounds for map slice")
		}
		slicedMap := make(map[string]interface{})
		for _, k := range keys[:end] {
			slicedMap[k] = arg[k]
		}
		return slicedMap
	case string:
		if end > len(arg) || end < 0 {
			return fmt.Errorf("index out of bounds for string slice")
		}
		return arg[:end]
	default:
		return fmt.Errorf("first argument must be an array, map, or string, got %T", args[0])
	}
}

func builtin_sliceRight(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("sliceLast requires exactly 2 arguments")
	}

	start, ok := args[1].(int)
	if !ok {
		return fmt.Errorf("second argument must be a valid start index")
	}

	switch arg := args[0].(type) {
	case []interface{}:
		if start < 0 || start > len(arg) {
			return fmt.Errorf("index out of bounds for array slice")
		}
		return arg[start:]
	case map[string]interface{}:
		keys := make([]string, 0, len(arg))
		for k := range arg {
			keys = append(keys, k)
		}
		if start < 0 || start > len(keys) {
			return fmt.Errorf("index out of bounds for map slice")
		}
		slicedMap := make(map[string]interface{})
		for _, k := range keys[start:] {
			slicedMap[k] = arg[k]
		}
		return slicedMap
	case string:
		if start < 0 || start > len(arg) {
			return fmt.Errorf("index out of bounds for string slice")
		}
		return arg[start:]
	default:
		return fmt.Errorf("first argument must be an array, map, or string, got %T", args[0])
	}
}

func builtin_Len(args ...interface{}) interface{} {
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

func builtin_New(args ...interface{}) interface{} {
	if len(args) != 1 {
		return fmt.Errorf("new requires exactly 1 argument")
	}

	switch v := args[0].(type) {
	case []interface{}:
		return deepCopyArray(v)
	case map[string]interface{}:
		return deepCopyMap(v)
	default:
		return fmt.Errorf("new can only create copies of arrays or tables, got %T", args[0])
	}
}

// Assert that a condition is true, errors with given message if the assert fails
func builtin_Assert(args ...interface{}) interface{} {
	if len(args) != 2 {
		return fmt.Errorf("assert requires exactly 1 argument")
	}

	// Using type assertions to check if the first argument is a boolean
	condition, ok := args[0].(bool)
	if !ok {
		return fmt.Errorf("first argument must be of type bool, got: %T", args[0])
	}

	// Using type assertions to check if the first argument is a boolean
	msg, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("second argument must be of type string, got: %T", args[0])
	}

	// Check the condition and return an error if it is false
	if !condition {
		return fmt.Errorf("assertion failed: %s", msg)
	}

	return nil
}

// //////////////////////////////////////////////////////////////////////////////
// Helper Functions
// //////////////////////////////////////////////////////////////////////////////

// Helper function to format arrays for pretty printing
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

// Helper function to format maps for pretty printing
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

func deepCopyArray(arr []interface{}) []interface{} {
	newArr := make([]interface{}, len(arr))
	for i, v := range arr {
		newArr[i] = deepCopyValue(v)
	}
	return newArr
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})
	for k, v := range m {
		newMap[k] = deepCopyValue(v)
	}
	return newMap
}

func deepCopyValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		return deepCopyArray(v)
	case map[string]interface{}:
		return deepCopyMap(v)
	default:
		return v
	}
}
