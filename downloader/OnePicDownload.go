package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func RequestsTextP(url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	//fmt.Println(resp.StatusCode, "请求代码")
	if resp.StatusCode != 200 {
		return ""
	}
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
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
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
func DownloadImage(imgaeName string, downloadLink string) {

	req, err := http.Get(downloadLink)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := req.Body.Close()
		if err != nil {
			fmt.Println("文件关闭错误")
		}
	}()
	content, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("下载文件：", imgaeName+`.jpg`)
	err = os.WriteFile(DownLoadPath+imgaeName+`.jpg`, content, 0666)
	if err != nil {
		panic(err)
	}
}
func OnePicDownload(picName, link string) {
	fmt.Println("this is OnePicDownload")
	// picsText := RequestsTextG(link)
}
