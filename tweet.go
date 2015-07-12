package main

import (
	"os"
	"log"

	"github.com/ChimeraCoder/anaconda"

	"./dbyoutube"
	"./dbserif"
	"./kanachan"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))


	ydb := dbyoutube.NewDBYoutubeMovie()
	var youtube_movie dbyoutube.YoutubeMovie = ydb

	movie := youtube_movie.SelectRandom()
	if movie == nil {
		log.Fatalf("DBYoutube random select error")
	}

	sdb := dbserif.NewDBSerif()
	var serif dbserif.Serif = sdb
	tweet_serif, err := serif.SelectRandom()
	if err != "" {
		log.Fatalf("DBSerif random select error: %v", err)
	}

	aKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = aKana
	_, error := api.PostTweet(
		kana.BuildTweet("",
			tweet_serif,
			movie.Title,
			movie.ConvertYoutubeID()),
		nil)
	if error != nil {
		log.Fatalf("twitter api error: %v", error)
	}
}
