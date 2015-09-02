package main

import(
	"fmt"
	"os"
	"net/url"
	"strings"
	"io/ioutil"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"../kanachan"
	"../models/dbyoutube"
	"../models/dbserif"
	"../models/dbtweet"
	"../models/dbuser"
	"../models/dbretweet"
	"../models/dbfav"
)

func main() {
	preparePidFile()
	defer removePidFile()
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
			ydb := dbyoutube.NewDBYoutubeMovie()
			var youtubedb dbyoutube.YoutubeMovie = ydb


			aKana := &kanachan.Kanachan{}
			var kana kanachan.Kana = aKana
			if kana.IncludeCheck(tweet.Text) {
				// contain
				tdb := dbtweet.NewDBTweet()
				var tweetdb dbtweet.Tweet = tdb

				udb := dbuser.NewDBUser()
				var userdb dbuser.User = udb

				user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

				if user.Id != 0 {
					_ = tweetdb.Add(user.Id, tweet.Text, tweet.Id)
					fmt.Printf("Add Tweet Log: %v \n", tweet.Text)
				}
			}

			if isRetweet(tweet, self) {
				// retweet
				movie := youtubedb.ScanYoutubeMovie(tweet)
				if movie == nil {
					continue
				}
				rdb := &dbretweet.DBRetweet{}
				var retweetdb dbretweet.Retweet = rdb

				udb := dbuser.NewDBUser()
				var userdb dbuser.User = udb

				user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

				if user.Id != 0 {
					_ = retweetdb.Add(user.Id, movie.Id)
					fmt.Printf("Add Rewteet Log: %v \n", movie.Id)
				}

			} else if isReply(tweet, self) {
				// reply
				movie := youtubedb.SelectRandom()
				if movie == nil {
					fmt.Printf("DBYoutube random select error: %v\n")
					continue
				}
				tweet_value := url.Values{}
				tweet_value.Set("in_reply_to_status_id", tweet.IdStr)

				sdb := dbserif.NewDBSerif()
				var serif dbserif.Serif = sdb
				tweet_serif, err := serif.SelectRandom()
				if err != "" {
					fmt.Printf("DBSerif random select error: %v\n", err)
					continue
				}

				_, error := api.PostTweet(
					kana.BuildTweet("@" + tweet.User.ScreenName + " ",
						tweet_serif,
						movie.Title,
						movie.ConvertYoutubeID()),
					tweet_value)
				if error != nil {
					fmt.Printf("twitter api error: %v\n", error)
				} else {
					fmt.Printf("reply to @%v: %v\n", tweet.User.ScreenName, "@" + tweet.User.ScreenName)
				}

			} else {
				continue
			}
		case anaconda.EventTweet:
			event_tweet := event.(anaconda.EventTweet)
			if event_tweet.Event.Event == "favorite" {
				// fav
				ydb := dbyoutube.NewDBYoutubeMovie()
				var youtubedb dbyoutube.YoutubeMovie = ydb
				tweet := *event_tweet.TargetObject
				movie := youtubedb.ScanYoutubeMovie(tweet)
				if movie == nil {
					continue
				}

				fdb := dbfav.NewDBFav()
				var favdb dbfav.Fav = fdb

				udb := dbuser.NewDBUser()
				var userdb dbuser.User = udb

				user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

				if user.Id != 0 {
					_ = favdb.Add(user.Id, movie.Id)
					fmt.Printf("Add Fav Log: %v \n", movie.Id)
				}
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

func preparePidFile() {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)
	pidByte := []byte(pidStr + "\n")
	ioutil.WriteFile("../tmp/pids/userstream.pid", pidByte, os.ModePerm)
}

func removePidFile() {
	os.Remove("../tmp/pids/userstream.pid")
}
