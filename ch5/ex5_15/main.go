package main

import (
	"fmt"
	"io/ioutil"
)

type CrawlDir func(path string) []string

func main() {
	dir := `/`
	if dir == `/` || dir == `\` {
		dir = "." + dir
	}
	breadthFirst(dir, crawlDir)
}

func breadthFirst(path string, f CrawlDir) {

	worklist := []string{path}
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, dir := range items {
			for _, subDir := range f(dir) {
				worklist = append(worklist, dir+"/"+subDir)
			}
		}
	}

}

var crawlDir CrawlDir = func(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	folders := make([]string, 0, len(files))
	for _, f := range files {
		isDir := f.IsDir()
		printDir(isDir, f.Name())
		if isDir {
			folders = append(folders, f.Name())
		}
	}
	return folders
}

func printDir(isDir bool, dir string) {
	dirChar := ""
	if isDir {
		dirChar = "_/"
	}
	fmt.Printf("% 2s%s\n", dirChar, dir)
}
