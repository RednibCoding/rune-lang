package runevm

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const MaxRecursionDepth = 3000 // Between 1000 and 2000 is considered to be a reasonable recursion depth for recursive functions

type Evaluator struct {
	// Keep track of file path that have been imported by the import statement.
	importedPaths  map[string]bool
	recursionDepth int
}

func newEvaluator() *Evaluator {
	e := &Evaluator{importedPaths: make(map[string]bool), recursionDepth: 0}
	return e
}

type ReturnValue struct {
	Value interface{}
}

type BreakValue struct {
	Value bool
}

type ContinueValue struct {
	Value bool
}

func (e *Evaluator) evaluate(exp *expression, env *Environment) interface{} {
	if exp == nil {
		Error(exp, "Null expression error, this is a bug and should never happen!. Please file a bug!")
		return nil
	}

	e.recursionDepth++
	defer func() { e.recursionDepth-- }()

	if e.recursionDepth > MaxRecursionDepth {
		Error(exp, "Maximum recursion depth exceeded")
	}

	switch exp.Type {
	case numExpr:
		return parseNumber(exp.Value.(string), exp)

	case strExpr, boolExpr:
		return exp.Value

	case varExpr:
		value := env.get(exp.Value.(string), exp)
		if exp.Index != nil {
			switch v := value.(type) {
			case []interface{}:
				index := e.evaluate(exp.Index, env).(int)
				if index < 0 || index >= len(v) {
					Error(exp, "Index '%d' out of bounds for array '%v[%d]'", index, exp.Value, len(v))
				}
				return v[index]
			case map[string]interface{}:
				key := e.evaluate(exp.Index, env).(string)
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

	case assignExpr:
		if exp.Left.Type == varExpr && exp.Left.Index != nil {
			arrayOrTable := e.evaluate(exp.Left.Left, env)
			if arrayOrTable == nil {
				Error(exp, "Cannot assign to an index on a null expression")
				return nil
			}
			index := e.evaluate(exp.Left.Index, env)
			value := e.evaluate(exp.Right, env)

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
		} else if exp.Left.Type == varExpr && exp.Left.Value != nil {
			if exp.Left.Left != nil {
				// Handle field assignment
				obj := e.evaluate(exp.Left.Left, env)
				field := exp.Left.Value.(string)
				value := e.evaluate(exp.Right, env)
				if m, ok := obj.(map[string]interface{}); ok {
					m[field] = value
					return value
				} else {
					Error(exp, "Cannot assign to field %v on non-table object", field)
					return nil
				}
			}
			// Handle variable assignment
			return env.set(exp.Left.Value.(string), e.evaluate(exp.Right, env))
		}

		if exp.Left.Type != varExpr {
			Error(exp, "Cannot assign to %v", exp.Left)
		}
		return env.set(exp.Left.Value.(string), e.evaluate(exp.Right, env))

	case binaryExpr:
		a := e.evaluate(exp.Left, env)
		b := e.evaluate(exp.Right, env)
		if a_, ok := a.(ReturnValue); ok {
			a = a_.Value
		}
		if b_, ok := b.(ReturnValue); ok {
			b = b_.Value
		}
		result := applyBinaryOp(exp.Operator, a, b, exp)
		return result

	case unaryExpr:
		return applyUnaryOp(exp.Operator,
			e.evaluate(exp.Right, env), exp)

	case funExpr:
		return e.makeFun(env, exp)

	case ifExpr:
		cond := e.evaluate(exp.Cond, env)
		if cond != false {
			return e.evaluate(exp.Then, env)
		}
		if exp.Else != nil {
			if exp.Else.Type == blockExpr {
				for _, elif := range exp.Else.Block {
					if elif.Type == ifExpr {
						elifCond := e.evaluate(elif.Cond, env)
						if elifCond != false {
							return e.evaluate(elif.Then, env)
						}
					} else {
						return e.evaluate(elif, env)
					}
				}
			} else {
				return e.evaluate(exp.Else, env)
			}
		}
		return false

	case whileExpr:
		for {
			cond := e.evaluate(exp.Cond, env)
			if !cond.(bool) {
				break
			}
			shouldContinue := false
			for _, exp := range exp.Body.Block {
				result := e.evaluate(exp, env)
				if _, ok := result.(BreakValue); ok {
					return false
				} else if _, ok := result.(ContinueValue); ok {
					shouldContinue = true
					break
				}
			}
			if shouldContinue {
				shouldContinue = false
				continue
			}
		}
		return false

	case arrayExpr:
		var arr []interface{}
		for _, element := range exp.Block {
			arr = append(arr, e.evaluate(element, env))
		}
		return arr

	case tableExpr:
		m := make(map[string]interface{})
		for _, pair := range exp.Block {
			key := e.evaluate(pair.Left, env)
			value := e.evaluate(pair.Right, env)
			m[fmt.Sprint(key)] = value
		}
		return m

	case blockExpr:
		var val interface{} = false
		for _, ex := range exp.Block {
			result := e.evaluate(ex, env)
			if ret, ok := result.(ReturnValue); ok {
				return ret.Value
			}
			val = result
		}
		return val

	case callExpr:
		fn, ok := e.evaluate(exp.Func, env).(func(args ...interface{}) interface{})
		if !ok {
			Error(exp, "'%s' is not a function", exp.Func.Value)
		}

		var args []interface{}

		// Check if caller is a go-map/rune-table ...
		callerName := exp.Func.Value
		caller, ok := env.get(callerName.(string), exp).(map[string]interface{})
		// if so:
		if ok {
			// inject its reference as the first argument (similar to pythons 'self' argument on methods)
			args = append(args, caller)
		}

		for _, arg := range exp.Args {
			args = append(args, e.evaluate(arg, env))
		}

		ret := fn(args...)
		if err, ok := ret.(error); ok {
			Error(exp, "Error in function call: '%v'", err)
		}
		return ret

	case returnExpr:
		return ReturnValue{Value: e.evaluate(exp.Right, env)}

	case breakExpr:
		return BreakValue{Value: false}

	case continueExpr:
		return ContinueValue{Value: false}

	case importExpr:
		path := e.evaluate(exp.Left, env).(string) + ".rune"
		if _, alreadyImported := e.importedPaths[path]; alreadyImported {
			Error(exp, "Duplicate import detected: '%s' was already imported", path)
		}

		e.importedPaths[path] = true

		importedSource, err := os.ReadFile(path)
		if err != nil {
			Error(exp, "Failed to import file '%s': %v", path, err)
		}

		importStream := newInputStream(string(importedSource), path)
		importTokenStream := newTokenStream(importStream)
		importParser := newParser(importTokenStream)
		importAST := importParser.parseProgram()

		e.evaluate(importAST, env)
		return nil

	default:
		Error(exp, "I don't know how to evaluate %v", exp.Type)
		return nil
	}
}

func parseNumber(val string, exp *expression) interface{} {
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

func applyUnaryOp(op string, a interface{}, exp *expression) interface{} {
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
			if x != nil {
				return false
			} else {
				Error(exp, fmt.Sprintf("Unary operator '%s' needs a valid operand", op))
			}
		}
		return false
	}
	switch op {
	case "not":
		return !boolVal(a)
	default:
		Error(exp, "Can't apply unary operator %s", op)
		return nil
	}
}

func applyBinaryOp(op string, a, b interface{}, exp *expression) interface{} {
	num := func(x interface{}) float64 {
		switch v := x.(type) {
		case string:
			return parseNumber(v, exp).(float64)
		case int:
			return float64(v)
		case int32:
			return float64(v)
		case int64:
			return float64(v)
		case float32:
			return float64(v)
		case float64:
			return float64(v)
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
	roundIfInt := func(value float64) interface{} {
		if math.Abs(value-math.Round(value)) < 1e-9 {
			return int(value)
		}
		return value
	}

	switch op {
	case "+":
		return roundIfInt(num(a) + num(b))
	case "-":
		return roundIfInt(num(a) - num(b))
	case "*":
		return roundIfInt(num(a) * num(b))
	case "/":
		return roundIfInt(num(a) / div(b))
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
		return num(a) == num(b)
	case "!=":
		return num(a) != num(b)
	default:
		Error(exp, "Can't apply operator %s", op)
		return nil
	}
}

func (e *Evaluator) makeFun(env *Environment, exp *expression) func(args ...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		scope := env.extend()
		// Collect and define all param names in the current scope/environment so they are known within the function
		for i, name := range exp.Params {
			if i < len(args) {
				scope.def(name, args[i])
			} else {
				scope.def(name, false)
			}
		}

		return e.evaluate(exp.Body, scope)
	}
}

func Error(exp *expression, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if exp != nil {
		fmt.Printf("error (%s:%d:%d): %s\n", exp.File, exp.Line, exp.Col, msg)
	} else {
		fmt.Printf("error: %s\n", msg)
	}
	os.Exit(0)
}
