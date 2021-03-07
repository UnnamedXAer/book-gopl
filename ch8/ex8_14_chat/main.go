package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp4", ":3030")
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
			for idx := range clients {
				cli.ch <- fmt.Sprintf("idx %v, v: %v", idx, idx.who)
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

var timeoutDuration = 15 * time.Second

// var cnt int
// var mu = sync.Mutex{}

// comments with text "...- END" are for learning purpose to ensure everything closes.

func handleConn(conn net.Conn) {
	cliTimeout := timeoutDuration

	ch := make(chan string)          // outgoing client messages
	closeSignal := make(chan string) // chanel to send close message
	timer := time.NewTimer(cliTimeout)
	var who string = conn.RemoteAddr().String() + " (temporary name)"
	fmt.Println(who + " connected")
	go clientWriter(conn, ch, &who)

	name := make(chan string)

	go askForName(conn, ch, name, timer, cliTimeout)
	select {
	case <-timer.C:
		go timeoutDisconnect(ch, closeSignal)
	case who = <-name:
		timer.Stop()
	}
	fmt.Printf("client with address %s, has set his name to %q", conn.RemoteAddr().String(), who)

	go setTimeout(ch, closeSignal, timer, cliTimeout, who)

	ch <- "You are: " + who
	messages <- who + " has arrived"
	cli := &client{
		ch:  ch,
		who: who,
	}
	entering <- cli

	go sendMessages(conn, closeSignal, timer, cliTimeout, who)

	leaveMsg := <-closeSignal
	leaving <- cli
	messages <- who + leaveMsg
	conn.Close()
	fmt.Println(who + " disconnected")
	fmt.Println("1 - END, handler", who)
}

// clientWriter sends messages to connected client
func clientWriter(conn net.Conn, ch chan string, who *string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ingoring network errors
	}

	fmt.Println("2 - END, clientWriter", *who)
}

// askForName prompts client for a name.
func askForName(
	conn net.Conn,
	cliCh chan<- string,
	name chan<- string,
	timer *time.Timer,
	timeout time.Duration) {
	cliCh <- "Enter your name: "
	input := bufio.NewScanner(conn)
	for input.Scan() {
		timer.Reset(timeout)
		txt := input.Text()
		if txt == "" {
			cliCh <- "name cannot be empty"
			continue
		}
		name <- txt
		fmt.Println("5 - END, askForName", txt, conn.RemoteAddr().String())
		return
	}
	close(name)
	// name <- "unknown"
	fmt.Println("5 - END, askForName", conn.RemoteAddr().String())
}

// setTimeout disconnects client if idle for `duration`
func setTimeout(ch, closeSignal chan string, timer *time.Timer, duration time.Duration, who string) {
	func() {
		timer.Reset(duration)
		select {
		case <-timer.C:
			timeoutDisconnect(ch, closeSignal)
		case <-closeSignal:
			timer.Stop()
		}

		fmt.Println("3 - END, setTimeout", who)
	}()
}

// timeoutDisconnect sends messages due to client idle timeout and closes closeSignal chanel
func timeoutDisconnect(ch, closeSignal chan<- string) {
	ch <- "You have been disconnected due to idleness"
	closeSignal <- " have been disconnected due to idleness"
	close(closeSignal)
}

// sendMessages sends client messages to the other connected clients
func sendMessages(
	conn net.Conn,
	closeSignal chan string,
	timer *time.Timer,
	timeout time.Duration,
	who string) {
	input := bufio.NewScanner(conn)
	for input.Scan() {
		timer.Reset(timeout)
		messages <- who + ": " + input.Text()
	}
	// note ignoring potential errors from input.Err()

	select {
	case <-closeSignal:
	default:
		closeSignal <- " has left"
		close(closeSignal)
	}

	fmt.Println("4 - END, sendMessages", who)
}
