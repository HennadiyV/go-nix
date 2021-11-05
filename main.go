package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Post struct {
	UserId int
	Id     int
	Title  string
	Body   string
}

func postCreateFromJson(b []byte) (p Post, err error) {

	err = json.Unmarshal(b, &p)
	if err != nil {

		return p, err
	}
	return p, err

}

func main() {

	url := "https://jsonplaceholder.typicode.com/posts"
	dir := "storage/posts/"
	pattern := "*.txt"
	countPost := 100

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		errorLog.Fatal(err)
	}
	var wg sync.WaitGroup
	ch := make(chan []byte)

	for i := 1; i <= countPost; i++ {
		wg.Add(2)
		go getPostToChanel(ch, &url, i, &wg, errorLog)
		go savePostToFile(ch, dir+pattern, &wg, errorLog)
	}
	wg.Wait()

}

func getPostToChanel(c chan []byte, purl *string, postId int, pwg *sync.WaitGroup, errorLog *log.Logger) {

	defer pwg.Done()
	url := *purl
	url = url + "/" + strconv.Itoa(postId)

	resp, err := http.Get(url)

	if err != nil {
		errorLog.Fatal(err)
		return
	}
	defer resp.Body.Close()

	bresp, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLog.Fatal(err)
		return
	}

	c <- bresp

}
func savePostToFile(c chan []byte, pattern string, pwg *sync.WaitGroup, errorLog *log.Logger) {

	defer pwg.Done()

	post, err := postCreateFromJson(<-c)
	if err != nil {
		errorLog.Print(err)
		return
	}
	bpost := []byte(fmt.Sprint(post))
	strid := strconv.Itoa(post.Id)
	fileName := strings.Replace(pattern, "*", strid, 1)

	err = ioutil.WriteFile(fileName, bpost, 0777)
	if err != nil {
		errorLog.Print(err)
		return
	}

}
