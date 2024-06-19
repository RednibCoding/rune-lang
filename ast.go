package runevm

type ExprType string

type KeyValuePair struct {
	Key   *Expr
	Value *Expr
}

const (
	Num    ExprType = "num"
	Str    ExprType = "str"
	Bool   ExprType = "bool"
	Var    ExprType = "var"
	Assign ExprType = "assign"
	Binary ExprType = "binary"
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
	Type     ExprType
	Value    interface{}
	Left     *Expr
	Right    *Expr
	Operator string
	Cond     *Expr
	Then     *Expr
	Else     *Expr
	Prog     []*Expr
	Func     *Expr
	Args     []*Expr
	Vars     []string
	Body     *Expr
	File     string
	Line     int
	Col      int
	Length   int
	Index    *Expr
}
