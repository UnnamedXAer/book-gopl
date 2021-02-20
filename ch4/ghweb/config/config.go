package config

import (
	"os"
	"strconv"

	// config is the very first thing that needs env vars
	_ "github.com/joho/godotenv/autoload"
)

// config keeps app configuration data
type config struct {
	Debbug bool
	PORT   uint16
}

// C keeps global app configuration
var C config

func init() {
	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	debug := os.Getenv("debug") == "TRUE" || os.Getenv("debug") == "true"
	if err != nil {
		if debug == false {
			panic(`missing "PORT" environment variable`)
		}
		port = 3030
	}

	C = config{
		Debbug: debug,
		PORT:   uint16(port),
	}
}
