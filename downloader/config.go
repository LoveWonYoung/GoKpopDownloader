package downloader

import "fmt"

// https://kpopping.com/profiles/idol/Yujin3
const (
	IDOLNAME = `Karina2`
	DOWNLOADPAGE = 2
)

var (
	DownLoadPath = fmt.Sprintf("/root/myGoCode/myfile/GoKpopDownloader/result/%v/", IDOLNAME)
)
