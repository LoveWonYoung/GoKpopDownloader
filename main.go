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
	"io"
	"karina/downloader"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
	//	"time"
)

var (
	path      = "Picture"
	idol      = "Wonyoung"
	totalPage = 1

	firstWg     sync.WaitGroup
	findAllPick sync.WaitGroup
)

type urlJson struct {
	Items        map[string]string `json:"items"`
	Max          int               `json:"max"`
	AppendTarget string            `json:"appendTarget"`
	Content      string            `json:"content"`
	ReplaceHTML  string            `json:"replaceHtml"`
}

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
		//fmt.Println("onePageAllPicLink", i[1])
	}
	//fmt.Println("onePageAllPicLink", onePageAllPicLink)
	return onePageAllPicLink
}

func MakeDir(s string) {
	// 定义你想要检查或创建的文件夹路径
	dirName := fmt.Sprintf(".\\%v\\%v", path, s)

	// 使用Stat函数获取文件夹信息
	_, err := os.Stat(dirName)

	// 判断错误类型
	if os.IsNotExist(err) {
		// 如果文件夹不存在，则创建文件夹
		err := os.MkdirAll(dirName, 0755) // 0755是目录的权限
		if err != nil {
			fmt.Println("创建文件夹失败:", err)
			return
		}
		fmt.Println("文件夹已创建:", dirName)
	} else if err != nil {
		// 如果发生其他错误，打印错误信息
		fmt.Println("获取文件夹信息时出错:", err)
		return
	} else {
		// 如果文件夹存在，则打印文件夹信息
		//fmt.Printf("文件夹已存在，信息如下: %+v\n", info)
	}
}
func DownloadImage(imageName string, downloadLink string) {

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
	fmt.Println("下载文件：", imageName+`.jpg`)
	err = os.WriteFile(imageName+`.jpg`, content, 0666)
	if err != nil {
		panic(err)
	}
}
func getDownloadLink() {
	var res, _ = allRespList()
	for _, i := range res {
		// https://kpopping.com/kpics/241010-Wonyoung-Instagram-Update
		r := strings.Split(i, `/`)
		// 链接拆分，得到目录名字
		if reflect.DeepEqual(r, []string{}) {
			continue
		} else {
			dirName := r[len(r)-1]
			MakeDir(dirName)
			// 对链接发起请求，得到下载链接
			resOne := respOnePicLink(i)
			// 处理每一个下载链接，得到文件名
			findAllPick.Add(len(resOne))
			for _, down := range resOne {
				go func() {
					defer findAllPick.Done()
					firstName := strings.Split(down, ".")[1]
					imageName := strings.Split(firstName, "/")
					//路径拼接
					finalPath := `.\` + path + `\` + dirName + `\` + imageName[len(imageName)-1]
					DownloadImage(finalPath, down)
				}()
			}
			findAllPick.Wait()
		}

	}
}
func main() {
	start := time.Now()

	getDownloadLink()

	end := time.Since(start)
	fmt.Println(end)
}
