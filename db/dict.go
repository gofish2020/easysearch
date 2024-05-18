package db

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gofish2020/easysearch/jieba"
	"golang.org/x/text/width"
	"gorm.io/gorm"
)

var maxId uint = 0

type wordInfo struct {
	count     int      // 词频
	positions []string // 词在文档中的位置
}

var once = sync.Once{}

var blackWord = make(map[string]struct{})

// craw_done = 1 and dict_done = 0 分词多少条记录
func DoDict(limit int) {

	// 1.使用 spider 库
	mysqlDB.Exec("use " + SpiderDBName)

	// 黑名单词汇
	once.Do(func() {
		blacks := []BlackList{}
		mysqlDB.Model(&blackList).Order("id asc").Find(&blacks)
		for _, black := range blacks {
			blackWord[black.Word] = struct{}{}
		}

	})
	// 2.查询
	var srcs []Source
	mysqlDB.Model(&Source{}).Where("dict_done = ? and craw_done=1 and id > ?", 1, maxId).Limit(limit).Order("id asc").Find(&srcs)

	batchUpdate := make(map[string]string)

	for _, src := range srcs {

		fmt.Println("文档id:", src.ID)
		htmlLen := utf8.RuneCountInString(src.Html)               // 页面字符数
		words := jieba.JiebaInstance.CutForSearch(src.Html, true) // 页面分词

		// 比如：我爱你你爱我，分词后就是：我/爱/你/你/爱/我 -->> 分词结果就是有2个我2个你2个爱
		pos := 0
		statWork := make(map[string]wordInfo) // 一个网页中，每个词汇出现的频次以及位置
		for _, word := range words {

			// 不区分大小写 比如 baidu 和 BAIDU是 一个词汇
			newWord := strings.ToLower(width.Narrow.String(word)) // 转半角
			wordLen := utf8.RuneCountInString(newWord)

			// 判断是否为黑名单词汇
			if _, ok := blackWord[newWord]; ok {
				pos += wordLen
				continue
			}

			if winfo, ok := statWork[newWord]; !ok {
				statWork[newWord] = wordInfo{
					count:     1,
					positions: []string{strconv.Itoa(pos)},
				}
			} else {
				winfo.count++
				winfo.positions = append(winfo.positions, strconv.Itoa(pos))
				statWork[newWord] = winfo
			}

			pos += wordLen
		}

		// 汇总结果
		for word, wordInfo := range statWork {

			// 文档id,文档长度,词频,位置1,位置2-
			infoString := strconv.FormatUint(uint64(src.ID), 10) + "," + strconv.Itoa(htmlLen) + "," + strconv.Itoa(wordInfo.count) + "," + strings.Join(wordInfo.positions, ",") + "-"

			batchUpdate[word] += infoString

		}

		// 更新t_source表分词状态
		src.DictDone = "1"
		src.UpdatedAt = time.Now()
		mysqlDB.Model(&src).Select("dict_done", "updated_at").Updates(src)

		if src.ID > maxId {
			maxId = src.ID
		}

	}

	// 批量更新 t_dict
	mysqlDB.Connection(func(tx *gorm.DB) error {
		tx.Exec("use " + DictDBName)

		tx = tx.Begin()
		for word, positions := range batchUpdate {
			tx.Exec(" insert into "+dictTable.TableName()+" (word,positions) values (?,?) on duplicate key update positions = concat(ifnull(positions,''), ?) ", word, positions, positions)
		}
		tx.Commit()
		return nil
	})

}
