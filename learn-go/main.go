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
	}

	// 버퍼드 채널 생성
	resultChan := make(chan struct{ url, status string }, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			status := "OK"
			if err := hitURL(url); err != nil {
				status = "FAILED"
			}
			resultChan <- struct{ url, status string }{url, status}
		}(url)
	}

	// 모든 Go 루틴 완료 대기
	wg.Wait()
	// 채널 닫기
	close(resultChan)

	// 메인 Go 루틴에서 결과 수집
	results := make(map[string]string)
	for res := range resultChan {
		// <-resultChan으로 받은 데이터를 results 맵에 저장
		results[res.url] = res.status
	}

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
