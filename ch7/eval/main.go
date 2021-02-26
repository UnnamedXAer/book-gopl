package eval

import (
	"fmt"
	"math"
)

// An Expr in an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
}

// Var identifies a variable, e.g., x.
type Var string

// literal is a numeric constant, e.g., 3.121.
type literal float64

// unary represents a unary operator expresson, e.g., -x.
type unary struct {
	op rune // one of '+', '-'
	x  Expr
}

// binary represents a binary operator expresson, e.g., y+x.
type binary struct {
	op   rune // one of 'x', '-', '*', '/'
	x, y Expr
}

// call represents a function call expression e.g., sin(x)
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}

// Env is an environment that maps variable names to values
type Env map[Var]float64

// Eval - see Expr interface
func (v Var) Eval(env Env) float64 {
	return env[v]
}

// Eval - see Expr interface
func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %s", c.fn))
}
