// Netcat1 is a read-only TCP clent.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var port = flag.String("port", "3030", "port number")

func main() {
	flag.Parse()
	fmt.Println("port : " + *port)
	conn, err := net.Dial("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalln(err)
	}
}
