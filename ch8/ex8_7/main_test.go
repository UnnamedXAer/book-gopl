package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

const link = "http://localhost:3091/bench"

func TestMain(m *testing.M) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		f, err := ioutil.ReadFile("test.html")
		if err != nil {
			log.Fatalln(err)
		}

		http.HandleFunc("/bench", func(w http.ResponseWriter, r *http.Request) {
			w.Write(f)
		})

		wg.Done()
		err = http.ListenAndServe(":3091", nil)
		if err != nil {
			log.Fatalln(err)
		}
	}()
	wg.Wait()

	time.Sleep(10 * time.Millisecond)
	os.Exit(m.Run())
}

func BenchmarkFetchAndSave(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fetchAndSave(link)
	}
}
