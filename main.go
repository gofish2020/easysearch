package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofish2020/easysearch/conf"
	"github.com/gofish2020/easysearch/db"
	"github.com/gofish2020/easysearch/jieba"

	"github.com/gin-gonic/gin"
)

func main() {

	//1. 初始化配置信息
	conf.InitConf()

	//2. 初始化数据库
	db.InitDB()

	//3. 初始化分词
	jieba.InitJieba("/Users/mac/source/easysearch/jieba/dict")

	flag.Parse()

	args := os.Args

	if len(args) > 1 {
		switch strings.ToLower(args[1]) { // ./easysearch spider 10
		case "spider": //爬取连接
			limit := 100 // 限制一次处理的url数量
			if len(args) > 2 {
				v, _ := strconv.Atoi(string(args[2]))
				if v > 0 {
					limit = v
				}
			}
			db.DoSpider(limit)
		case "dict": // ./easysearch dict 2
			//分词
			limit := 2
			if len(args) > 2 {
				v, _ := strconv.Atoi(string(args[2]))
				if v > 0 {
					limit = v
				}
			}
			db.DoDict(limit)

			// c := cron.New(cron.WithSeconds())
			// c.AddFunc("25 * * * * *", f)

			// go c.Start()

			//select {}

		}
		return
	}

	//启动web服务，搜索结果

	e := gin.Default()
	e.GET("/", func(ctx *gin.Context) {
		fmt.Println(ctx.Params)
		ctx.JSON(200, gin.H{"message": "ok"})
	})
	e.Run(":8080")
}
