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

// -----------------------------------------------------

type Variable struct {
	Name Token
}

func (v *Variable) isExpr() {}

// -----------------------------------------------------

type Assign struct {
	Name  Token
	Value Expr
}

func (a *Assign) isExpr() {}

// -----------------------------------------------------

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l *Logical) isExpr() {}
