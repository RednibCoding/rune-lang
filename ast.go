package runevm

type exprType string

const (
	numExpr      exprType = "num"
	strExpr      exprType = "str"
	boolExpr     exprType = "bool"
	varExpr      exprType = "var"
	assignExpr   exprType = "assign"
	binaryExpr   exprType = "binary"
	unaryExpr    exprType = "unary"
	funExpr      exprType = "fun"
	ifExpr       exprType = "if"
	blockExpr    exprType = "block"
	callExpr     exprType = "call"
	returnExpr   exprType = "return"
	whileExpr    exprType = "while"
	breakExpr    exprType = "break"
	continueExpr exprType = "continue"
	arrayExpr    exprType = "array"
	tableExpr    exprType = "table"
	pairExpr     exprType = "pair"
	indexExpr    exprType = "Index"
	importExpr   exprType = "import"
)

type expression struct {
	Type exprType
	// Multipurpose field for storing a value
	Value interface{}

	Left  *expression
	Right *expression

	// Operator of binary expressions
	Operator string

	// If/While
	Cond *expression
	Then *expression
	Else *expression

	// Function decl
	Func *expression
	// Function decl param names
	Params []string

	// Entire block
	Block []*expression

	// Function call arguments
	Args []*expression

	// Function / while bodies
	Body *expression

	// Index access / Field access
	Index *expression

	// Token infos
	File   string
	Line   int
	Col    int
	Length int
}
