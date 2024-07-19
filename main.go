package main

import (
	"encoding/json"
	"fmt"
	"karina/downloader"
	"regexp"
	"strings"
	"sync"
)

type urlJson struct {
	Items        map[string]string `json:"items"`
	Max          int               `json:"max"`
	AppendTarget string            `json:"appendTarget"`
	Content      string            `json:"content"`
	ReplaceHTML  string            `json:"replaceHtml"`
}

var (
	wg     sync.WaitGroup
	wgReq  sync.WaitGroup
	wgPage sync.WaitGroup
)

func kpopPageList(page int) []string {
	var retImageList []string
	c := downloader.RequestsTextP(fmt.Sprintf("https://kpopping.com/profiles/idol/%v/latest-pictures/%v", downloader.IDOLNAME, page))
	var html urlJson
	err := json.Unmarshal([]byte(c), &html)
	if err != nil {
		fmt.Println(err)
	}

	// 定义一个正则表达式
	re := regexp.MustCompile(`<a href="(.*?)" class="cell" aria-label="album">`)
	imageList := re.FindAllStringSubmatch(html.Content, -1)
	for _, image := range imageList {
		retImageList = append(retImageList, "https://kpopping.com"+image[1])
	}
	return retImageList
}

// 从以上返回的链接列表中再次请求，找到每个元素指向所有图片的链接 并发太快了 io跟不上，老是报错
func OneLinkAllPic() {
	// var allPicLink []string
	for page := range 5 {
		wgPage.Add(1)
		go func(p int) {
			for _, l := range kpopPageList(p) {
				wgReq.Add(1)
				go func(reqLink string) {
					tempReq := downloader.RequestsTextP(reqLink)
					// 用正则表达式从返回的html中找到所有图片的地址
					re := regexp.MustCompile(`<a href="/documents/(.*?)" data`)
					//多线程下载所有图片
					for _, oneLink := range re.FindAllStringSubmatch(tempReq, -1) {
						wg.Add(1)
						go func(name, link string) {
							downloader.DownloadImage(name, link)
							defer wg.Done()
						}(strings.Split(strings.Split(oneLink[1], "/")[3], ".")[0], "https://kpopping.com/documents/"+oneLink[1])
					}
					wg.Wait()
					defer wgReq.Done()
				}(l)
			}
			wgReq.Wait()
			defer wgPage.Done()
		}(page)
	}
	wgPage.Wait()
}

// 从以上返回的链接列表中再次请求，找到每个元素指向所有图片的链接
func OneLinkAllPic2() {
	// var allPicLink []string
	for page := range 2 {
		for _, l := range kpopPageList(page) {
			wgReq.Add(1)
			go func(reqLink string) {
				tempReq := downloader.RequestsTextP(reqLink)
				// 用正则表达式从返回的html中找到所有图片的地址
				re := regexp.MustCompile(`<a href="/documents/(.*?)" data`)
				//多线程下载所有图片
				for _, oneLink := range re.FindAllStringSubmatch(tempReq, -1) {
					wg.Add(1)
					go func(name, link string) {
						downloader.DownloadImage(name, link)
						defer wg.Done()
					}(strings.Split(strings.Split(oneLink[1], "/")[3], ".")[0], "https://kpopping.com/documents/"+oneLink[1])
				}
				wg.Wait()
				defer wgReq.Done()
			}(l)
		}
		wgReq.Wait()
		defer wgPage.Done()
	}
}
func main() {
	OneLinkAllPic2()
}
