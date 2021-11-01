package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	client := http.Client{}
	resp, err := client.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s \n", body)
}
