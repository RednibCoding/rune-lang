package runevm

type Environment struct {
	vars   map[string]interface{}
	parent *Environment
}

func newEnvironment(parent *Environment) *Environment {
	vars := make(map[string]interface{})
	return &Environment{vars: vars, parent: parent}
}

func (env *Environment) extend() *Environment {
	return newEnvironment(env)
}

func (env *Environment) lookup(name string) *Environment {
	for scope := env; scope != nil; scope = scope.parent {
		if _, found := scope.vars[name]; found {
			return scope
		}
	}
	return nil
}

func (env *Environment) get(name string, exp *expression) interface{} {
	if value, found := env.vars[name]; found {
		return value
	}
	if env.parent != nil {
		return env.parent.get(name, exp)
	}
	Error(exp, "Undefined variable '%s'", name)
	return nil
}

// func (env *Environment) set(name string, value interface{}, exp *Expr) interface{} {
func (env *Environment) set(name string, value interface{}) interface{} {
	scope := env.lookup(name)
	// if scope == nil && env.parent != nil {
	// 	Error(exp, "Undefined variable '%s'", name)
	// }
	if scope != nil {
		scope.vars[name] = value
	} else {
		env.vars[name] = value
	}
	return value
}

func (env *Environment) def(name string, value interface{}) interface{} {
	env.vars[name] = value
	return value
}
