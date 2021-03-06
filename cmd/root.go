package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd has commands for hanazawabot
var RootCmd = &cobra.Command{
	Use:           "hanazawabot",
	Short:         "Twitter bot for hanazawabot",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize()
	RootCmd.AddCommand(
		migrateCmd(),
		seedCmd(),
		userstreamCmd(),
		tweetCmd(),
		refollowCmd(),
		allCrawlYoutubeCmd(),
		dailyCrawlYoutubeCmd(),
		dailyTweetCmd(),
	)
}
