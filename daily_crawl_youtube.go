package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
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

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	vKana := &kanachan.Kanachan{}
	var kana kanachan.Kana = vKana

	for _, query := range query_words {
		call := service.Search.List("id,snippet").
			Q(query).
			MaxResults(*maxResults).
			Order("date").
			PublishedAfter(yesterday.Format(time.RFC3339))
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

		printIDs("videos", videos)

		myDb := &dbyoutube.DBYoutubeMovie{}
		var db dbyoutube.YoutubeMovie = myDb

		for id, youtube := range videos {
			db.Add(youtube.title, id, youtube.description)
		}
	}
}

func printIDs(sectionName string, matches map[string]Youtube) {
	fmt.Printf("%v:\n", sectionName)
	for id, youtube := range matches {
		fmt.Printf("[%v] %v : %v \n", id, youtube.title, youtube.description)
	}
	fmt.Printf("\n\n")
}
