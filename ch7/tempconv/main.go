package main

import (
	"flag"
	"fmt"

	"github.com/unnamedxaer/book-gopl/ch1/tempconv"
)

// *celsiusFlag satisfies the flag.Value interface
type celsiusFlag struct {
	tempconv.Celsius
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // (surprisingly) no error check needed?
	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	}

	return fmt.Errorf("invalid temperatur %q", s)
}

//CelsiusFlag...

func main() {
	flag.Parse()

}
