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

type InputField struct {
	Url, Dir, Pattern string
	ch                chan []byte
	SaveKey           bool // true - In file save key and value
}

func main() {

	ch := make(chan []byte)
	inputfild := InputField{

		"https://jsonplaceholder.typicode.com/posts",
		"storage/posts/",
		"*.txt",
		ch,
		true,
	}

	countPost := 10

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	err := os.MkdirAll(inputfild.Dir, 0777)
	if err != nil {
		errorLog.Fatal(err)
	}
	var wg sync.WaitGroup

	for i := 1; i <= countPost; i++ {
		wg.Add(2)
		go inputfild.getPostToChanel(i, &wg, errorLog)
		go inputfild.savePostToFile(&wg, errorLog)
	}
	wg.Wait()

}

func (i *InputField) getPostToChanel(count int, pwg *sync.WaitGroup, errorLog *log.Logger) {

	defer pwg.Done()
	url := i.Url
	url = url + "/" + strconv.Itoa(count)

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

	i.ch <- bresp

}
func (i *InputField) savePostToFile(pwg *sync.WaitGroup, errorLog *log.Logger) {

	defer pwg.Done()
	ch := <-i.ch

	post, err := postCreateFromJson(ch)
	if err != nil {
		errorLog.Print(err)
		return
	}

	strid := strconv.Itoa(post.Id)
	fileName := strings.Replace(i.Pattern, "*", strid, 1)
	if i.SaveKey {
		err = ioutil.WriteFile(fileName, ch, 0777)
		if err != nil {
			errorLog.Print(err)

		}
		return
	}
	bpost := []byte(fmt.Sprint(post))
	err = ioutil.WriteFile(fileName, bpost, 0777)
	if err != nil {
		errorLog.Print(err)

	}

}
