package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func contac1tHandler(w http.ResponseWriter, r *http.Request) {
	// validate params
	// if (param == "" || ...) {...}
	data := map[string]map[string]string{
		"values": {
			"mail":    "example@@email.com",
			"message": "",
		},
		"errors": {
			"mail":    "email address is not correct",
			"message": "message is required",
		},
	}
	// I do not call tpl.ExecuteTemplate directly to http.ResponseWriter because
	// in case of error the error message we want to show would
	// be appended to rendered template (the part before error occurs)
	b := bytes.Buffer{}
	err := contactTpl.ExecuteTemplate(&b, "bootstrap", data)
	if err != nil {
		errText := http.StatusText(http.StatusInternalServerError)
		if config.Dev {
			errText = err.Error()
		}
		// or render error page
		http.Error(
			w,
			errText,
			http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, &b)
}
