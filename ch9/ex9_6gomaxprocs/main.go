package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/unnamedxaer/book-gopl/ch8/thumbnail"
)

func main() {
	fnames, err := getFileNames("assets")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(fnames)
	filenames := make(chan string, len(fnames))
	go func() {
		for _, f := range fnames {
			filenames <- f
		}
		close(filenames)
	}()
	t := time.Now()
	totalSize := makeThumbnails5(filenames)
	log.Println("done", time.Now().Sub(t))
	log.Println(totalSize)
}

// makeThumbnails5 makes thumbnails for each file received from the chanel.
// It returns the number of bytes occupied by the files it creates.
func makeThumbnails5(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines

	for f := range filenames {
		wg.Add(1)

		// worker
		go func(f string) {
			defer wg.Done()
			thumbfile, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumbfile) // ok to ignore error

			sizes <- info.Size()

		}(f)

	}

	//closer
	go func() {
		wg.Wait()
		close(sizes)
	}()
	var total int64
	for size := range sizes {
		total += size
	}
	return total
}

// // makeThumbnails4 makes thumbnails for the specified files in parallel.
// // It returns the generated file names in arbitrary order,
// // or an error if any step failed.
// func makeThumbnails4(filenames []string) (thumbfiles []string, err error) {
// 	type item struct {
// 		thumbfile string
// 		err       error
// 	}

// 	ch := make(chan item, len(filenames))
// 	for _, f := range filenames {
// 		go func(f string) {
// 			var it item
// 			it.thumbfile, it.err = thumbnail.ImageFile(f)
// 			ch <- it
// 		}(f)
// 	}

// 	for range filenames {
// 		it := <-ch
// 		if it.err != nil {
// 			return nil, it.err
// 		}
// 		thumbfiles = append(thumbfiles, it.thumbfile)
// 	}
// 	return thumbfiles, nil
// }

// // makeThumbnails3 makes thumbnails for the specified files in parallel.
// // It returns an error if any step failed.
// func makeThumbnails3(filenames []string) error {
// 	errors := make(chan error)
// 	for _, f := range filenames {
// 		go func(f string) {
// 			_, err := thumbnail.ImageFile(f)
// 			errors <- err
// 		}(f)
// 	}

// 	for range filenames {
// 		if err := <-errors; err != nil {
// 			return err // NOTE: incorrect: goroutine leak!
// 		}
// 	}
// 	return nil
// }

// func makeThumbnails2(filenames []string) {
// 	ch := make(chan struct{})
// 	for _, f := range filenames {
// 		go func(f string) {
// 			if _, err := thumbnail.ImageFile(f); err != nil {
// 				log.Println(err)
// 			}
// 			ch <- struct{}{}
// 		}(f)
// 	}

// 	for range filenames {
// 		<-ch
// 	}
// }

// func makeThumbnails(filenames []string) {
// 	for _, f := range filenames {
// 		if _, err := thumbnail.ImageFile(f); err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func getFileNames(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fns := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.Contains(f.Name(), "thumb") {
			continue
		}
		fns = append(fns, filepath.Join(dir, f.Name()))
	}
	return fns, err
}
