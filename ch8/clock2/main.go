package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

var port = flag.String("port", "3030", "port number")
var tz = flag.String("tz", "Europe/Warsaw", "timezone name")
var loc *time.Location

func main() {
	flag.Parse()
	var err error
	loc, err = time.LoadLocation(*tz)
	if err != nil {
		log.Fatalf("clockwall %v", err)
	}

	listener, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	log.Printf("new connection to port %s by %s", *port, c.RemoteAddr())

	for {
		t := time.Now().In(loc)

		_, err := io.WriteString(c, *tz+" -> "+t.Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
