package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	storageFile = "comics.json"
	api         = "https://xkcd.com/%num%/info.0.json"
)

type notFoundError struct {
}

func (e *notFoundError) Error() string {
	return "not_found"
}

type comic struct {
	Year       int `json:"year,string"`
	Month      int `json:"month,string"`
	Day        int `json:"day,string"`
	Num        int
	SafeTitle  string `json:"safe_title"`
	Title      string
	Link       string
	Transcript string
	Img        string
	Alt        string
}

var titleFlag *string = flag.String("t", "", "pass comic name")
var numFlag *int = flag.Int("n", 0, "pass comic number")

func main() {
	flag.Parse()
	makeDB()
	if err := checkFlags(); err != nil {
		fmt.Println(err)
		return
	}

	printComic()
}

func printComic() {
	c, err := findComicInStorage()

	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Printf("The comic was found by unable to pretty print. Simple print:\n%v\n", c)
		return
	}
	fmt.Printf("\nYour comic is:\n%s", b)
}

func findComicInStorage() (*comic, error) {
	comics, err := getStorageData()
	if err != nil {
		return nil, err
	}

	if *numFlag != 0 {
		comic := findComicByNum(comics, *numFlag)
		if comic != nil {
			return comic, nil
		}
	}

	if *titleFlag != "" {
		comic := findComicByTitle(comics, *titleFlag)
		return comic, nil
	}
	return nil, fmt.Errorf("there is no comic matching to your criteria")
}

func findComicByNum(comics []comic, num int) *comic {
	for i := len(comics) - 1; i >= 0; i-- {
		if comics[i].Num == num {
			return &comics[i]
		}
	}
	return nil
}

func findComicByTitle(comics []comic, title string) *comic {
	lcTitle := strings.ToLower(title)
	for i := len(comics) - 1; i >= 0; i-- {
		if strings.HasPrefix(strings.ToLower(comics[i].Title), lcTitle) {
			return &comics[i]
		}
	}
	return nil
}

func getStorageData() ([]comic, error) {
	f, err := os.Open(storageFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	comics := make([]comic, 2000, 4000)
	in := bufio.NewScanner(f)
	var c comic
	var b []byte
	var n int = -1
	for in.Scan() {
		n++
		b = in.Bytes()
		err = json.Unmarshal(b, &c)
		if err != nil {
			fmt.Printf("line %d, unable to unmarshal data: %q, error: %v\n", n, b, err)
			continue
		}
		comics = append(comics, c)
	}

	return comics, nil
}

func checkFlags() error {

	if (*titleFlag == "") && (*numFlag == 0) {
		return fmt.Errorf("pass comic title prefix or comic number > 0")
	}
	if (*titleFlag != "") && (*numFlag != 0) {
		fmt.Println("Both title and comic number is passed, program will prioritize comic number.")
	}
	return nil
}

func makeDB() {
	existed := ensureFileExists(storageFile)
	if existed == false {
		fmt.Println("There was not comics storage, the comics have to be downloaded. This may take a few moments...\n ")
		downloadComics()
	}
}

func downloadComics() {
	failedComics := []int{}
	var n int = 1 // the 0 indexed comic does not exists
	f, err := os.OpenFile(storageFile, os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err404Cnt := 0
	for {
		// if n > 5 {
		// 	fmt.Println("you should now have 5 comics in the storage file")
		// 	return
		// }
		b, err := fetchComic(n)
		if err != nil {
			if errors.Is(err, &notFoundError{}) {
				fmt.Printf("comic with num %q not found\n", n)
				err404Cnt++
				if err404Cnt > 5 {
					fmt.Print("5 comics not found in a row, programs stops futher downloading\n")
					break
				}
			}
			failedComics = append(failedComics, n)
			fmt.Printf("comic number %q added to retires due to error: %q\n", n, err)
			n++
			continue
		}
		err404Cnt = 0
		err = saveComic(f, b)
		if err != nil {
			failedComics = append(failedComics, n)
		}
		n++
	}

}

// fetchComic fetches comic with given number
func fetchComic(n int) ([]byte, error) {
	res, err := http.Get(strings.Replace(api, "%num%", strconv.Itoa(n), 1))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, &notFoundError{}
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n%d.\n%s\n", n, b)

	res.Body.Close()
	return b, nil
}

// saveComic appends comic data to storage file
func saveComic(f *os.File, b []byte) error {
	b = append(b, '\n')
	_, err := f.Write(b)
	return err
}

// ensureFileExists checks if file at given path exists
// if not creates it.
//
// Returns `true` if the file already existed, `false` when the file had to be created.
//
// It will panic if en error occurs different then `ErrNotExist`.
//
// @todo: rewrite it to use OpenFile and check if first line contains some kind of mark
// @todo: that indicates the data was already fetched.
func ensureFileExists(filePath string) (existed bool) {
	f, err := os.Open(storageFile)
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("fail to close file %q, error: %q\n", filePath, err)
			return
		}
	}()
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(storageFile)
			if err != nil {
				panic(err)
			}
			return false
		}
		panic(err)
	}
	return true
}
