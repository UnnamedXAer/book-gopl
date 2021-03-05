package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// type file struct {
// 	size int64
// 	name string
// }

// func (f file) String() string {
// 	return fmt.Sprintf("% 15d - %s", f.size, f.name)
// }

func main() {
	showProgress := flag.Bool("v", false, "show verbose progress message")
	flag.Parse()
	readFileSizes(flag.Args(), *showProgress)
}

func readFileSizes(dirs []string, showProgress bool) {
	fmt.Println("showProgress", showProgress)
	fileSizes := make(chan int64)
	// done := make(chan struct{})

	var sizeSum int64
	var fileCount int64

	wg := sync.WaitGroup{}
	go func() {
		for _, d := range dirs {
			wg.Add(1)
			go walkDir(&wg, d, fileSizes)
		}
		wg.Wait()
		close(fileSizes)
	}()

	var ticker <-chan time.Time
	if showProgress {
		ticker = time.Tick(300 * time.Millisecond)
	}
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if ok == false {
				break loop
			}
			fileCount++
			sizeSum += size
		case <-ticker:
			printDiskUsage(fileCount, sizeSum)
		}
	}

	fmt.Println("\nTotal:")
	printDiskUsage(fileCount, sizeSum)
}

func printDiskUsage(fileCount, sizeSum int64) {
	fmt.Printf("\rprogress: % 10d files, % 15d bytes, % 5.2f GB\n", fileCount, sizeSum, float64(sizeSum)/1024/1024/1024)
}

// walkDir recursively walks the file three rooted at dir
// and sends the size of each found file on file sizes
func walkDir(wg *sync.WaitGroup, dir string, fileSizes chan<- int64) {
	defer wg.Done()

	for _, f := range dirents(dir) {
		if f.IsDir() {
			wg.Add(1)
			go walkDir(wg, filepath.Join(dir, f.Name()), fileSizes)
			continue
		}
		fileSizes <- f.Size()
	}
}

// dirents returns the entires of the directory
func dirents(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	return files
}
