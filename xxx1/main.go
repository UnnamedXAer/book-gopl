package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {

	p := "./web/pages"
	err := os.MkdirAll(p, 0644)
	if err != nil {
		log.Fatalln("cannot create path", err)
	}
	u := strings.ReplaceAll(url.PathEscape("https://linuxize.com/post/what-does-chmod-777-mean/"), ":", "")
	_, err = os.OpenFile(p+"/"+u+".file", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			log.Fatalln("file already exists, great")
		}
		log.Fatalln(err)
	}
	fmt.Println("file created, also great")
}
