// Reverb is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, disconnect chan time.Time) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
	go func() {
		x := <-time.After(10 * time.Second)
		select {
		case disconnect <- x:
			c.(*net.TCPConn).Close()
			log.Println(c.RemoteAddr().String(), " - CloseWrite")
		}

	}()
}

func handleConn(c net.Conn) {
	log.Println(c.RemoteAddr().String(), " - connected")
	disconnect := make(chan time.Time)
	input := bufio.NewScanner(c)

l:
	for input.Scan() {
		// disconnect = nil
		go echo(c, input.Text(), 1*time.Second, disconnect)
		select {
		case <-disconnect:
			log.Println("break")
			break l
		default:
		}
	}
	// NOTE: ignoring potential errors from input.Err()
	log.Println("final closing")
	err := c.Close()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(c.RemoteAddr().String(), " - closed")
}

func main() {
	l, err := net.Listen("tcp", "localhost:3030")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
