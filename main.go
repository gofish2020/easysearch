package main

import (
	"flag"
	"os"
	"strings"

	"github.com/gofish2020/easysearch/conf"
	"github.com/gofish2020/easysearch/controllers"
	"github.com/gofish2020/easysearch/db"
	"github.com/gofish2020/easysearch/jieba"

	"github.com/gin-gonic/gin"
)

func main() {

	flag.Parse()
	args := os.Args
	//1. 初始化配置信息
	conf.InitConf()

	//2. 初始化数据库
	db.InitDB()

	//3. 初始化分词
	jieba.InitJieba("./jieba/dict") // path.Join(filepath.Dir(args[0]), "jieba/dict")

	if len(args) > 1 {
		switch strings.ToLower(args[1]) { // ./easysearch spider
		case "spider": //爬取连接
			for i := 0; i < 5; i++ {
				db.DoSpider(10)
			}
		case "dict": // ./easysearch dict
			//分词
			for i := 0; i < 5; i++ {
				db.DoDict(10)
			}

			// 用定时任务一直执行
			// c := cron.New(cron.WithSeconds())
			// c.AddFunc("25 * * * * *", f)
			// go c.Start()
			//select {}

		}
		return
	}

	//启动web服务，搜索结果

	e := gin.Default()
	e.LoadHTMLGlob("views/*")

	e.GET("/", controllers.Search)
	e.Run(":8080")
}
