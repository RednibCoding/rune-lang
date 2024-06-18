package runevm

type RuneVM struct {
	filepath string
	source   string
	env      *Environment
}

func NewRuneVM() *RuneVM {
	vm := &RuneVM{}

	vm.env = NewEnvironment(nil)
	vm.env.Def("print", builtin_Print)
	vm.env.Def("println", builtin_Println)

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

func (r *RuneVM) RegisterFunction(name string, fn func(args ...interface{}) interface{}) {
	r.env.Def(name, fn)
}

func (r *RuneVM) RegisterVariable(name string, value interface{}) {
	r.env.Def(name, value)
}
