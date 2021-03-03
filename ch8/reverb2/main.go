package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
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
	wg := sync.WaitGroup{}

	fmt.Println("new connection:", c.RemoteAddr())
	input := bufio.NewScanner(c)
	for input.Scan() {
		wg.Add(1)
		fmt.Println("ADD:", 1, c.RemoteAddr().String())
		go echo(c, input.Text(), 1*time.Second, &wg)
	}
	fmt.Println("after loop")

	go func() {
		wg.Wait()
		log.Println("wg.Waited long enough !", c.RemoteAddr().String())
		c.(*net.TCPConn).CloseWrite()
		fmt.Println("c - CloseWriter", c.RemoteAddr().String())
	}()

}

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Done:", c.RemoteAddr().String())
		wg.Done()
	}()
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
