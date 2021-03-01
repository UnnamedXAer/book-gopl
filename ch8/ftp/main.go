// Package Ftp
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var basePath string

var cwd map[string]string

func main() {
	d, err := os.Getwd()
	if err != nil {
		log.Fatalf("%q", err)
		return
	}
	basePath = filepath.Clean(filepath.Join(d, "base"))
	cwd = make(map[string]string, 10)

	// as a server
	listener, err := net.Listen("tcp", "localhost:3031")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handleCommand(conn)
	}
}

func handleCommand(w net.Conn) {
	if cwd[w.RemoteAddr().String()] == "" {
		d, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(w, "error: internal filesystem problem %q\n> ", err)
			return
		}
		cwd[w.RemoteAddr().String()] = filepath.Join(d, "base")
	}
	fmt.Println(basePath)
	log.Println("Welcome to ftp")
	fmt.Print("\n> ")
	scanner := bufio.NewScanner(w)
	for scanner.Scan() {
		cmd := scanner.Text()
		execCmd(w, cmd)
	}
}

func execCmd(w net.Conn, cmd string) {
	log.Printf("command: %q\n", cmd)
	cmdParts := strings.Split(cmd, "=")
	if len(cmdParts) > 2 {
		fmt.Fprintf(w, "error: not recognized command %q\n> ", cmd)
		return
	}
	if len(cmdParts) == 0 || (len(cmdParts) == 1 && cmdParts[0] == "") {
		return
	}
	name := cmdParts[0]

	if len(cmdParts) == 1 {
		switch name {
		case "close":
			log.Fatalln(w, "disconnecting...")
			os.Exit(0)
		case "ls":
			s, err := ls(cwd[w.RemoteAddr().String()], "")
			if err != nil {
				fmt.Fprintf(w, "error: %s\n> ", err)
				return
			}
			fmt.Fprintf(w, "%s\n> ", s)
		default:
			fmt.Fprintf(w, "error: not recognized command %q\n> ", cmdParts[0])
		}
		return
	}
	value := cmdParts[1]
	switch name {
	case "ls":
		s, err := ls(cwd[w.RemoteAddr().String()], value)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}
		fmt.Fprintf(w, "%s\n", s)
	case "cd":
		path := filepath.Join(cwd[w.RemoteAddr().String()], value)
		if err := checkAgainstBasePath(path); err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}
		_, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(w, "error", err)
			return
		}
		cwd[w.RemoteAddr().String()] = path
		fmt.Fprintf(w, "%s\n> ", path)
	case "mkdir":
		path := filepath.Join(cwd[w.RemoteAddr().String()], value)
		if err := checkAgainstBasePath(path); err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}
		err := os.MkdirAll(path, 0666)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}
		fmt.Fprintf(w, "%s\n> ", path)

	case "get":
		valParts := strings.Split(value, ">")
		if len(valParts) != 2 {
			fmt.Fprintf(w, "error: invalid value for get command %q\n> ", value)
			return
		}
		srcPath := filepath.Join(cwd[w.RemoteAddr().String()], valParts[0])
		if err := checkAgainstBasePath(srcPath); err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}

		b, err := ioutil.ReadFile(srcPath)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}

		dstPath := filepath.Join("./", valParts[1])
		err = ioutil.WriteFile(dstPath, b, 0644)
		if err != nil {
			fmt.Fprintf(w, "error: %s\n> ", err)
			return
		}
		fmt.Fprintf(w, "%s\n> ", dstPath)
	default:
		fmt.Fprintf(w, "error: not recognized command %q\n> ", cmdParts[0])
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
