package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/unnamedxaer/book-gopl/ch7/eval"
)

func main() {

	fmt.Println("Expression:")
	var exprStr string = ""
	_, err := fmt.Scanln(&exprStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("E: ", exprStr)
	x, err := eval.Parse(exprStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	env := eval.Env{}
	fmt.Println("Enter variables one by one in format 'var_name=var_value',\ntype 'EOF' (ctrl + Z) to end supplying variables ")
	for {
		txt := ""
		_, err := fmt.Scanln(&txt)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintln(os.Stderr, "x "+err.Error())
			os.Exit(1)
		}
		if txt == "EOF" {
			break
		}
		if txt == "" {
			continue
		}

		if strings.Contains(txt, "=") == false {
			fmt.Fprintln(os.Stderr, "missing '=', please reatry")
			continue
		}

		vv := strings.Split(txt, "=")
		val, err := strconv.ParseFloat(vv[1], 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "incorrect value %q, please reatry\n", vv[0])
			continue
		}
		env[eval.Var(vv[0])] = val
	}

	fmt.Println(x)
	fmt.Println(env)

	fmt.Println(x.Eval(env))
}
