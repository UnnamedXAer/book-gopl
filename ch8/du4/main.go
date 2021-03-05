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

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func sendCancel() {
	// input := bufio.NewScanner(os.Stdin)
	// for input.Scan() {
	// 	close(done)
	// 	return
	// }

	os.Stdin.Read(make([]byte, 1)) // read just one bute
	close(done)
}

func main() {
	showProgress := flag.Bool("v", false, "show verbose progress message")
	flag.Parse()
	go sendCancel()
	readFileSizes(flag.Args(), *showProgress)
}

func readFileSizes(dirs []string, showProgress bool) {
	fmt.Println("START")
	fileSizes := make(chan int64)

	var sizeSum int64
	var fileCount int64

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 20)

	go func() {
		for _, d := range dirs {
			wg.Add(1)
			go walkDir(&wg, semaphore, d, fileSizes)
		}
		wg.Wait()
		close(fileSizes)
	}()

	var ticker <-chan time.Time
	// when thicker is nil it blocks forever therefore in select statement
	// under nil chanel case no code will execute
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
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish.
			for range fileSizes {
				// do nothing...
			}
			return
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
func walkDir(wg *sync.WaitGroup, semaphore chan struct{}, dir string, fileSizes chan<- int64) {
	defer wg.Done()
	if cancelled() {
		return
	}
	for _, f := range dirents(dir, semaphore) {
		if f.IsDir() {
			wg.Add(1)
			go walkDir(wg, semaphore, filepath.Join(dir, f.Name()), fileSizes)
			continue
		}
		fileSizes <- f.Size()
	}
}

// dirents returns the entires of the directory
func dirents(dir string, sema chan struct{}) []os.FileInfo {
	select {
	case sema <- struct{}{}: // acquire token
		defer func() {
			<-sema
		}()
	case <-done:
		return nil // cancelled
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return nil
	}
	return files
}
