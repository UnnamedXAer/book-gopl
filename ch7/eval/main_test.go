package eval

import (
	"fmt"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A/pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x,3)", Env{"x": 12, "y": 1}, "1728"},
		{"5/9 * (F-32)", Env{"F": -40}, "-40"},
		{"5/9 * (F-32)", Env{"F": 32}, "0"},
		{"5/9 * (F-32)", Env{"F": 212}, "100"},
		{"sqrt(1,2)", Env{}, "1"},
		// {"-1 + -x", Env{"x": 1}, "-2"},
		// {"-1 - x", Env{"x": 1}, "-2"},
	}

	var txt string

	var prevExp string
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExp {
			fmt.Printf("\n%s\n", test.expr)
			prevExp = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		if got != test.want {
			txt += fmt.Sprintf("%s.Eval() in %v = %q, want %q\n", test.expr, test.env, got, test.want)
		}
	}

	if txt != "" {
		t.Errorf(txt)
	}
}

func TestCoverage(t *testing.T) {
	var tests = []struct {
		input string
		env   Env
		want  string // expected error from Parse/Check or result from Eval
	}{
		{"x % 2", nil, "unexpected '%'"},
		{"!true", nil, "unexpected '!'"},
		{"log(10)", nil, `unknown function "log"`},
		{"sqrt(1, 2)", nil, "call to sqrt has 2 args, want 1"},
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
	}

	for _, test := range tests {
		expr, err := Parse(test.input)
		if err == nil {
			err = expr.Check(map[Var]bool{})
		}
		if err != nil {
			if err.Error() != test.want {
				t.Errorf("%s: got %q, want %q", test.input, err, test.want)
			}
			continue
		}

		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		if got != test.want {
			t.Errorf("%s: %v => %s, want %s",
				test.input, test.env, got, test.want)
		}
	}
}
