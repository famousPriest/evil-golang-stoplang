package thelang

// TODO: i can just use type switch, its not java anyway
type Visitor interface {
	VisitBinaryExpr(expr *Binary) (any, error)
	VisitUnaryExpr(expr *Unary) (any, error)
	VisitGroupingExpr(expr *Grouping) (any, error)
	VisitLiteralExpr(expr *Literal) (any, error)
}

type Expr interface {
	Accept(visitor Visitor) (any, error)
}

// -----------------------------------------------------

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(v Visitor) (any, error) {
	return v.VisitBinaryExpr(b)
}

// -----------------------------------------------------

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(v Visitor) (any, error) {
	return v.VisitUnaryExpr(u)
}

// -----------------------------------------------------

type Grouping struct {
	Expression Expr
}

func (u *Grouping) Accept(v Visitor) (any, error) {
	return v.VisitGroupingExpr(u)
}

// -----------------------------------------------------

type Literal struct {
	Value any
}

func (u *Literal) Accept(v Visitor) (any, error) {
	return v.VisitLiteralExpr(u)
}
