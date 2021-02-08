package tempconv

// CToF coverts a Celsius temperature to Fahrenheit
func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// FToC coverts a Fahrenheit temperature to Celsius
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// CToK coverts a Celsius temperature to Kelvin
func CToK(c Celsius) Kelvin {
	return Kelvin(c + 275.15)
}

// FToK coverts a Kelvin temperature to Celsius
func FToK(k Kelvin) Celsius {
	return Celsius(k - 275.15)
}
