package main

import (
	"testing"
	"time"
)

var args = []string{
	"https://www.linkedin.com/",
	"https://stackoverflow.com/",
	"",
	"broken",
	"https://translate.google.com",
	"https://golang.org/",
	"https://www.google.pl/",
}

func TestDoWork(t *testing.T) {
	for i := 0; i < 100; i++ {
		go func() {
			t := 1 + i*10
			if t > 400 {
				t = 3 * i
			}
			time.Sleep(200 * time.Millisecond)
			// fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ 0.2s timeout")
			close(done)
		}()
		// not very good test, but its mostly to test the occurrence
		// of panics triggered by cancellation
		doWork(args)
		<-done
		done = make(chan struct{})
	}
}
