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
	vm.Set("print", builtin_Print)
	vm.Set("println", builtin_Println)
	vm.Set("wait", builtin_Wait)
	vm.Set("exit", builtin_Exit)
	vm.Set("typeof", builtin_TypeOf)

	return vm
}

func (r *RuneVM) Run(source string, filepath string) {
	r.filepath = filepath
	r.source = source

	stream := NewInputStream(string(source), filepath)
	tokenStream := NewTokenStream(stream)
	parser := NewParser(tokenStream)
	ast := parser.parseToplevel()

	evaluate(ast, r.env)
}

func (r *RuneVM) Set(name string, value interface{}) {
	r.env.Def(name, value)
}

func (r *RuneVM) get(name string) interface{} {
	return r.env.Get(name, nil)
}

func (r *RuneVM) GetFun(name string) (func(...interface{}) interface{}, error) {
	fn, ok := r.get(name).(func(...interface{}) interface{})
	if !ok {
		return nil, fmt.Errorf("'%s' is not a function", name)
	}
	return fn, nil
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
