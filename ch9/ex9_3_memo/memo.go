package memo

import (
	"fmt"
)

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
	done     Done
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
type Func func(key string, done Done) (interface{}, error)

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
			go e.call(f, req.key, req.done) // call f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string, done Done) {
	// Evaluate the function
	e.res.value, e.res.err = f(key, done)
	// Broadcast the ready condition.
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
	memo.requests <- request{key: key, response: response, done: done}
	select {
	case <-done:
		return nil, fmt.Errorf("memo.Get: call of the function for the %q key was cancelled", key)
	case res := <-response:
		return res.value, res.err
	}
}
