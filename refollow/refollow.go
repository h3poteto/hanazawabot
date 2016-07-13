package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(os.Getenv("TWITTER_OAUTH_TOKEN"), os.Getenv("TWITTER_OAUTH_SECRET"))

	// 自分の情報取得なので空
	empty_values := url.Values{}

	chanFollowers := api.GetFollowersIdsAll(empty_values)
	followers := <-chanFollowers

	chanFriends := api.GetFriendsIdsAll(empty_values)
	friends := <-chanFriends
	var diff []int64
	for _, follower := range followers.Ids {
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
			fmt.Printf("twitter follow error: %v", err)
		} else {
			fmt.Printf("follow new user: %d \n", user.Id)
		}
	}

}
