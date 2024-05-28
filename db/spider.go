package db

import (
	"errors"
	"fmt"
	"net/url"
	"time"
	"unicode/utf8"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofish2020/easysearch/tools"
	"golang.org/x/text/width"
	"gorm.io/gorm"
)

// craw_done = 0  从数据库一次获取limit个url，然后进行处理
func DoSpider(limit int) {
	//1. 指定使用spider数据库
	mysqlDB.Exec("use " + SpiderDBName)

	//2. 执行t_source表查询
	srcs := []Source{}
	//mysqlDB.Raw(" select id,url,md5 from t_source where craw_done=0 order by id asc limit 0,?", limit).Find(&res)
	mysqlDB.Model(&Source{}).Select([]string{"id,url,host,md5"}).Where("craw_done = ?", 0).Limit(limit).Order("id asc").Find(&srcs)

	for _, src := range srcs {
		//爬取信息
		doc := queryUrl(src.Url, 3)
		if doc == nil { // 爬取失败
			src.CrawDone = "2"
			if src.Host == "" {
				u, _ := url.Parse(src.Url)
				src.Host = strings.ToLower(u.Host)
			}
			if src.Md5 == "" {
				src.Md5 = tools.GetMD5Hash(src.Url)
			}
			src.HtmlLen = 0
			src.UpdatedAt = time.Now()
			mysqlDB.Model(&src).Select("craw_done", "host", "md5", "html", "html_len", "title", "updated_at").Updates(src)
		} else {

			src.CrawDone = "1"
			src.Html = tools.StringStrip(strings.TrimSpace(doc.Text()))
			src.HtmlLen = uint(utf8.RuneCountInString(src.Html))
			src.Title = tools.StringStrip(strings.TrimSpace(doc.Find("title").Text()))
			src.UpdatedAt = time.Now()
			if src.Md5 == "" {
				src.Md5 = tools.GetMD5Hash(src.Url)
			}
			if src.Host == "" {
				u, _ := url.Parse(src.Url)
				src.Host = strings.ToLower(u.Host)
			}

			// 爬取成功,更新数据库
			res := mysqlDB.Model(&src).Select("craw_done", "html", "html_len", "title", "update_at", "md5", "host").Updates(src)
			if res.RowsAffected > 0 {
				fmt.Println("已成功爬取", src.Url)
			}

			uniqueUrl := make(map[string]struct{}) // 当前页面中的a链接去重
			// 解析html，提取更多的a链接，保存到数据库
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				//
				href := width.Narrow.String(strings.Trim(s.AttrOr("href", ""), " \n"))
				_url, _, _ := strings.Cut(href, "#") // 去掉 #后面的
				_url = strings.ToLower(_url)

				if _, ok := uniqueUrl[_url]; ok { // 重复
					return
				}
				uniqueUrl[_url] = struct{}{}

				if tools.IsUrl(_url) { // 有效 url

					// 这里还需要做一次全局去重复( 如果数量大，可以用redis来记录url是否记录过)
					md5 := tools.GetMD5Hash(_url)
					u, _ := url.Parse(_url)

					src := Source{Md5: md5, Url: _url, Host: strings.ToLower(u.Host), Model: gorm.Model{CreatedAt: time.Now(), UpdatedAt: time.Now()}}

					result := mysqlDB.First(&src, "md5 = ?", md5)

					if errors.Is(result.Error, gorm.ErrRecordNotFound) { // 说明不存在
						// 保存到数据库中
						result := mysqlDB.Create(&src)

						if result.Error == nil && result.RowsAffected > 0 {
							fmt.Println("插入成功", _url)
						}
					}
				}
			})
		}

	}
}
