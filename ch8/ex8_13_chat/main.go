package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server listens on :3030")
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
			// fmt.Println("enter: ", cli)
			for idx := range clients {
				cli.ch <- fmt.Sprintf("idx %v, v: %v", idx, idx.who)
			}
			clients[cli] = true

		case cli := <-leaving:
			// fmt.Println("leave: ", cli)
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

var timeoutDuration = 10 * time.Second

func handleConn(conn net.Conn) {
	ch := make(chan string)          // outgoing client messages
	closeSignal := make(chan string) // chanel to send close message

	who := conn.RemoteAddr().String()
	fmt.Println(who + " connected")
	go clientWriter(conn, ch, who)
	ch <- "You are " + who
	messages <- who + " has arrived"
	cli := &client{
		ch:  ch,
		who: who,
	}
	entering <- cli

	timer := time.NewTimer(timeoutDuration)

	go func() {
		select {
		case <-timer.C:
			ch <- "You have been disconnected due to idleness"
			closeSignal <- " have been disconnected due to idleness"
			close(closeSignal)
		case <-closeSignal:
			timer.Stop()
		}

		fmt.Println("3 - END", who)
	}()

	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			timer.Reset(timeoutDuration)
			messages <- who + ": " + input.Text()
		}

		select {
		case <-closeSignal:
		default:
			closeSignal <- " has left"
			close(closeSignal)
		}

		// note ignoring potential errors from input.Err()
		fmt.Println("4 - END", who)
	}()

	leaveMsg := <-closeSignal
	leaving <- cli
	messages <- who + leaveMsg
	conn.Close()
	fmt.Println(who + " disconnected")

	fmt.Println("1 - END", who)
}

func clientWriter(conn net.Conn, ch chan string, who string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ingoring network errors
	}

	fmt.Println("2 - END", who)
}
