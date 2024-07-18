package main

import (
	"encoding/json"
	"fmt"
	"karina/downloader"
	"regexp"
)

type urlJson struct {
	Items        map[string]string `json:"items"`
	Max          int               `json:"max"`
	AppendTarget string            `json:"appendTarget"`
	Content      string            `json:"content"`
	ReplaceHTML  string            `json:"replaceHtml"`
}

func kpopPageList() []string {
	var retImageList []string
	c := downloader.RequestsTextP("https://kpopping.com/profiles/idol/Wonyoung/latest-pictures/1")
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
func main() {
	for _, url := range kpopPageList() {
		fmt.Println(url)
	}
}
