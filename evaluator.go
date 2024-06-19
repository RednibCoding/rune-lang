package runevm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func evaluate(exp *Expr, env *Environment) interface{} {
	if exp == nil {
		Error(exp, "Null expression error, this is a bug and should never happen!. Please file a bug!")
		return nil
	}
	switch exp.Type {
	case Num:
		return parseNumber(exp.Value.(string), exp)

	case Str, Bool:
		return exp.Value

	case Var:
		value := env.Get(exp.Value.(string), exp)
		if exp.Index != nil {
			switch v := value.(type) {
			case []interface{}:
				index := evaluate(exp.Index, env).(int)
				if index < 0 || index >= len(v) {
					Error(exp, "Index '%d' out of bounds for array '%v[%d]'", index, exp.Value, len(v))
				}
				return v[index]
			case map[string]interface{}:
				key := evaluate(exp.Index, env).(string)
				val, ok := v[key]
				if !ok {
					Error(exp, "Key '%s' not found in table '%v'", key, exp.Value)
				}
				return val
			default:
				Error(exp, "Variable %v is not an array or table", exp.Value)
			}
		}
		return value

	case Assign:
		if exp.Left.Type == Var && exp.Left.Index != nil {
			arrayOrTable := evaluate(exp.Left.Left, env)
			if arrayOrTable == nil {
				Error(exp, "Cannot assign to an index on a null expression")
				return nil
			}
			index := evaluate(exp.Left.Index, env)
			value := evaluate(exp.Right, env)

			switch arr := arrayOrTable.(type) {
			case []interface{}:
				idx, ok := index.(int)
				if !ok {
					Error(exp, "Array index must be an integer")
					return nil
				}
				if idx < 0 || idx >= len(arr) {
					Error(exp, "Array index out of bounds")
					return nil
				}
				arr[idx] = value
				return value
			case map[string]interface{}:
				key, ok := index.(string)
				if !ok {
					Error(exp, "Table key must be a string")
					return nil
				}
				arr[key] = value
				return value
			default:
				Error(exp, "Cannot index into type %T", arrayOrTable)
				return nil
			}
		} else if exp.Left.Type == Var && exp.Left.Value != nil {
			if exp.Left.Left != nil {
				// Handle field assignment
				obj := evaluate(exp.Left.Left, env)
				field := exp.Left.Value.(string)
				value := evaluate(exp.Right, env)
				if m, ok := obj.(map[string]interface{}); ok {
					m[field] = value
					return value
				} else {
					Error(exp, "Cannot assign to field %v on non-table object", field)
					return nil
				}
			}
			// Handle variable assignment
			return env.Set(exp.Left.Value.(string), evaluate(exp.Right, env), exp)
		}

		if exp.Left.Type != Var {
			Error(exp, "Cannot assign to %v", exp.Left)
		}
		return env.Set(exp.Left.Value.(string), evaluate(exp.Right, env), exp)

	case Binary:
		return applyOp(exp.Operator,
			evaluate(exp.Left, env),
			evaluate(exp.Right, env), exp)

	case Fun:
		return makeFun(env, exp)

	case If:
		cond := evaluate(exp.Cond, env)
		if cond != false {
			return evaluate(exp.Then, env)
		}
		if exp.Else != nil {
			if exp.Else.Type == Prog {
				for _, elif := range exp.Else.Prog {
					if elif.Type == If {
						elifCond := evaluate(elif.Cond, env)
						if elifCond != false {
							return evaluate(elif.Then, env)
						}
					} else {
						return evaluate(elif, env)
					}
				}
			} else {
				return evaluate(exp.Else, env)
			}
		}
		return false

	case While:
		for {
			cond := evaluate(exp.Cond, env)
			if !cond.(bool) {
				break
			}
			evaluate(exp.Body, env)
		}
		return false

	case Array:
		var arr []interface{}
		for _, element := range exp.Prog {
			arr = append(arr, evaluate(element, env))
		}
		return arr

	case Table:
		m := make(map[string]interface{})
		for _, pair := range exp.Prog {
			key := evaluate(pair.Left, env)
			value := evaluate(pair.Right, env)
			m[fmt.Sprint(key)] = value
		}
		return m

	case Prog:
		var val interface{} = false
		for _, e := range exp.Prog {
			val = evaluate(e, env)
		}
		return val

	case Call:
		fn, ok := evaluate(exp.Func, env).(func(args ...interface{}) interface{})
		if !ok {
			Error(exp, "'%s' is not a function", exp.Func.Value)
		}

		var args []interface{}

		// Check if caller is a go-map/rune-table ...
		callerName := exp.Func.Value
		caller, ok := env.Get(callerName.(string), exp).(map[string]interface{})
		// if so:
		if ok {
			// inject its reference as the first argument (similar to pythons 'self' argument on methods)
			args = append(args, caller)
		}

		for _, arg := range exp.Args {
			args = append(args, evaluate(arg, env))
		}

		ret := fn(args...)
		if err, ok := ret.(error); ok {
			Error(exp, "Error in function call: '%v'", err)
		}
		return ret

	case Import:
		filePath := evaluate(exp.Left, env).(string) + ".rune"
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			Error(exp, "Failed to import file '%s': %v", filePath, err)
		}
		stream := NewInputStream(string(fileContent), filePath)
		tokenStream := NewTokenStream(stream)
		parser := NewParser(tokenStream)
		ast := parser.parseToplevel()
		evaluate(ast, env)
		return nil

	default:
		Error(exp, "I don't know how to evaluate %v", exp.Type)
		return nil
	}
}

