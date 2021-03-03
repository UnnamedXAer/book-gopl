package main

import "testing"

func BenchmarkDoImg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doImg()
	}
}

func BenchmarkDoImgParaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doImgParaller()
	}
}

func BenchmarkDoImgWG(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doImgWG()
	}
}
