package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println(err)
		}

		go handleConn(conn)
	}
}

type client struct {
	ch  chan<- string // an outgoing message chanel
	who string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	messages = make(chan string) // all incoming client messages
)

func broadcaster() {
	clients := make(map[*client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// client's outgoing message channels
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			fmt.Println("enter: ", cli)
			for idx := range clients {
				cli.ch <- fmt.Sprintf("idx %v, v: %v", idx, idx.who)
			}
			clients[cli] = true

		case cli := <-leaving:
			fmt.Println("leave: ", cli)
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	me := &client{
		ch:  ch,
		who: who,
	}
	entering <- me

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// note ignoring potential errors from input.Err()

	leaving <- me
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ingoring network errors
	}
}
