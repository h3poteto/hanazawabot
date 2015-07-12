package kanachan

import (
	"strings"
	"unicode/utf8"
)

var (
	except_words = [...]string{"歌ってみた", "踊ってみた"}
	include_words = [...]string{"花澤 香菜", "花澤香菜", "花澤"}
)

type Kana interface {
	ExceptCheck(string) bool
	IncludeCheck(string) bool
	BuildTweet(string, string, string, string) string
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

func (u *Kanachan) BuildTweet(reply string, serif string, title string, url string) string {
	tweet := reply + serif + " " + "【" + title + "】"
	// twitterはURLを短縮するので，URLは必ず23文字以下になる
	// http://tech.aainc.co.jp/archives/6472
	over_count := utf8.RuneCountInString(tweet) - (140 - 25)
	if over_count > 0 {
		title_length := utf8.RuneCountInString(title)
		r := []rune(title)
		new_title := string(r[0:(title_length - over_count)])
		tweet = reply + serif + " " + "【" + new_title + "…】" + url
	} else {
		tweet += url
	}
	return tweet
}
