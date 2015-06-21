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
	"./dbtweet"
	"./dbuser"
	"./dbretweet"
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

			if isMyself(tweet, self) {
				continue
			}
			ydb := &dbyoutube.DBYoutubeMovie{}
			var youtubedb dbyoutube.YoutubeMovie = ydb


			aKana := &kanachan.Kanachan{}
			var kana kanachan.Kana = aKana
			if kana.IncludeCheck(tweet.Text) {
				// contain
				tdb := &dbtweet.DBTweet{}
				var tweetdb dbtweet.Tweet = tdb

				udb := &dbuser.DBUser{}
				var userdb dbuser.User = udb

				user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

				if user.Id != 0 {
					_ = tweetdb.Add(user.Id, tweet.Text, tweet.Id)
				}
			}

			if isRetweet(tweet, self) {
				// retweet
				movie := youtubedb.ScanYoutubeMovie(tweet)

				rdb := &dbretweet.DBRetweet{}
				var retweetdb dbretweet.Retweet = rdb

				udb := &dbuser.DBUser{}
				var userdb dbuser.User = udb

				user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

				if user.Id != 0 {
					_ = retweetdb.Add(user.Id, movie.Id)
				}

			} else if isReply(tweet, self) {
				// reply
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
		case anaconda.EventTweet:
			event_tweet := event.(anaconda.EventTweet)
			if event_tweet.Event.Event == "favorite" {
				// TODO: favの処理
			}
		}
	}
}

func isReply(tweet anaconda.Tweet, user anaconda.User) bool {
	if tweet.InReplyToUserID == user.Id {
		return true
	}
	if strings.Contains(tweet.Text, "@" + user.ScreenName) {
		return true
	}
	return false
}

func isRetweet(tweet anaconda.Tweet, user anaconda.User) bool {
	if strings.Index(tweet.Text, "RT") == 0 && tweet.User.Id != user.Id && strings.Contains(tweet.Text, "@" + user.ScreenName) {
		return true
	} else {
		return false
	}
}

func isMyself(tweet anaconda.Tweet, user anaconda.User) bool {
	if tweet.User.Id == user.Id {
		return true
	} else {
		return false
	}
}
