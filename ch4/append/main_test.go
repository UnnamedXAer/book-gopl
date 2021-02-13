package main

import "testing"

const n = 40

func BenchmarkAppendInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := []int{}
		for j := 0; j < n; j++ {
			x = appendInt(x, j)
		}
	}
}

func BenchmarkAppendIntBookVer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := []int{}
		for j := 0; j < n; j++ {
			x = appendIntBookVer(x, j)
		}
	}
}
