package main

import "testing"

var dirs = []string{
	"../../../.",
}

func BenchmarkReadFileSizes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readFileSizes(dirs, false)
	}
}
