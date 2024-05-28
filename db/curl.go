package db

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

var client = req.C().SetTimeout(time.Second * 4).SetRedirectPolicy(req.NoRedirectPolicy())

// 爬取url，失败重试 reties 次数
func queryUrl(url string, retries int) *goquery.Document {

	if retries < 0 {
		retries = 3
	}

	var doc *goquery.Document = nil

	for i := 0; i < retries; i++ {
		resp, err := client.R().
			SetHeader("User-Agent", "Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)").
			Get(url)
		if err == nil {
			doc, err = goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				doc = nil
			}
			break
		}

	}

	return doc
}
