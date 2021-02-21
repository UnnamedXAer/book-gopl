package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := WaitForServer(os.Args[1])
	if err != nil {
		log.Println()
		log.Println()
		log.Fatalf("Server is down: %v\n", err)
	}

	log.Printf("server %q can be connected", os.Args[1])
}

// WaitForServer attempts to contact to the server of URL.
// It tries for one minute using exponential back-off
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	const timeout = time.Minute * 1
	deadline := time.Now().Add(timeout)
	log.SetPrefix("wait: ")
	f := log.Flags()
	log.SetFlags(0)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success, we can connect to server
		}
		log.Printf("try#%d, sever not responding (%s}; retrying...\n", tries+1, err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	log.SetFlags(f)
	log.SetPrefix("")

	return fmt.Errorf("server %q failed to respond after %s", url, timeout)
}
