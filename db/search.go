package db

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type StatResult struct {
	Count int64 `gorm:"column:count"`
	Total int64 `grom:"column:total"`
}

/*
词1 ： 文档1 文档2 文档3
词2 ： 文档2 文档3

先利用词1，计算出一个IDF。然后 词1 和 文档1 文档2 文档3挨个计算一个值RQiDj * IDF，结果作为 文档1 文档2 文档3 的得分

然后继续 词2:计算出一个IDF。然后 词2 和 文档2 文档3 挨个计算一个值RQiDj *IDF， 结果作为 文档2 文档3的得分

可以看到都有文档2/文档3，所以他们的积分需要累加起来
*/
func Search(words []string) []SearchResult {

	statRes := StatResult{}
	mysqlDB.Exec("use " + SpiderDBName)
	mysqlDB.Raw("select count(1) as count, sum(html_len) as total from "+srcTable.TableName()+" where dict_done=?", 1).Find(&statRes)

	// N 已经完成分词的文档数量
	N := statRes.Count
	if N == 0 { // do nothing
		return nil
	}

	// 平均文档长度
	avgDocLength := float64(statRes.Total) / float64(statRes.Count)

	// 记录文档的总得分
	docScore := make(map[string]float64)

	mysqlDB.Exec("use " + DictDBName)

	// 每次针对一个词进行处理
	for _, word := range words {

		var dict Dict
		mysqlDB.Model(&dictTable).Where("word = ?", word).Find(&dict)
		if dict.ID == 0 {
			continue
		}

		positions := strings.Split(dict.Positions, "-")
		positions = positions[:len(positions)-1] // 去掉最后一个【符号-】

		// NQi 含有这个词的文档有多少个
		NQi := int64(len(positions)) // -1因为尾部多了一个【符号-】
		// IDF ： 如果一个词在很多页面里面都出现了，那说明这个词不重要，例如“的”字，哪个页面都有，说明这个词不准确，进而它就不重要。
		IDF := math.Log10((float64(N-NQi) + 0.5) / (float64(NQi) + 0.5))

		// 固定值
		k1 := 2.0
		b := 0.75
		for _, position := range positions { // 这里其实就是和 word 关联的 所有文档

			docInfo := strings.Split(position, ",") //  文档id,文档长度,词频,位置1,位置2-

			//Dj :文档长度
			Dj, _ := strconv.Atoi(docInfo[1])
			//Fi :该词在文档中出现的词频
			Fi, _ := strconv.Atoi(docInfo[2])

			// 本词和本文档的相关性
			RQiDj := (float64(Fi) * (k1 + 1)) / (float64(Fi) + k1*(1-b+b*(float64(Dj)/avgDocLength)))

			// 汇总所有文档的相关性总得分
			docScore[docInfo[0]] += RQiDj * IDF
		}
	}

	// 根据积分排序文档

	docIds := make([]string, 0, len(docScore)) // 文档ids
	for id := range docScore {
		docIds = append(docIds, id)
	}

	sort.SliceStable(docIds, func(i, j int) bool {
		return docScore[docIds[i]] > docScore[docIds[j]] // 按照分从【大到小排序】
	})

	// 默认取 20个
	limit := 20
	if len(docIds) < limit {
		limit = len(docIds)
	}
	docIds = docIds[0:limit]

	// 然后利用id回表查询文档内容
	sousouResult := make([]SearchResult, 0)
	mysqlDB.Exec("use " + SpiderDBName)
	for _, docId := range docIds {

		var src Source
		mysqlDB.Table(src.TableName()).Where("id = ?", docId).Find(&src)

		// 从Html中截取 一段文本作为简介（去掉了其中的ascii编码字符）
		rep := regexp.MustCompile("[[:ascii:]]")
		brief := rep.ReplaceAllLiteralString(src.Html, "")

		briefLen := utf8.RuneCountInString(brief)
		if briefLen > 100 {
			brief = string([]rune(brief)[:100])
		}

		sousouResult = append(sousouResult, SearchResult{
			Title: src.Title,
			Brief: brief,
			Url:   src.Url,
		})
	}

	return sousouResult
}

type SearchResult struct {
	Title string
	Brief string
	Url   string
}
