package kanachan

import (
	"strings"
)

var (
	except_words = [...]string{"歌ってみた", "踊ってみた"}
)

type Kana interface {
	ExceptCheck(string) bool
}

type Kanachan struct {
}

func (u *Kanachan) ExceptCheck(word string) bool {
	for _, except := range except_words {
		if strings.Contains(word, except) {
			return false
		}
	}
	return true
}
