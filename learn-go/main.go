package main

import (
	"errors"
	"fmt"
	"net/http"
)

var errRequestFailed = errors.New("request failed")

func main() {

	urls := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.airbnb.com",
		"https://www.amazon.com",
		"https://www.reddit.com",
	}

	results := make(map[string]string)

	for _, url := range urls {
		result := "OK"
		err := hitURL(url)
		if err != nil {
			result = "FAILED"
		}
		results[url] = result
	}

}

func hitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return errRequestFailed
	}

	defer resp.Body.Close()

	return nil
}
