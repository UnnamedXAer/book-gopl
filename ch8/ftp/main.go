// Package Ftp
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// var (
// 	cd  = flag.String("cd", "", "path")
// 	ls  = flag.String("ls", "", "path")
// 	get = flag.String("get", "", "file path")
// )

var basePath string

func main() {
	var w io.Writer = os.Stdout
	var r io.Reader = os.Stdin
	cwd, err := os.Getwd()
	basePath = filepath.Clean(filepath.Join(cwd, "base"))
	cwd = basePath
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(basePath)
	log.Println("Welcome to ftp")
	fmt.Print("\n> ")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		cmd := scanner.Text()
		cmdParts := strings.Split(cmd, "=")
		if len(cmdParts) > 2 {
			fmt.Fprintf(w, "error: not recognized command %q\n> ", cmd)
			continue
		}
		if len(cmdParts) == 0 {
			continue
		}
		name := cmdParts[0]

		if len(cmdParts) == 1 {
			switch name {
			case "close":
				log.Fatalln(w, "disconnecting...")
				os.Exit(0)
			case "ls":
				s, err := ls(cwd, "")
				if err != nil {
					fmt.Fprintf(w, "error: %s\n> ", err)
					continue
				}
				fmt.Fprintf(w, "%s\n> ", s)
			default:
				fmt.Fprintf(w, "error: not recognized command %q\n> ", cmdParts[0])
			}
			continue
		}
		value := cmdParts[1]
		switch name {
		case "ls":
			s, err := ls(cwd, value)
			if err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}
			fmt.Fprintf(w, "%s\n", s)
		case "cd":
			path := filepath.Join(cwd, value)
			if err := checkAgainstBasePath(path); err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}
			_, err := os.Stat(path)
			if err != nil {
				fmt.Fprintln(w, "error", err)
				continue
			}
			cwd = path
			fmt.Fprintf(w, "%s\n> ", path)
		case "mkdir":
			path := filepath.Join(cwd, value)
			if err := checkAgainstBasePath(path); err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}
			err := os.MkdirAll(path, 0644)
			if err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}
		case "get":
			valParts := strings.Split(value, ">")
			if len(valParts) != 2 {
				fmt.Fprintf(w, "error: invalid value for get command %q\n> ", value)
				continue
			}
			srcPath := filepath.Join(cwd, valParts[0])
			if err := checkAgainstBasePath(srcPath); err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}

			b, err := ioutil.ReadFile(srcPath)
			if err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}

			dstPath := filepath.Join("./", valParts[1])
			err = ioutil.WriteFile(dstPath, b, 0644)
			if err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				continue
			}
		default:
			fmt.Fprintf(w, "error: not recognized command %q\n> ", cmdParts[0])
		}
	}
}

func ls(cwd, value string) (string, error) {
	path := filepath.Join(cwd, value)
	if err := checkAgainstBasePath(path); err != nil {
		return "", err
	}
	d, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	s := ""
	for i, f := range d {
		dirChar := ""
		if f.IsDir() {
			dirChar = "/"
		}
		s += fmt.Sprintf("% 3d. %s%s\n", i, f.Name(), dirChar)
	}

	return s, nil
}

func checkAgainstBasePath(path string) error {
	fmt.Println(path)
	if strings.HasPrefix(filepath.Clean(path), basePath) == false {
		return fmt.Errorf("nope, cannot go there %q", path)
	}
	return nil
}
