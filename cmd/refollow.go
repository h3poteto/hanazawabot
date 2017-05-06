package cmd

import (
	"net/url"
	"os"

	"github.com/h3poteto/hanazawabot/modules/logging"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/cobra"
)

func refollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refollow",
		Short: "Refollow all user",
		Run:   refollow,
	}

	return cmd
}

func refollow(cmd *cobra.Command, args []string) {
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
			logging.SharedInstance().MethodInfo("refollow").Errorf("twitter follow error: %v", err)
		}
		logging.SharedInstance().MethodInfo("refollow").Infof("follow new user: %d", user.Id)
	}

}
