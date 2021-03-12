package memo

import (
	"fmt"
	"time"
)

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

type result struct {
	value interface{}
	err   error
}

// Memo caches th results of calling a Func.
type Memo struct {
	requests chan request
}

// Done allows to cancel f function
type Done chan struct{}

// Func is the type og the function to memoize.
type Func func(key string) (interface{}, error)

// New returns a memoization of f. Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (m *Memo) Close() { close(m.requests) }

func (m *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range m.requests {
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	fmt.Printf("calling: f(%q)\n", key)
	// Evaluate the function
	time.Sleep(4 * time.Second)
	e.res.value, e.res.err = f(key)
	fmt.Printf("f(%q) resolved\n", key)
	time.Sleep(4 * time.Second)
	// Broadcast the ready condition.
	fmt.Printf("closing e.ready for f(%q)\n", key)
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready
	// Send the result to the client.
	response <- e.res
}

// Get returns a results from Memo f function for given key
// either by calling it or from the cache.
// concurrency-safe
func (memo *Memo) Get(key string, done Done) (interface{}, error) {
	response := make(chan result)
	// var res result
	memo.requests <- request{key: key, response: response}
	select {
	case <-done:
		// fmt.Printf("function for key: %q was cancelled.\n", key)
		return nil, fmt.Errorf("memo.Get: call of the function for the %q key was cancelled", key)
	case res := <-response:
		return res.value, res.err
	}
}
