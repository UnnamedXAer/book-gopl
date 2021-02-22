package main

import (
	"fmt"
	"log"
)

type event struct {
	name        string
	listeneters []func(args ...string)
}

func newEvent(name string) *event {
	return &event{
		name: name,
	}
}

// EventEmitter emits events
type EventEmitter struct {
	events map[string]*event
}

func (ee EventEmitter) findOrCreateEvent(event string) *event {
	if _, ok := ee.events[event]; ok == false {
		ee.events[event] = newEvent(event)
	}
	return ee.events[event]
}

func (ee EventEmitter) on(event string, listener func(args ...string)) {
	e := ee.findOrCreateEvent(event)
	e.listeneters = append(e.listeneters, listener)
}

func (ee EventEmitter) emit(event string, args ...string) {
	if e, ok := ee.events[event]; ok {
		for _, listener := range e.listeneters {
			listener(args...)
		}
	}
}

type server struct {
	connection    string
	connected     bool
	eventsEmitter EventEmitter
}

func (s server) init() {
	s.eventsEmitter.on("connect", handleConnection)
	s.eventsEmitter.on("connect", logMessage)
	s.eventsEmitter.on("data", logMessage)
	s.eventsEmitter.on("disconnect", handleDisconnection)
	s.eventsEmitter.on("disconnect", logMessage)
}

func handleConnection(args ...string) {
	fmt.Println(args)
}

func handleDisconnection(args ...string) {
	fmt.Println(args[0])
}

func logMessage(args ...string) {
	log.Print(args[1])
}

// Function for connecting to server
func (s *server) connect(connectionString string) {
	s.connection = connectionString
	s.connected = true
	s.eventsEmitter.emit("connect", s.connection, "Connected successfully")
}

// Getting data from users
func (s *server) data(userData string) {
	if s.connected {
		s.eventsEmitter.emit("data", s.connection, userData)
	}
}

// Function for disconnecting from server
func (s *server) disconnect() {
	if s.connected {
		s.connected = false
		s.eventsEmitter.emit("disconnect", s.connection, "Disconnected successfully")
	} else {
		fmt.Print("Already disconnected")
	}
}

func main() {
	var s = server{
		eventsEmitter: EventEmitter{
			events: make(map[string]*event),
		},
	}
	s.init()

	s.connect("our/connection/string")

	s.data("Message from user 1")
	s.data("Message from user 2")
	s.data("Message from user 3")
	s.data("Message from user 4")

	s.disconnect()
}
