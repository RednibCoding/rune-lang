package runevm

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func evaluate(exp *Expr, env *Environment) interface{} {
	switch exp.Type {
	case Num:
		return parseNumber(exp.Value.(string), exp)

	case Str, Bool:
		return exp.Value

	case Var:
		value := env.Get(exp.Value.(string), exp)
		if exp.Index != nil {
			array, ok := value.([]interface{})
			if !ok {
				Error(exp, "Variable %v is not an array", exp.Value)
			}
			index := evaluate(exp.Index, env).(int)
			if index < 0 || index >= len(array) {
				Error(exp, "Index '%d' out of bounds for array '%v[%d]'", index, exp.Value, len(array))
			}
			return array[index]
		}
		return value

	case Assign:
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
		for _, arg := range exp.Args {
			args = append(args, evaluate(arg, env))
		}
		ret := fn(args...)
		if err, ok := ret.(error); ok {
			Error(exp, "Error in function call: '%v'", err)
		}
		return ret

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
