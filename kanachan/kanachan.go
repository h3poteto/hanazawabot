package kanachan

import (
	"strings"
)

var (
	except_words = [...]string{"歌ってみた", "踊ってみた"}
	include_words = [...]string{"花澤 香菜", "花澤香菜", "花澤"}
)

type Kana interface {
	ExceptCheck(string) bool
	IncludeCheck(string) bool
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


func (u *Kanachan) IncludeCheck(word string) bool {
	for _, include := range include_words {
		if strings.Contains(word, include) {
			return true
		}
	}
	return false
}
