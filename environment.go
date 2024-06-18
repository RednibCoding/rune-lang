package runevm

type Environment struct {
	vars   map[string]interface{}
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	vars := make(map[string]interface{})
	return &Environment{vars: vars, parent: parent}
}

func (env *Environment) Extend() *Environment {
	return NewEnvironment(env)
}

func (env *Environment) Lookup(name string) *Environment {
	for scope := env; scope != nil; scope = scope.parent {
		if _, found := scope.vars[name]; found {
			return scope
		}
	}
	return nil
}

func (env *Environment) Get(name string, exp *Expr) interface{} {
	if value, found := env.vars[name]; found {
		return value
	}
	if env.parent != nil {
		return env.parent.Get(name, exp)
	}
	Error(exp, "Undefined variable '%s'", name)
	return nil
}

func (env *Environment) Set(name string, value interface{}, exp *Expr) interface{} {
	scope := env.Lookup(name)
	if scope == nil && env.parent != nil {
		Error(exp, "Undefined variable '%s'", name)
	}
	if scope != nil {
		scope.vars[name] = value
	} else {
		env.vars[name] = value
	}
	return value
}

func (env *Environment) Def(name string, value interface{}) interface{} {
	env.vars[name] = value
	return value
}
