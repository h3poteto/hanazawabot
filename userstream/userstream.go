package main

import (
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
	"../modules/logging"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pkg/errors"
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
				err := saveKanaTweet(tweet)
				if err != nil {
					logging.SharedInstance().MethodInfo("userstream").Error(err)
					continue
				}
			}

			if isRetweet(tweet, self) {
				// retweet
				err := saveRetweet(tweet)
				if err != nil {
					logging.SharedInstance().MethodInfo("userstream").Error(err)
					continue
				}

			} else if isReply(tweet, self) {
				// reply
				err := replyKanaMovie(tweet, api, kana)
				if err != nil {
					logging.SharedInstance().MethodInfo("userstream").Error(err)
					continue
				}
			} else {
				continue
			}
		case anaconda.EventTweet:
			event_tweet := event.(anaconda.EventTweet)
			if event_tweet.Event.Event == "favorite" {
				// fav
				tweet := *event_tweet.TargetObject
				err := saveFav(tweet)
				if err != nil {
					logging.SharedInstance().MethodInfo("userstream").Error(err)
					continue
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
	logging.SharedInstance().MethodInfo("userstream").Infof("Add tweet log: %v", tweet.Text)
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
	err := retweetdb.Add(user.Id, movie.Id)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("userstream").Infof("Add retweet log: %v", movie.Id)
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
		return err
	}

	_, err = api.PostTweet(
		kana.BuildTweet("@"+tweet.User.ScreenName+" ",
			tweet_serif,
			movie.Title,
			movie.ConvertYoutubeID()),
		tweet_value)
	if err != nil {
		return errors.Wrap(err, "twitter api error")
	}
	logging.SharedInstance().MethodInfo("userstream").Infof("reply to @%v", tweet.User.ScreenName)
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
	err := favdb.Add(user.Id, movie.Id)
	if err != nil {
		return err
	}
	logging.SharedInstance().MethodInfo("userstream").Infof("Add fav log: %v", movie.Id)
	return nil
}
