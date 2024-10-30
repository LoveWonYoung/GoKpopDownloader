/*
 * @Author: LoveWonYoung leeseoimnida@gmail.com
 * @Date: 2024-07-19 09:35:43
 * @LastEditors: LoveWonYoung leeseoimnida@gmail.com
 * @LastEditTime: 2024-10-30 15:36:25
 * @FilePath: \GoKpopDownloader\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
            佛祖保佑       永无BUG
*/

package main

import (
	"encoding/json"
	"fmt"
	"karina/downloader"
	"regexp"
	"sync"
	"time"

	//	"time"
)

type urlJson struct {
	Items        map[string]string `json:"items"`
	Max          int               `json:"max"`
	AppendTarget string            `json:"appendTarget"`
	Content      string            `json:"content"`
	ReplaceHTML  string            `json:"replaceHtml"`
}

var (
	idol      = "Wonyoung"
	totalPage = 1

	firstWg     sync.WaitGroup
	findallPick sync.WaitGroup
)

// 找到一页的所有链接
func firstRespStr(page int) []string {
	var retImageList []string
	c := downloader.RequestsTextP(fmt.Sprintf("https://kpopping.com/profiles/idol/%v/latest-pictures/%v", idol, page))
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

// 并发找到所有页面的所有链接
func allRespList() ([]string, int) {
	firstWg.Add(totalPage)
	var allLink []string
	for i := range totalPage {
		go func(p int) {
			defer firstWg.Done()
			allLink = append(allLink, firstRespStr(p+1)...)
		}(i)
	}
	firstWg.Wait()
	return allLink, len(allLink)
}

// 请求一个以上的链接返回一个页面下所有图片的链接
func respOnePicLink(l string) []string {
	var onePageAllPicLink []string
	r := downloader.RequestsTextG(l) // 返回的是一个html网页，直接用正则得到所有图片的链接
	re := regexp.MustCompile(`<a href="/documents/(.*?)" data`).FindAllStringSubmatch(r, -1)
	for _, i := range re {
		onePageAllPicLink = append(onePageAllPicLink, "https://kpopping.com/documents/"+i[1]) // 拼接这些路径
		fmt.Println("onePageAllPicLink", i[1])
	}
	fmt.Println("onePageAllPicLink", onePageAllPicLink)
	return onePageAllPicLink
}

//从上面返回的链接中找到所有图片的下载链接
func allPicList() []string {
	var res, count = allRespList()
	var allPicLink []string
	findallPick.Add(count)
	for _, l := range res {
		go func(s string) {
			defer findallPick.Done()
			allPicLink = append(allPicLink, respOnePicLink(s)...)
		}(l)
	}
	findallPick.Wait()
	return allPicLink
}
func main() {
	start := time.Now()
	fmt.Println(allPicList())
	end := time.Since(start)
	fmt.Println(end)
}
