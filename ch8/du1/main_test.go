package main

import "testing"

var dirs = []string{
	"../../../.",
	"d:/flarrow/",
}

func BenchmarkReadFileSizes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readFileSizes(dirs)
	}
}
