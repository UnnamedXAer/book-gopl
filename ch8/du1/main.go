package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// type file struct {
// 	size int64
// 	name string
// }

// func (f file) String() string {
// 	return fmt.Sprintf("% 15d - %s", f.size, f.name)
// }

func main() {
	flag.Parse()
	readFileSizes(flag.Args())
}

func readFileSizes(dirs []string) {
	fileSizes := make(chan int64)

	var sizeSum int64
	var fileCount int64

	go func() {
		for _, d := range dirs {
			walkDir(d, fileSizes)
		}
		close(fileSizes)
	}()

	for n := range fileSizes {
		fileCount++
		sizeSum += n
	}

	fmt.Println("\nTotal:")
	printDiscUsage(fileCount, sizeSum)
}

func printDiscUsage(fileCount, sizeSum int64) {
	fmt.Printf("\rprogress: % 10d files, % 15d bytes, % 5.2f GB\n", fileCount, sizeSum, float64(sizeSum)/1024/1024/1024)
}

// walkDir recursively walks the file three rooted at dir
// and sends the size of each found file on file sizes
func walkDir(dir string, fileSizes chan<- int64) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			walkDir(filepath.Join(dir, f.Name()), fileSizes)
			continue
		}
		fileSizes <- f.Size()
	}
}
