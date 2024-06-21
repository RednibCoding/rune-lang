package runevm

type ExprType string

const (
	Num    ExprType = "num"
	Str    ExprType = "str"
	Bool   ExprType = "bool"
	Var    ExprType = "var"
	Assign ExprType = "assign"
	Binary ExprType = "binary"
	Unary  ExprType = "unary"
	Fun    ExprType = "fun"
	If     ExprType = "if"
	Prog   ExprType = "prog"
	Call   ExprType = "call"
	While  ExprType = "while"
	Array  ExprType = "array"
	Table  ExprType = "table"
	Pair   ExprType = "pair"
	Index  ExprType = "Index"
	Import ExprType = "import"
)

type Expr struct {
	Type ExprType
	// Multipurpose field for storing a value
	Value interface{}

	Left  *Expr
	Right *Expr

	// Operator of binary expressions
	Operator string

	// If/While
	Cond *Expr
	Then *Expr
	Else *Expr

	// Function decl
	Func *Expr
	// Function decl param names
	Params []string

	// Entire block
	Prog []*Expr

	// Function call arguments
	Args []*Expr

	// Function / while bodies
	Body *Expr

	// Index access / Field access
	Index *Expr

	// Token infos
	File   string
	Line   int
	Col    int
	Length int
}
