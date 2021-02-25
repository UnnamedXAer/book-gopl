package tempflag

import (
	"flag"
	"fmt"

	"github.com/unnamedxaer/book-gopl/ch1/tempconv"
)

// *celsiusFlag satisfies the flag.Value interface
type celsiusFlag struct {
	Celsius tempconv.Celsius
}

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = tempconv.Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = tempconv.FToC(tempconv.Fahrenheit(value))
		return nil
	case "K":
		f.Celsius = tempconv.KToC(tempconv.Kelvin(value))
		return nil
	}

	return fmt.Errorf("invalid temperature %q", s)
}

func (f celsiusFlag) String() string {
	return f.Celsius.String()
}

//CelsiusFlag defines a Celsius flat with the specified name
// default value, and usage, and returns the address of the flag varable.
// The flag argument must have a quantity and a unit, e.g., "100C"
func CelsiusFlag(name string, value tempconv.Celsius, usage string) *tempconv.Celsius {
	f := celsiusFlag{Celsius: value}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

// var temp = CelsiusFlag("temp", 20.0, "the temperature")

// func main() {
// 	flag.Parse()
// 	fmt.Println(temp)
// }
