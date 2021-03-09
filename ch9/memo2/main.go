// Package memo provides a concurrency-unsafe(!)
// memoization of a function of type Func.
package memo

import (
	"fmt"
	"sync"
)

// Memo caches th results of calling a Func.
type Memo struct {
	f     Func
	mu    sync.Mutex
	cache map[string]result
}

// Func is the type og the function to memoize.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

// New returns a pointer to the new Memo object
func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Get returns a results from Memo f function for given key
// either by calling it or from the cache.
// concurrency-safe
func (memo *Memo) Get(key string) (interface{}, error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	if ok == false {
		fmt.Println("calling: ", key)
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	memo.mu.Unlock()

	return res.value, res.err
}