func parseNumber(val string, exp *Expr) interface{} {
	if strings.Contains(val, ".") {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			Error(exp, "Expected number but got %T", val)
		}
		return f
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		Error(exp, "Expected number but got %T", val)
	}
	return i
}

func applyOp(op string, a, b interface{}, exp *Expr) interface{} {
	num := func(x interface{}) float64 {
		switch v := x.(type) {
		case string:
			return parseNumber(v, exp).(float64)
		case int:
			return float64(v)
		case float64:
			return v
		default:
			Error(exp, "Expected number but got %T", x)
			return 0
		}
	}
	div := func(x interface{}) float64 {
		if num(x) == 0 {
			Error(exp, "Divide by zero")
		}
		return num(x)
	}
	boolVal := func(x interface{}) bool {
		switch v := x.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		case int:
			return v != 0
		case float64:
			return v != 0
		default:
			Error(exp, "Expected bool but got %v", x)
		}
		return false
	}
	switch op {
	case "+":
		return num(a) + num(b)
	case "-":
		return num(a) - num(b)
	case "*":
		return num(a) * num(b)
	case "/":
		return num(a) / div(b)
	case "%":
		return int(num(a)) % int(num(b))
	case "&&":
		return boolVal(a) && boolVal(b)
	case "||":
		return boolVal(a) || boolVal(b)
	case "<":
		return num(a) < num(b)
	case ">":
		return num(a) > num(b)
	case "<=":
		return num(a) <= num(b)
	case ">=":
		return num(a) >= num(b)
	case "==":
		return a == b
	case "!=":
		return a != b
	default:
		Error(exp, "Can't apply operator %s", op)
		return nil
	}
}

func makeFun(env *Environment, exp *Expr) func(args ...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		scope := env.Extend()
		for i, name := range exp.Vars {
			if i < len(args) {
				scope.Def(name, args[i])
			} else {
				scope.Def(name, false)
			}
		}

		return evaluate(exp.Body, scope)
	}
}

func Error(exp *Expr, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if exp != nil {
		fmt.Printf("error (%s:%d:%d): %s\n", exp.File, exp.Line, exp.Col, msg)
	} else {
		fmt.Printf("error: %s\n", msg)
	}
	os.Exit(0)
}
