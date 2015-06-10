package main

import (
	"fmt"
	"os"
	"net/url"
	"log"


	"github.com/ChimeraCoder/anaconda"

)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))

	// 自分の情報取得なので空
	empty_values := url.Values{}

	values := url.Values{}
	first_followers, err := api.GetFollowersIds(values)
	if err != nil {
		log.Fatalf("twitter api error: %v", err)
	}
	followers := first_followers.Ids
	next_cursor := first_followers.Next_cursor_str
	for next_cursor != "0" {
		values.Set("cursor", next_cursor)
		f, _ := api.GetFollowersIds(values)
		followers = append(followers, f.Ids...)
		next_cursor = f.Next_cursor_str
	}

	friends, _ := api.GetFriendsIdsAll(empty_values)
	var diff []int64
	for _, follower := range followers {
		found := false
		for _, friend := range friends.Ids {
			if friend == follower {
				found = true
				break
			}
		}
		if found == false {
			diff = append(diff, follower)
		}
	}

	for _, follower := range diff {
		user, err := api.FollowUserId(follower, empty_values)
		if err != nil {
			log.Fatalf("twitter follow error: %v", err)
		} else {
			fmt.Printf("follow new user: %d \n", user.Id)
		}
	}


}
