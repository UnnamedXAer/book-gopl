package eval

import (
	"fmt"
	"math"
	"strings"
)

// An Expr in an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) float64
	// Check reports error in this Expr and adds its Vars to the set.
	Check(vars map[Var]bool) error
	String() string
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

// postUnary represents a unary operator expresson, e.g., -x.
type postUnary struct {
	op rune // one of '!'
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
	if x, ok := env[v]; ok {
		return x
	}
	panic(fmt.Errorf("undefined %s", v))
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

func (p postUnary) Eval(env Env) float64 {
	switch p.op {
	case '!':
		return math.Gamma(p.x.Eval(env))
	}
	panic(fmt.Sprintf("unsupported postunary operator: %q", p.op))
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

func (v Var) Check(vars map[Var]bool) error {
	fmt.Println("var -> ", v)
	vars[v] = true
	return nil
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-<>=", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (p postUnary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("!", p.op) {
		return fmt.Errorf("unexpected postunary op %q", p.op)
	}
	return p.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/<>=", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s gas %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1, "min": 2}

func (v Var) String() string {
	return string(v[0])

}
func (v literal) String() string {
	s := fmt.Sprintf("%f", float64(v))
	return s
}
func (v unary) String() string {
	s := fmt.Sprintf("%s%s", string(v.op), v.x)
	return s
}

func (v postUnary) String() string {
	s := fmt.Sprintf("%s%s", v.x, string(v.op))
	return s
}

func (v binary) String() string {
	s := fmt.Sprintf("%s %s %s", v.x, string(v.op), v.y)
	return s
}

func (v call) String() string {
	if len(v.args) == 2 {
		s := fmt.Sprintf("%s(%s, %s)", v.fn, v.args[0], v.args[1])
		return s
	}
	s := fmt.Sprintf("%s(%s)", v.fn, v.args[0])
	return s
}
