package cmd

import (
	"os"

	"../kanachan"
	"../models/dbserif"
	"../models/dbyoutube"
	"../modules/logging"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
)

func tweetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tweet",
		Short: "Tweet a new youtube movie",
		Run:   tweet,
	}

	return cmd
}

func tweet(cmd *cobra.Command, args []string) {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))

	ydb := dbyoutube.NewDBYoutubeMovie()
	var youtube_movie dbyoutube.YoutubeMovie = ydb

	movie, err := youtube_movie.SelectRandom()
	if err != nil {
		logging.SharedInstance().MethodInfo("tweet").Fatalf("DBYoutube random select error: %v", err)
	}

	sdb := dbserif.NewDBSerif()
	var serif dbserif.Serif = sdb
	tweet_serif, err := serif.SelectRandom()
	if err != nil {
		logging.SharedInstance().MethodInfo("tweet").Fatalf("DBSerif random select error: %v", err)
	}

	aKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = aKana
	movieID, err := movie.ConvertYoutubeID()
	if err != nil {
		logging.SharedInstance().MethodInfo("tweet").Fatal(err)
	}
	_, err = api.PostTweet(
		kana.BuildTweet("",
			tweet_serif,
			movie.Title,
			movieID),
		nil)
	if err != nil {
		logging.SharedInstance().MethodInfo("tweet").Fatalf("Twitter API error: %v", err)
	}
	logging.SharedInstance().MethodInfo("tweet").Infof("tweet success: %v", movie.Title)
}
