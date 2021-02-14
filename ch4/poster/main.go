package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var apiKey string

type Released time.Time

func (r *Released) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	// s := string(b)
	// fmt.Println("Released s:", s)
	t, err := time.Parse("02 Jan 2006", s)
	if err != nil {
		return fmt.Errorf("unable to parse time %q, error: %v", s, err)
	}
	*r = Released(t)
	return nil
}

func (r Released) MarshalJSON() ([]byte, error) {
	// var s string
	return []byte{}, fmt.Errorf("not implemented (r Released) MarshalJSON()")
}

func (r *Released) String() string {
	// var s string
	return time.Time(*r).String()
}

type movieData struct {
	Title    string
	Year     string
	Released Released
	Poster   string
}

const omdbapiUrL = "https://omdbapi.com/?"

var nFlag = flag.String("n", "", "movie name")
var yFlag = flag.String("y", "", "movie year")

func main() {
	flag.Parse()
	getPoster()

}

func getPoster() {

	err := checkParams()
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := getData()
	if err != nil {
		fmt.Println("Unable to retrive poster, error:", err)
		return
	}

	if data.Title == "" {
		fmt.Println("No movies found for given params")
		return
	}

	if data.Poster == "" {
		fmt.Printf("Movie %q, released at %s was found, but there is no poster\n", data.Title, data.Released.String())
		fmt.Println("if it's not the movie you were looking for check correctness of your input")
		if *yFlag == "" {
			fmt.Println("\nyou may also pass a year of the movie using '-y' parameter")
		}
		return
	}

	fmt.Println("poster: \n", data.Poster)

	b, err := getPosterImgData(data.Poster)
	if err != nil {
		fmt.Printf("Unable to retrive poster, url: %q error: %v\n", data.Poster, err)
		return
	}

	ext := filepath.Ext(data.Poster)
	if ext == "" {
		ext = ".png"
	}

	err = ioutil.WriteFile(data.Title+ext, b, 0644)
	if err != nil {
		fmt.Println("Unable to save poster, error: ", err)
		return
	}
}

func checkParams() error {
	fmt.Println("apikey:", apiKey, `os.Getenv("APIKEY"):`, os.Getenv("APIKEY"))
	if apiKey == "" {
		apiKey = os.Getenv("APIKEY")
		if apiKey == "" {
			panic("'env' missing APIKEY")
		}
	}

	if *nFlag == "" && *yFlag == "" {
		return fmt.Errorf("please pass at least movie name like -n \"Bad Boys\"")
	}
	if *nFlag == "" {
		return fmt.Errorf("please pass movie name like '-n \"<movie name>\"' eg. -n \"Bad Boys\"")
	}
	return nil
}

func getData() (*movieData, error) {
	encoded := url.PathEscape(fmt.Sprintf("apikey=%s&t=%s&y=%s", apiKey, *nFlag, *yFlag))
	url := omdbapiUrL + encoded

	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	m := &movieData{}
	err = json.Unmarshal(body, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func getPosterImgData(posterURL string) ([]byte, error) {
	res, err := http.Get(posterURL + "&apikey=" + apiKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	return ioutil.ReadAll(res.Body)
}
