package urlchecker

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrRequestFailed = errors.New("request failed")

func HitURL(url string) error {
	fmt.Println("Checking:", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode >= 400 {
		return ErrRequestFailed
	}
	defer resp.Body.Close()
	return nil
}
