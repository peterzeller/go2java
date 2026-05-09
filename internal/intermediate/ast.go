package intermediate

// Program is a language-neutral AST carrying both Go and Java-friendly constructs.
type Program struct {
	PackageName string
	Records     []*Record
	Functions   []*Function
}

type Record struct {
	Name    string
	Fields  []Parameter
	Methods []*Function
}

type Function struct {
	Name         string
	Parameters   []Parameter
	ReturnType   Type
	Body         []Stmt
	ReceiverName string
	ReceiverType Type
}

type Parameter struct {
	Name   string
	Type   Type
	Public bool
}

type Type string

const (
	TypeVoid   Type = "void"
	TypeInt    Type = "int"
	TypeBool   Type = "boolean"
	TypeString Type = "String"
)

type Stmt interface{ isStmt() }
type Expr interface{ isExpr() }

type ExprStmt struct{ Expr Expr }

func (*ExprStmt) isStmt() {}

type ReturnStmt struct{ Value Expr }

func (*ReturnStmt) isStmt() {}

type AssignStmt struct {
	Name  string
	Type  Type
	Value Expr
}

type FieldAssignStmt struct {
	Target Expr
	Field  string
	Value  Expr
}

func (*FieldAssignStmt) isStmt() {}

func (*AssignStmt) isStmt() {}

type IfStmt struct {
	Cond Expr
	Then []Stmt
	Else []Stmt
}

func (*IfStmt) isStmt() {}

type ForStmt struct {
	Cond Expr
	Body []Stmt
}

func (*ForStmt) isStmt() {}

type BinaryExpr struct {
	Op          string
	Left, Right Expr
}

func (*BinaryExpr) isExpr() {}

type IdentExpr struct{ Name string }

func (*IdentExpr) isExpr() {}

type SelectorExpr struct {
	Target Expr
	Field  string
}

func (*SelectorExpr) isExpr() {}

type IntLiteral struct{ Value string }

func (*IntLiteral) isExpr() {}

type StringLiteral struct{ Value string }

func (*StringLiteral) isExpr() {}

type BoolLiteral struct{ Value bool }

func (*BoolLiteral) isExpr() {}

type CompositeLiteral struct {
	TypeName string
	Args     []Expr
}

func (*CompositeLiteral) isExpr() {}

type CallExpr struct {
	Func Expr
	Args []Expr
}

func (*CallExpr) isExpr() {}
