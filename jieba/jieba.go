package jieba

import (
	"path"

	"github.com/yanyiwu/gojieba"
)

var JiebaInstance *gojieba.Jieba

func InitJieba(dictDir string) {

	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")
	JiebaInstance = gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)

}
