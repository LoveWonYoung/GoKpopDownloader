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

func picturesDownload() {
	var wgDownloadPage sync.WaitGroup
	wgDownloadPage.Add(downloader.DOWNLOADPAGE)
	// 并发请求所有page页面
	for page := 1; page < downloader.DOWNLOADPAGE+1; page++ {
		go func(p int) {
			var kpl sync.WaitGroup
			k := kpopPageList(p) // 请求一页，返回12个的图片链接
			kpl.Add(len(k))
			for _, i := range k {
				go func(image_link_12 string) {
					// fmt.Println(image_link_12)
					tempReq := downloader.RequestsTextP(image_link_12) // 对12个链接其中之一进行请求 ，找到每个链接下所有的图片链接
					re := regexp.MustCompile(`<a href="/documents/(.*?)" data`)
					everyPicLink := re.FindAllStringSubmatch(tempReq, -1) // 这是每个连接下所有图片的地址
					var wgEveryPic sync.WaitGroup
					wgEveryPic.Add(len(everyPicLink))
					// fmt.Println("这页链接有", len(everyPicLink), "张图片")
					// 开始下载
					for _,dl:=range everyPicLink{
						go func(n,d string){
							downloader.DownloadImage(n, d)
							defer wgEveryPic.Done()
						}(strings.Split(strings.Split(dl[1], "/")[3], ".")[0], "https://kpopping.com/documents/"+dl[1])
					}
					wgEveryPic.Wait()
					defer kpl.Done() // 没请求一次就done一次
				}(i)
			}
			kpl.Wait()

			defer wgDownloadPage.Done()
		}(page)

	}
	wgDownloadPage.Wait()
}
func main() {
	picturesDownload()
	fmt.Println("下载完成")
}
