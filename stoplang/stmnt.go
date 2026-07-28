package stoplang

type Stmt interface {
	Accept(visitor StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) (any, error)
	VisitPrintStmt(stmt *PrintStmt) (any, error)
	VisitVarStmt(stmt *VarStmt) (any, error)
	VisitBlockStmt(stmt *BlockStmt) (any, error)
	VisitIfStmt(stmt *IfStmt) (any, error)
	VisitWhileStmt(stmt *WhileStmt) (any, error)
}

type ExpressionStmt struct {
	Expression Expr
}

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

type ForStmt struct {
	Initializer Expr
	Condition   Expr
	Body        Stmt
}

func (b *BlockStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitBlockStmt(b)
}

func (s *ExpressionStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (v *VarStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitVarStmt(v)
}

func (s *PrintStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitPrintStmt(s)
}

func (s *IfStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitIfStmt(s)
}

func (s *WhileStmt) Accept(visitor StmtVisitor) (any, error) {
	return visitor.VisitWhileStmt(s)
}
