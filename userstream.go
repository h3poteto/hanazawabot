package main

import(
	"fmt"
	"os"
	"net/url"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"./kanachan"
	"./dbyoutube"
	"./dbserif"

)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))

	values := url.Values{}
	self, _ := api.GetSelf(values)
	stream := api.UserStream(values)
	for {
		event := <-stream.C
		switch event.(type) {
		case anaconda.Tweet:
			tweet := event.(anaconda.Tweet)

			ydb := &dbyoutube.DBYoutubeMovie{}
			var youtubedb dbyoutube.YoutubeMovie = ydb


			aKana := &kanachan.Kanachan{}
			var kana kanachan.Kana = aKana
			if kana.IncludeCheck(tweet.Text) {
				// DB保存
			}
			// RTチェック
			// TODO: リプライチェック
			if checkRetweet(tweet, self) {
			} else if checkReply(tweet, self) {
				youtube_url, err := youtubedb.SelectRandom()
				if err != "" {
					fmt.Printf("DBYoutube random select error: %v\n", err)
					continue
				}
				tweet_value := url.Values{}
				tweet_value.Set("in_reply_to_status_id", tweet.IdStr)

				sdb := &dbserif.DBSerif{}
				var serif dbserif.Serif = sdb
				tweet_serif, err := serif.SelectRandom()
				if err != "" {
					fmt.Printf("DBSerif random select error: %v\n", err)
					continue
				}

				_, error := api.PostTweet("@" + tweet.User.ScreenName + " " + tweet_serif + " " + youtube_url, tweet_value)
				if error != nil {
					fmt.Printf("twitter api error: %v\n", error)
				} else {
					fmt.Printf("reply to @%v: %v\n", tweet.User.ScreenName, tweet_serif)
				}

			} else {
				continue
			}
			// favチェック
		}
	}
}

func checkReply(tweet anaconda.Tweet, user anaconda.User) bool {
	if tweet.InReplyToUserID == user.Id {
		return true
	}
	if strings.Contains(tweet.Text, "@" + user.ScreenName) {
		return true
	}
	return false
}

func checkRetweet(tweet anaconda.Tweet, user anaconda.User) bool {
	if strings.Index(tweet.Text, "RT") == 0 && tweet.User.Id != user.Id && strings.Contains(tweet.Text, "@" + user.ScreenName) {
		return true
	} else {
		return false
	}
}
