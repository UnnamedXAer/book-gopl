package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp4", ":3030")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn) // temporarily ignore errors
		if err != nil {
			fmt.Fprintln(os.Stderr, "copy err: ", err)
		}
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.(*net.TCPConn).CloseWrite()
	log.Println("writer closed")
	<-done // wait for background goroutine to finish
	log.Println("done main")
}

func mustCopy(dst io.Writer, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	log.Println("must copy - end")
}
