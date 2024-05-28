package controllers

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gofish2020/easysearch/db"
	"github.com/gofish2020/easysearch/jieba"
)

func Search(c *gin.Context) {
	// 开始时间
	t := time.Now()

	keyword := c.Query("keyword")
	fmt.Println(keyword)

	result := []db.SearchResult{}

	if utf8.RuneCountInString(keyword) > 0 {
		// 对搜索词进行分词
		words := jieba.JiebaInstance.CutForSearch(keyword, true)
		result = db.Search(words)
	}
	latency := time.Since(t)

	c.HTML(200, "search.tpl", gin.H{
		"title":   "Go搜搜Go",
		"time":    time.Now().Format("2006-01-02 15:04:05"), // 当前时间
		"values":  result,                                   // 结果列表
		"keyword": keyword,
		"latency": latency, // 耗时
	})

}
