package main

import (
	"errors"
	"fmt"
	"net/http"
)

type RequsetResult struct {
	url    string
	status string
}

var errRequestFailed = errors.New("request failed")

func main() {
	results := map[string]string{}
	c := make(chan RequsetResult)
	urls := []string{"https://www.google.com", "https://www.bing.com", "https://www.yahoo.com", "https://www.duckduckgo.com", "https://www.baidu.com", "https://www.yandex.com", "https://www.ask.com", "https://www.aol.com", "https://www.wolframalpha.com", "https://www.ecosia.org", "https://www.startpage.com"}

	for _, url := range urls {
		go hitURL(url, c)
	}

	for range urls {
		result := <-c
		results[result.url] = result.status
	}

	for url, status := range results {
		fmt.Printf("%s: %s\n", url, status)
	}
}

func hitURL(url string, c chan<- RequsetResult) {
	fmt.Println("Hitting URL:", url)
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- RequsetResult{url: url, status: status}

}
