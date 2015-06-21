package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"

	"./dbyoutube"
	"./kanachan"
)

var (
	maxResults = flag.Int64("max-results", 50, "Max Youtube results")
	query_words = [...]string{"花澤香菜", "花澤病"}
)


type Youtube struct {
	title string
	description string
	thumbnail string
}

func main() {
	flag.Parse()

	developerKey := os.Getenv("DEVELOPER_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating Youtube client: %v", err)
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
				log.Fatalf("error making search API call: %v", err)
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

			myDb := &dbyoutube.DBYoutubeMovie{}
			var db dbyoutube.YoutubeMovie = myDb

			for id, youtube := range videos {
				db.Add(youtube.title, id, youtube.description)
				fmt.Printf("add youtube_movies to: %v \n", youtube.title)
			}
		}
	}
}
