package main

import (
	"fmt"
	"runtime"
	"testing"
)

func Benchmark(b *testing.B) {
	runtime.GOMAXPROCS(12)
	fnames, _ := getFileNames("assets")
	filenames := make(chan string, len(fnames)*b.N)
	go func() {
		for i := 0; i < b.N; i++ {
			for _, f := range fnames {
				filenames <- f
			}
		}
		close(filenames)
	}()
	totalSize := makeThumbnails5(filenames)
	fmt.Println("total:", totalSize)
	fmt.Println()
}
