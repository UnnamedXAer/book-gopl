package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", "localhost:3030")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn) // temporarily ignore errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.(*net.TCPConn).CloseWrite()
	log.Println("before done", 1)
	<-done // wait for background goroutine to finish
	log.Println("done main")
}

func mustCopy(dst io.Writer, src io.Reader) {
	fmt.Println("must copy")
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("must copy - end")
}
