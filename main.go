package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	url := "https://jsonplaceholder.typicode.com/posts"
	var wg sync.WaitGroup
	countPost := 100
	wg.Add(countPost)
	for i := 1; i <= countPost; i++ {
		go getPost(&url, i, &wg)
	}
	wg.Wait()
}

func getPost(purl *string, postId int, pwg *sync.WaitGroup) {

	defer pwg.Done()
	url := *purl
	url = url + "/" + strconv.Itoa(postId)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	sbt, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s  \n", sbt)

}
