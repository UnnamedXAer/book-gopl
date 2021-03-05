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

func main() {
	showProgress := flag.Bool("v", false, "show verbose progress message")
	flag.Parse()
	readFileSizes(flag.Args(), *showProgress)
}

type dirInfo struct {
	dir   string
	count int64
	size  int64
}

func readFileSizes(dirs []string, showProgress bool) {
	fmt.Println("showProgress", showProgress)
	totals := make([]*dirInfo, len(dirs))

	wg := sync.WaitGroup{}
	wg.Add(len(dirs))
	for i, d := range dirs {
		go func(i int, d string) {
			totals[i] = readFileSizesInRootDir(d, showProgress)
			wg.Done()
		}(i, d)
	}

	wg.Wait()
	for _, total := range totals {
		fmt.Printf("\nTotal of %q is:\n%s", total.dir, printDiskUsage(total.count, total.size))

	}
}

func readFileSizesInRootDir(dir string, showProgress bool) *dirInfo {
	fileSizes := make(chan int64)

	var sizeSum int64
	var fileCount int64

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 8)

	go func() {
		wg.Add(1)
		go walkDir(&wg, semaphore, dir, dir, fileSizes)
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
		case n, ok := <-fileSizes:
			if ok == false {
				break loop
			}
			fileCount++
			sizeSum += n
		case <-ticker:
			fmt.Printf("\r%s\n%s:", dir, printDiskUsage(fileCount, sizeSum))

		}
	}
	return &dirInfo{
		dir:   dir,
		count: fileCount,
		size:  sizeSum,
	}
}

func printDiskUsage(fileCount, sizeSum int64) string {
	return fmt.Sprintf("\r% 10d files, % 15d bytes, % 5.2f GB\n",
		fileCount,
		sizeSum,
		float64(sizeSum)/1024/1024/1024)
}

// walkDir recursively walks the file three rooted at dir
// and sends the size of each found file on file sizes
func walkDir(
	wg *sync.WaitGroup,
	semaphore chan struct{},
	baseDir,
	dir string,
	fileSizes chan<- int64) {

	defer wg.Done()
	semaphore <- struct{}{}
	for _, f := range dirents(dir) {
		if f.IsDir() {
			wg.Add(1)
			go walkDir(wg, semaphore, baseDir, filepath.Join(dir, f.Name()), fileSizes)
			continue
		}
		fileSizes <- f.Size()
	}
	<-semaphore
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
