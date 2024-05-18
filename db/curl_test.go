package db

import (
	"testing"
)

func TestCrul(t *testing.T) {
	doc := queryUrl("https://tgdb.37.com/?uid=hao123dt", 3)

	t.Log(doc.Html())
	//assert.Equal(t, nil, doc)
}
