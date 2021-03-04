// Netcat is a simple read/write client for TCP servers.
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
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "copy err: ", err)
		}
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	mustCopy(conn, os.Stdin)
	conn.(*net.TCPConn).CloseWrite()
	log.Println("writer closed")
	<-done // wait for background goroutine to finish
	log.Println("done main")
}

func mustCopy(dst io.Writer, src io.Reader) {
	log.Println("must copy")
	_, err := io.Copy(dst, src)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	log.Println("must copy - end")
}
