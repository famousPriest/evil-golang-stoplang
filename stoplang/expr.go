package stoplang

type Expr interface {
	isExpr()
}

// -----------------------------------------------------

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) isExpr() {}

// -----------------------------------------------------

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) isExpr() {}

// -----------------------------------------------------

type Grouping struct {
	Expression Expr
}

func (g *Grouping) isExpr() {}

// -----------------------------------------------------

type Literal struct {
	Value any
}

func (l *Literal) isExpr() {}
