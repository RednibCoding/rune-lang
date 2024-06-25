package runevm

import (
	"fmt"
	"strconv"
)

const Version = "v0.1.49"

type RuneVM struct {
	filepath string
	source   string
	env      *Environment
}

func NewRuneVM() *RuneVM {
	vm := &RuneVM{}

	vm.env = newEnvironment(nil)
	vm.set("version", builtin_VmVersion)
	vm.set("print", builtin_Print)
	vm.set("println", builtin_Println)
	vm.set("wait", builtin_Wait)
	vm.set("millis", builtin_Millisecs)
	vm.set("exit", builtin_Exit)
	vm.set("readfile", builtin_ReadFileStr)
	vm.set("writefile", builtin_WriteFileStr)
	vm.set("fileexist", builtin_FileExists)
	vm.set("direxists", builtin_DirExists)
	vm.set("isfileordir", builtin_IsFileOrDir)
	vm.set("strsplit", builtin_StrSplit)
	vm.set("strtrim", builtin_StrTrim)
	vm.set("trimleft", builtin_TrimLeft)
	vm.set("trimright", builtin_TrimRight)
	vm.set("isdigit", builtin_IsDigit)
	vm.set("isalpha", builtin_IsAlpha)
	vm.set("iswhite", builtin_IsWhite)
	vm.set("strreplace", builtin_Replace)
	vm.set("strcontains", builtin_Contains)
	vm.set("strhasprefix", builtin_HasPrefix)
	vm.set("strhassuffix", builtin_HasSuffix)
	vm.set("cutprefix", builtin_CutPrefix)
	vm.set("cutsuffix", builtin_CutSuffix)
	vm.set("strlower", builtin_StrToLower)
	vm.set("strupper", builtin_StrToUpper)
	vm.set("typeof", builtin_TypeOf)
	vm.set("append", builtin_append)
	vm.set("remove", builtin_remove)
	vm.set("haskey", builtin_hasKey)
	vm.set("slice", builtin_slice)
	vm.set("sliceleft", builtin_sliceLeft)
	vm.set("sliceright", builtin_sliceRight)
	vm.set("len", builtin_Len)
	vm.set("new", builtin_New)
	vm.set("exec", builtin_Exec)
	vm.set("assert", builtin_Assert)

	return vm
}

// Executes the Rune source code from the provided source string. Filepath is used for error reporting.
func (r *RuneVM) Run(source string, filepath string) {
	r.filepath = filepath
	r.source = source

	stream := newInputStream(string(source), filepath)
	tokenStream := newTokenStream(stream)
	parser := newParser(tokenStream)
	ast := parser.parseProgram()
	evaluator := newEvaluator()

	evaluator.evaluate(ast, r.env)
}

func (r *RuneVM) set(name string, value interface{}) {
	r.env.def(name, value)
}

func (r *RuneVM) get(name string) interface{} {
	return r.env.get(name, nil)
}

// Defines a function in the Rune environment.
func (r *RuneVM) SetFun(name string, value func(...interface{}) interface{}) {
	r.set(name, value)
}

// Defines a boolean variable in the Rune environment.
func (r *RuneVM) SetBool(name string, value bool) {
	r.set(name, value)
}

// Defines an integer variable in the Rune environment.
func (r *RuneVM) SetInt(name string, value int) {
	r.set(name, value)
}

// Defines a float variable in the Rune environment.
func (r *RuneVM) SetFloat(name string, value float64) {
	r.set(name, value)
}

// Defines a string variable in the Rune environment.
func (r *RuneVM) SetString(name string, value string) {
	r.set(name, value)
}

// Defines an array variable in the Rune environment.
func (r *RuneVM) SetArray(name string, value []interface{}) {
	r.set(name, value)
}

// Defines a table (map) variable in the Rune environment.
func (r *RuneVM) SetTable(name string, value map[string]interface{}) {
	r.set(name, value)
}

// Retrieves a boolean variable from the Rune environment.
func (r *RuneVM) GetBool(name string) (bool, error) {
	value := r.get(name)
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("'%s' is not a bool", name)
}

// Retrieves a string variable from the Rune environment.
func (r *RuneVM) GetString(name string) (string, error) {
	value := r.get(name)
	if s, ok := value.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("'%s' is not a string", name)
}

// GetInt retrieves an integer variable from the Rune environment.
func (r *RuneVM) GetInt(name string) (int, error) {
	value := r.get(name)
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i, nil
		}
	}
	return 0, fmt.Errorf("'%s' is not an int", name)
}

// Retrieves a float variable from the Rune environment.
func (r *RuneVM) GetFloat(name string) (float64, error) {
	value := r.get(name)
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, nil
		}
	}
	return 0, fmt.Errorf("'%s' is not a float", name)
}

// Retrieves an array variable from the Rune environment.
func (r *RuneVM) GetArray(name string) ([]interface{}, error) {
	val := r.get(name)
	if arr, ok := val.([]interface{}); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("variable '%s' is not an array", name)
}

// Retrieves a table (map) variable from the Rune environment.
func (r *RuneVM) GetTable(name string) (map[string]interface{}, error) {
	val := r.get(name)
	if arr, ok := val.(map[string]interface{}); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("variable '%s' is not a table", name)
}

// Retrieves a function from the Rune environment.
func (r *RuneVM) GetFun(name string) (func(...interface{}) interface{}, error) {
	fn, ok := r.get(name).(func(...interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a function", name)
	}
	return fn, nil
}

// Retrieves a function from a table (map) in the Rune environment.
func (r *RuneVM) GetTableFun(tableName string, funName string) (map[string]interface{}, func(...interface{}) interface{}, error) {
	table, err := r.GetTable(tableName)
	if err != nil {
		return nil, nil, fmt.Errorf(tableName, " is not a rune table")
	}

	fun, ok := table[funName].(func(args ...interface{}) interface{})
	if !ok {
		return nil, nil, fmt.Errorf(funName, " is not a function on table ", tableName)
	}
	return table, fun, nil
}
