package kanachan

import (
	"testing"
	"."
)

func TestExceptCheck(t *testing.T) {
	aKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = aKana

	if kana.ExceptCheck("花澤香菜 歌ってみた") != false {
		t.Error("歌ってみたを含むワードが除外されない")
	}

	if kana.ExceptCheck("花澤香菜踊ってみた") != false {
		t.Error("踊ってみたを含むワードが除外されない")
	}
}

func TestIncludeCheck(t *testing.T) {
	aKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = aKana

	if kana.IncludeCheck("撫子ちゃん，花澤さん") != true {
		t.Error("花澤というワードに反応していない")
	}

	if kana.IncludeCheck("撫子ちゃん，花澤香菜") != true {
		t.Error("花澤香菜というワードに反応していない")
	}

	if kana.IncludeCheck("撫子ちゃん，花澤 香菜") != true {
		t.Error("花澤 香菜というワードに反応していない")
	}

	if kana.IncludeCheck("撫子ちゃん，花澤　香菜") != true {
		t.Error("花澤　香菜というワードに反応していない")
	}

	if kana.IncludeCheck("なでこちゃん") != false {
		t.Error("関係ないワードに反応している")
	}
}

func TestBuildTweet(t *testing.T) {
	aKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = aKana

	if kana.BuildTweet("", "台詞", "タイトル", "https://youtube.com/") != "台詞 【タイトル】https://youtube.com/" {
		t.Error("リプライなしのツイート文字列生成が上手く行ってない")
	}

	if kana.BuildTweet("@h3_poteto ", "台詞", "タイトル", "https://youtube.com/") != "@h3_poteto 台詞 【タイトル】https://youtube.com/" {
		t.Error("リプライのツイート文字列生成が上手く行ってない")
	}

	if kana.BuildTweet("@h3_poteto ", "あたしもう25歳だよ！お前あたしと結婚できんのか！", "【責任取れ！】花澤香菜「あたしもう25歳だよ！お前あたしと結婚できんのか！」豊永利行「責任取ればいいの？///」杉田智和の好き勝手な発言に花澤香菜が絶叫www「花澤さんがおもいっきり吹いたこの紙笛はこちらの宛先まで♪」www", "https://youtube.com/") != "@h3_poteto あたしもう25歳だよ！お前あたしと結婚できんのか！ 【【責任取れ！】花澤香菜「あたしもう25歳だよ！お前あたしと結婚できんのか！」豊永利行「責任取ればいいの？///」杉田智和の好き勝手な発言に花澤香菜が絶叫…】https://youtube.com/" {
		t.Error("140字を超えるツイートを短縮できていない")
	}
}
