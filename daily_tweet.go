package main

import (
	"fmt"
	"os"
	"log"

	"github.com/ChimeraCoder/anaconda"

	"./dbyoutube"
	"./dbserif"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))


	ydb := &dbyoutube.DBYoutubeMovie{}
	var youtube_movie dbyoutube.YoutubeMovie = ydb

	youtube_url := youtube_movie.GetTodayMovieURL()
	if youtube_url == "" {
		log.Fatalf("DBYoutube today select error")
	}


	sdb := &dbserif.DBSerif{}
	var serif dbserif.Serif = sdb
	tweet_serif, err := serif.SelectRandom()
	if err != "" {
		log.Fatalf("DBSerif random select error: %v", err)
	}

	_, error := api.PostTweet(tweet_serif + " " + youtube_url, nil)
	if error != nil {
		log.Fatalf("twitter api error: %v", error)
	}
	fmt.Printf("tweet: %v \n", tweet_serif + " " + youtube_url)
}
