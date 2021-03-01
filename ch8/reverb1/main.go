package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("started accepting connections")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn)
	}
	fmt.Println("main end")
}

func handleConn(c net.Conn) {
	fmt.Println("new connection:", c.RemoteAddr())
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 1*time.Second)
	}

	checkErr(input.Err())
	c.Close()
}

func echo(c net.Conn, shout string, delay time.Duration) {
	_, err := fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	checkErr(err)
	time.Sleep(delay)
	_, err = fmt.Fprintln(c, "\t", strings.Title(shout))
	checkErr(err)
	time.Sleep(delay)
	_, err = fmt.Fprintln(c, "\t", strings.ToLower(shout))
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
