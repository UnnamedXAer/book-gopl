package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type employee struct {
	name      string
	inserted  bool
	createtAt time.Time
	pic       []byte
	picname   string
}

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	b := make([]byte, 3)
	username := "Ann"
	e := employee{
		name: username,
		pic:  b,
	}
	go createCustomer(wg, &e)
	go updatePic(wg, &e)
	wg.Wait()
	fmt.Println("\n", e)
	// time.Sleep(1 * time.Second)
	fmt.Printf("All done.\n")
}

func createCustomer(wg *sync.WaitGroup, e *employee) {
	defer wg.Done()
	fmt.Printf("inserting customer %q to db...\n", e.name)
	time.Sleep(time.Millisecond * 10)
	e.inserted = true
	t := time.Now()
	fmt.Println(t)

	fmt.Printf("customer %q inserted to db\n", e.name)
}

func updatePic(wg *sync.WaitGroup, e *employee) {
	defer wg.Done()
	fmt.Printf("updating pic for the %q user...\n", e.name)
	time.Sleep(time.Millisecond * 2)
	e.picname = e.name + strconv.Itoa(time.Now().Second())
	fmt.Printf("%q pic for %q updated, %d bytes written\n", e.picname, e.name, len(e.pic))
}
