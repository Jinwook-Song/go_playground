package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
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

	// 결과를 저장할 맵
	results := make(map[string]string)
	// 맵 동기화를 위한 뮤텍스
	var mu sync.Mutex
	// Go 루틴 동기화를 위한 WaitGroup
	var wg sync.WaitGroup

	// 각 URL에 대해 Go 루틴 실행
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done() // Go 루틴 종료 시 WaitGroup 카운터 감소
			result := "OK"
			err := hitURL(url)
			if err != nil {
				result = "FAILED"
			}
			// 맵에 결과 쓰기 (뮤텍스로 동기화)
			mu.Lock()
			results[url] = result
			mu.Unlock()
		}(url) // URL을 매개변수로 전달
	}

	// 모든 Go 루틴이 완료될 때까지 대기
	wg.Wait()

	// 결과 출력
	for url, result := range results {
		fmt.Printf("%s: %s\n", url, result)
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
