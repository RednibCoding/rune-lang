package runevm

import (
	"fmt"
	"strconv"
)

type RuneVM struct {
	filepath string
	source   string
	env      *Environment
}

func NewRuneVM() *RuneVM {
	vm := &RuneVM{}

	vm.env = NewEnvironment(nil)
	vm.set("print", builtin_Print)
	vm.set("println", builtin_Println)
	vm.set("wait", builtin_Wait)
	vm.set("exit", builtin_Exit)
	vm.set("typeof", builtin_TypeOf)
	vm.set("append", builtin_append)
	vm.set("remove", builtin_remove)
	vm.set("haskey", builtin_hasKey)
	vm.set("slice", builtin_slice)
	vm.set("slicefirst", builtin_sliceFirst)
	vm.set("slicelast", builtin_sliceLast)
	vm.set("len", builtin_len)
	vm.set("new", builtin_New)

	return vm
}

func (r *RuneVM) Run(source string, filepath string) {
	r.filepath = filepath
	r.source = source

	stream := NewInputStream(string(source), filepath)
	tokenStream := NewTokenStream(stream)
	parser := NewParser(tokenStream)
	ast := parser.parseToplevel()
	evaluator := NewEvaluator()

	evaluator.evaluate(ast, r.env)
}

func (r *RuneVM) set(name string, value interface{}) {
	r.env.Def(name, value)
}

func (r *RuneVM) get(name string) interface{} {
	return r.env.Get(name, nil)
}

func (r *RuneVM) SetFun(name string, value func(...interface{}) interface{}) {
	r.set(name, value)
}

func (r *RuneVM) SetBool(name string, value bool) {
	r.set(name, value)
}

func (r *RuneVM) SetInt(name string, value int) {
	r.set(name, value)
}

func (r *RuneVM) SetFloat(name string, value float64) {
	r.set(name, value)
}

func (r *RuneVM) SetString(name string, value string) {
	r.set(name, value)
}

func (r *RuneVM) SetArray(name string, value []interface{}) {
	r.set(name, value)
}

func (r *RuneVM) GetBool(name string) (bool, error) {
	value := r.get(name)
	if b, ok := value.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("'%s' is not a bool", name)
}

func (r *RuneVM) GetString(name string) (string, error) {
	value := r.get(name)
	if s, ok := value.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("'%s' is not a string", name)
}

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

func (r *RuneVM) GetArray(name string) ([]interface{}, error) {
	val := r.get(name)
	if arr, ok := val.([]interface{}); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("variable '%s' is not an array", name)
}

func (r *RuneVM) GetTable(name string) (map[string]interface{}, error) {
	val := r.get(name)
	if arr, ok := val.(map[string]interface{}); ok {
		return arr, nil
	}
	return nil, fmt.Errorf("variable '%s' is not a table", name)
}

func (r *RuneVM) GetFun(name string) (func(...interface{}) interface{}, error) {
	fn, ok := r.get(name).(func(...interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a function", name)
	}
	return fn, nil
}

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
