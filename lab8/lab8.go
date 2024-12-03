package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// 定義 -max flag，預設值為 10
	max := flag.Int("max", 10, "限制印出的留言數量 (預設為 10)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  -max <number>: 限制印出的留言數量 (預設為 10)")
	}
	flag.Parse()

	// 驗證 max 的值
	if *max <= 0 {
		log.Fatal("Error: -max 必須為正整數")
	}

	// 使用 Colly 建立 Collector
	c := colly.NewCollector()

	// 保存留言資料
	type Comment struct {
		Name    string
		Content string
		Time    string
	}
	var comments []Comment

	// 選擇留言的 HTML 區塊
	c.OnHTML(".push", func(e *colly.HTMLElement) {
		name := strings.TrimSpace(e.ChildText(".push-userid"))
		content := strings.TrimSpace(e.ChildText(".push-content"))
		time := strings.TrimSpace(e.ChildText(".push-ipdatetime"))

		// 檢查必要資料是否存在
		if name != "" && content != "" && time != "" {
			comments = append(comments, Comment{
				Name:    name,
				Content: strings.TrimPrefix(content, ": "), // 去掉留言內容前的 ": "
				Time:    time,
			})
		}
	})

	// 設定錯誤處理
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Request failed:", err)
	})

	// 爬取 PTT 頁面
	err := c.Visit("https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html")
	if err != nil {
		log.Fatal(err)
	}

	// 格式化輸出留言，根據 max 限制數量
	for i, comment := range comments {
		if i >= *max {
			break
		}
		fmt.Printf("%d. 名字：%s，留言: %s, 時間： %s\n", i+1, comment.Name, comment.Content, comment.Time)
	}
}
