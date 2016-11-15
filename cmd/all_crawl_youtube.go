package cmd

import (
	"net/http"
	"os"

	"../kanachan"
	"../models/dbyoutube"
	"../modules/logging"

	"github.com/spf13/cobra"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	title       string
	description string
	thumbnail   string
}

// AllCrawl is struct for allCrawlYoutube
type AllCrawl struct {
	maxResults int64
	queryWords []string
}

func allCrawlYoutubeCmd() *cobra.Command {
	a := &AllCrawl{
		queryWords: []string{"花澤香菜", "花澤病"},
	}
	cmd := &cobra.Command{
		Use:   "allcrawl",
		Short: "Crawl youtube movies whole terms",
		Run:   a.allCrawlYoutube,
	}

	flags := cmd.Flags()
	flags.Int64VarP(&a.maxResults, "max", "m", int64(50), "Max youtube results")

	return cmd
}

func (a *AllCrawl) allCrawlYoutube(cmd *cobra.Command, args []string) {
	developerKey := os.Getenv("DEVELOPER_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		logging.SharedInstance().MethodInfo("all_crawl_youtube").Fatalf("Error createing youtube clinet: %v", err)
	}

	vKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = vKana

	for _, query := range a.queryWords {

		nextPageToken := "hanazawabot_token"

		for nextPageToken != "" {
			if nextPageToken == "hanazawabot_token" {
				nextPageToken = ""
			}
			call := service.Search.List("id,snippet").
				Q(query).
				MaxResults(a.maxResults).
				Type("video").
				PageToken(nextPageToken)

			response, err := call.Do()
			if err != nil {
				logging.SharedInstance().MethodInfo("all_crawl_youtube").Fatalf("Error making search API call: %v", err)
			}

			videos := make(map[string]Youtube)
			for _, item := range response.Items {
				switch item.Id.Kind {
				case "youtube#video":
					if kana.ExceptCheck(item.Snippet.Title) && kana.ExceptCheck(item.Snippet.Description) {
						videos[item.Id.VideoId] = Youtube{item.Snippet.Title, item.Snippet.Description, item.Snippet.Thumbnails.Default.Url}
					}
				}
			}
			nextPageToken = response.NextPageToken

			myDb := dbyoutube.NewDBYoutubeMovie()
			var db dbyoutube.YoutubeMovie = myDb

			for id, youtube := range videos {
				err := db.Add(youtube.title, id, youtube.description)
				if err != nil {
					logging.SharedInstance().MethodInfo("all_crawl_youtube").Info(err)
					continue
				}
				logging.SharedInstance().MethodInfo("all_crawl_youtube").Infof("add youtube_movies to: %v", youtube.title)
			}
		}
	}
}
