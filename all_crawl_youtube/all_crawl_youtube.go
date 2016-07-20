package main

import (
	"flag"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"../kanachan"
	"../models/dbyoutube"
	"../modules/logging"
)

var (
	maxResults  = flag.Int64("max-results", 50, "Max Youtube results")
	query_words = [...]string{"花澤香菜", "花澤病"}
)

type Youtube struct {
	title       string
	description string
	thumbnail   string
}

func main() {
	flag.Parse()

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

	for _, query := range query_words {

		nextPageToken := "hanazawabot_token"

		for nextPageToken != "" {
			if nextPageToken == "hanazawabot_token" {
				nextPageToken = ""
			}
			call := service.Search.List("id,snippet").
				Q(query).
				MaxResults(*maxResults).
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
