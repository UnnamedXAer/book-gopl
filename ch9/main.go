package main

import (
	"fmt"
	"runtime"
)

// func httpGetBody(url string) (interface{}, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	return ioutil.ReadAll(resp.Body)
// }

// func incomingURLs() []string {
// 	input := bufio.NewScanner(os.Stdin)
// 	s := []string{}
// 	for input.Scan() {
// 		s = append(s, input.Text())
// 	}

// 	return s
// }

// func main() {

// 	m := memo.New(httpGetBody)
// 	for _, url := range incomingURLs() {
// 		start := time.Now()
// 		value, err := m.Get("http://google.com")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		fmt.Printf("%s, %s, %d, bytes\n",
// 			url, time.Since(start), len(value.([]byte)))
// 	}
// }

func main() {
	fmt.Println(runtime.GOMAXPROCS(2))
	for {
		go fmt.Print("0")
		fmt.Print("1")
	}
}
