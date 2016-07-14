package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"../kanachan"
	"../models/dbfav"
	"../models/dbretweet"
	"../models/dbserif"
	"../models/dbtweet"
	"../models/dbuser"
	"../models/dbyoutube"
	"github.com/ChimeraCoder/anaconda"
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

			aKana := &kanachan.Kanachan{}
			var kana kanachan.Kana = aKana
			if kana.IncludeCheck(tweet.Text) {
				// contain
				if saveKanaTweet(tweet) != nil {
					// TODO: error log
				}
			}

			if isRetweet(tweet, self) {
				// retweet
				if saveRetweet(tweet) != nil {
					// TODO: error log
				}

			} else if isReply(tweet, self) {
				// reply
				if replyKanaMovie(tweet, api, kana) != nil {
					// TODO: error log
				}
			} else {
				continue
			}
		case anaconda.EventTweet:
			event_tweet := event.(anaconda.EventTweet)
			if event_tweet.Event.Event == "favorite" {
				// fav
				tweet := *event_tweet.TargetObject
				if saveFav(tweet) != nil {
					// TODO: error log
				}
			}
		}
	}
}

func isReply(tweet anaconda.Tweet, user anaconda.User) bool {
	if tweet.InReplyToUserID == user.Id {
		return true
	}
	if strings.Contains(tweet.Text, "@"+user.ScreenName) {
		return true
	}
	return false
}

func isRetweet(tweet anaconda.Tweet, user anaconda.User) bool {
	if strings.Index(tweet.Text, "RT") == 0 && tweet.User.Id != user.Id && strings.Contains(tweet.Text, "@"+user.ScreenName) {
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

func saveKanaTweet(tweet anaconda.Tweet) error {
	tdb := dbtweet.NewDBTweet()
	var tweetdb dbtweet.Tweet = tdb

	udb := dbuser.NewDBUser()
	var userdb dbuser.User = udb

	user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

	if user.Id == 0 {
		return errors.New("cannot select or add user")
	}
	tweetdb.Add(user.Id, tweet.Text, tweet.Id)
	fmt.Printf("Add Tweet Log: %v \n", tweet.Text)
	return nil
}

func saveRetweet(tweet anaconda.Tweet) error {
	ydb := dbyoutube.NewDBYoutubeMovie()
	var youtubedb dbyoutube.YoutubeMovie = ydb
	movie := youtubedb.ScanYoutubeMovie(tweet)
	if movie == nil {
		return nil
	}
	rdb := dbretweet.NewDBRetweet()
	var retweetdb dbretweet.Retweet = rdb

	udb := dbuser.NewDBUser()
	var userdb dbuser.User = udb

	user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

	if user.Id == 0 {
		return errors.New("cannot select or add user")
	}
	retweetdb.Add(user.Id, movie.Id)
	fmt.Printf("Add Rewteet Log: %v \n", movie.Id)
	return nil
}

func replyKanaMovie(tweet anaconda.Tweet, api *anaconda.TwitterApi, kana kanachan.Kana) error {
	ydb := dbyoutube.NewDBYoutubeMovie()
	var youtubedb dbyoutube.YoutubeMovie = ydb
	movie := youtubedb.SelectRandom()
	if movie == nil {
		return errors.New("DBYoutube random select error")
	}
	tweet_value := url.Values{}
	tweet_value.Set("in_reply_to_status_id", tweet.IdStr)

	sdb := dbserif.NewDBSerif()
	var serif dbserif.Serif = sdb
	tweet_serif, err := serif.SelectRandom()
	if err != nil {
		fmt.Printf("DBSerif random select error: %v\n", err)
		return err
	}

	_, err = api.PostTweet(
		kana.BuildTweet("@"+tweet.User.ScreenName+" ",
			tweet_serif,
			movie.Title,
			movie.ConvertYoutubeID()),
		tweet_value)
	if err != nil {
		fmt.Printf("twitter api error: %v\n", err)
		return err
	}
	fmt.Printf("reply to @%v: %v\n", tweet.User.ScreenName, "@"+tweet.User.ScreenName)
	return nil
}

func saveFav(tweet anaconda.Tweet) error {
	ydb := dbyoutube.NewDBYoutubeMovie()
	var youtubedb dbyoutube.YoutubeMovie = ydb
	movie := youtubedb.ScanYoutubeMovie(tweet)
	if movie == nil {
		return nil
	}

	fdb := dbfav.NewDBFav()
	var favdb dbfav.Fav = fdb

	udb := dbuser.NewDBUser()
	var userdb dbuser.User = udb

	user := userdb.SelectOrAdd(tweet.User.Id, tweet.User.ScreenName)

	if user.Id == 0 {
		return errors.New("cannot select or add user")
	}
	_ = favdb.Add(user.Id, movie.Id)
	fmt.Printf("Add Fav Log: %v \n", movie.Id)
	return nil
}
