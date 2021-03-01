package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {

	go showTime("3020")
	go showTime("3030")
	showTime("3040")

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalln(err)
	}
}

func showTime(port string) {
	// fmt.Println("port : " + port)
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}
