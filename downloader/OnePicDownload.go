package downloader

import (
	"fmt"
	"io"
	"net/http"
)

func RequestsTextP(url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	return string(body)
}
func RequestsTextG(url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	return string(body)
}
func OnePicDownload(picName, link string) {
	fmt.Println("this is OnePicDownload")
	picsText := RequestsTextG(link)
}
