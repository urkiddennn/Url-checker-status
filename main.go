package main

import (
	"fmt"
	"net/http"
	"sync"
)

type URL struct {
	url       string
	urlStatus int
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.github.com",
		"http://www.youtube.com",
		"https://www.asdefrfsfeasdwdwa.com", // return a status not ok
	}

	urlChan := make(chan string) // send a data
	getStatus := make(chan URL)
	numberWorkers := 3

	var wg sync.WaitGroup

	for i := 1; i <= numberWorkers; i++ {
		wg.Add(1)
		go checker(i, urlChan, getStatus, &wg)
	}

	go func() {
		for result := range getStatus {
			fmt.Printf("URL %s, Status: %d\n", result.url, result.urlStatus)
		}
	}()

	go func() {
		for _, url := range urls {
			urlChan <- url
		}
		close(urlChan)
	}()

	wg.Wait() // wait for the routines to get done
	close(getStatus)
}

func checker(id int, urlChan <-chan string, urlGetChan chan<- URL, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range urlChan {
		fmt.Println("Testing url", id, url)
		resp, err := http.Get(url)
		status := 0
		if err != nil {
			fmt.Println("Error or invalid URL", id, url, err)
			return
		} else {
			status = resp.StatusCode
			resp.Body.Close()
		}
		urlGetChan <- URL{url: url, urlStatus: status}
	}
}
