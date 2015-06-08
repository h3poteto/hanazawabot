package main

import (
	"fmt"
	"os"
	"log"

	"github.com/ChimeraCoder/anaconda"

	"./dbyoutube"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))


	myDb := &dbyoutube.DBYoutubeMovie{}
	var db dbyoutube.YoutubeMovie = myDb

	youtube_url, err := db.SelectRandom()
	if err != "" {
		log.Fatalf("DBYoutube random select error: %v", err)
	}
	fmt.Printf("%v \n", youtube_url)

	_, error := api.PostTweet("わーい" + youtube_url, nil)
	if error != nil {
		log.Fatalf("twitter api error: %v", error)
	}
}
