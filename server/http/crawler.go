package http

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func Crawler() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// req, err := http.NewRequest(http.MethodGet, "https://www.ptt.cc/bbs/Gossiping/index.html", nil)
	req, err := http.NewRequest(http.MethodGet, "https://tw.stock.yahoo.com/quote/2330.TW", nil)
	if err != nil {
		fmt.Printf("goCrawler err %v\n", err)
		return
	}
	// req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("goCrawler do err: %v\n", err)
		return
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Response Body: ", string(body))

	// doc, err := goquery.NewDocumentFromReader(resp.Body)
	// if err != nil {
	// 	fmt.Println("NewDocumentFromReader err:", err)
	// 	return
	// }

	// 提取文章標題
	// doc.Find("div.title a").Each(func(index int, item *goquery.Selection) {
	// 	title := item.Text()
	// 	fmt.Println(title)
	// })

	GetImg("http://image.baidu.com/search/index?tn=baiduimage&ps=1&ct=201326592&lm=-1&cl=2&nc=1&ie=utf-8&word=%E7%BE%8E%E5%A5%B3")
}

var reImg = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`

func GetImg(url string) {
	pageStr := GetPageStr(url)
	if pageStr == "" {
		return
	}
	fmt.Println("GetPageStr ", pageStr)

	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr, -1)
	fmt.Println("result size ", len(results))
	for _, result := range results {
		fmt.Println(result[0])
	}
}

// 抽取根据url获取内容
func GetPageStr(url string) (pageStr string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("GetPageStr err: %v\n", err)
		return
	}
	defer resp.Body.Close()
	// 2.读取页面内容
	pageBytes, _ := io.ReadAll(resp.Body)
	// 字节转字符串
	pageStr = string(pageBytes)
	return pageStr
}
