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

func echo(c *net.TCPConn, shout string, delay time.Duration, clearTimer chan struct{}) {
	_, err := fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	checkErr(1, err)
	time.Sleep(delay)
	_, err = fmt.Fprintln(c, "\t", shout)
	checkErr(2, err)
	time.Sleep(delay)
	_, err = fmt.Fprintln(c, "\t", strings.ToLower(shout))
	checkErr(3, err)

	log.Println(time.Now().Format("15:04:05"))
	go timeout(clearTimer, c)

}

func timeout(clearTimer chan struct{}, c *net.TCPConn) {
	timeout := 5 * time.Second
	endTime := time.Now().Add(timeout).Format("15:04:05")
	log.Println("will timetout around", endTime)

	select {
	case <-time.After(timeout):
		log.Println("closing by timeout", endTime)
		err := c.Close()
		if err != nil {
			log.Printf("conn close: %v\n", err)
		}
		clearTimer = nil
	case <-clearTimer:
		log.Println("timeout cleared", endTime)

	}

	log.Println("exiting timeout", endTime)
}

func handleConn(c *net.TCPConn) {
	log.Println(c.RemoteAddr().String(), " - connected")
	input := bufio.NewScanner(c)
	clearTimer := make(chan struct{})
	go timeout(clearTimer, c)
	for input.Scan() {
		log.Println("T:", input.Text())
		clearTimer <- struct{}{}
		go echo(c, input.Text(), 1*time.Second, clearTimer)
	}

	log.Println("input: ", input.Err().Error())

	// NOTE: ignoring potential errors from input.Err()
	log.Println("final closing", c.RemoteAddr().String())
	err := c.Close()
	if err != nil {
		log.Printf("conn close: %v\n", err)
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
		c := conn.(*net.TCPConn)
		go handleConn(c)
	}
}

func checkErr(d int, err error) {
	if err != nil {
		log.Printf("checkErr (%d): %v\n", d, err)
	}
}
