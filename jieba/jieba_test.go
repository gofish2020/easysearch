package jieba

import (
	"testing"
	"unicode/utf8"

	"golang.org/x/text/width"
)

func TestRune(t *testing.T) {

	s := "你好呀，我好呀。Baidu" // 这里的 标点符号是全角

	t.Log(utf8.RuneCountInString(s))
	t.Log(s)

	// 将符号转成半角
	word := width.Narrow.String(s)
	t.Log(utf8.RuneCountInString(word))
	t.Log(word)

	InitJieba("/Users/mac/source/easysearch/jieba/dict")
	t.Log(JiebaInstance.CutForSearch(s, true))
	t.Log(JiebaInstance.CutForSearch(word, true))

	point := width.Narrow.String("，")
	t.Log(point)
	if point == "," {
		t.Log("same")
	}

}
