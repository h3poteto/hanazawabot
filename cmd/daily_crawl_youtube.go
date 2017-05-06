package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/h3poteto/hanazawabot/kanachan"
	"github.com/h3poteto/hanazawabot/models/dbyoutube"
	"github.com/h3poteto/hanazawabot/modules/logging"

	"github.com/spf13/cobra"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// DailyCrawl is struct for dailyCrowdYoutube
type DailyCrawl struct {
	maxResults int64
	queryWords []string
}

func dailyCrawlYoutubeCmd() *cobra.Command {
	d := &DailyCrawl{
		queryWords: []string{"花澤香菜", "花澤病"},
	}
	cmd := &cobra.Command{
		Use:   "dailycrawl",
		Short: "Crawl today's youtube movies",
		Run:   d.dailyCrawlYoutube,
	}

	flags := cmd.Flags()
	flags.Int64VarP(&d.maxResults, "max", "m", int64(50), "Max youtube results")

	return cmd
}

func (d *DailyCrawl) dailyCrawlYoutube(cmd *cobra.Command, args []string) {
	developerKey := os.Getenv("DEVELOPER_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		logging.SharedInstance().MethodInfo("daily_crawl_youtube").Fatalf("Error creating youtube client: %v", err)
	}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	vKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = vKana

	for _, query := range d.queryWords {
		call := service.Search.List("id,snippet").
			Q(query).
			MaxResults(d.maxResults).
			Order("date").
			PublishedAfter(yesterday.Format(time.RFC3339))
		response, err := call.Do()

		if err != nil {
			logging.SharedInstance().MethodInfo("daily_crawl_youtube").Fatalf("Error making search API call: %v", err)
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

		myDb := dbyoutube.NewDBYoutubeMovie()
		var db dbyoutube.YoutubeMovie = myDb

		for id, youtube := range videos {
			err := db.Add(youtube.title, id, youtube.description)
			if err != nil {
				logging.SharedInstance().MethodInfo("daily_crawl_youtube").Info(err)
				continue
			}
			logging.SharedInstance().MethodInfo("daily_crawl_youtube").Infof("Add youtube_movies to: %v", youtube.title)
		}
	}
}
