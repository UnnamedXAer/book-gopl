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
	go func() {
		wg.Wait()
		log.Println("wg.Waited long enough !")
	}()
	defer wg.Done()
	fmt.Println("new connection:", c.RemoteAddr())
	input := bufio.NewScanner(c)
	for input.Scan() {
		wg.Add(1)
		go echo(c, input.Text(), 1*time.Second, &wg)
	}

	checkErr(input.Err())
	c.(*net.TCPConn).CloseWrite()
}

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
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
