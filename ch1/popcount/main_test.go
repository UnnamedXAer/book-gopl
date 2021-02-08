package main

import "testing"

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCount(uint64(i))
	}
}
func BenchmarkPopCountLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PopCountLoop(uint64(i))
	}
}
func TestPopCount(t *testing.T) {
	for i := 0; i < 100; i++ {
		if PopCount(uint64(i)) != PopCountLoop(uint64(i)) {
			t.Errorf("expected results to be equal got %d, %d, for i = %d", PopCount(uint64(i)), PopCountLoop(uint64(i)), i)
		}
	}
}
