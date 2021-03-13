package main

import (
	"testing"
)

func BenchmarkPipline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		in, _ := pipline(2 * 1000000)
		in <- uint64(0)
	}
}
