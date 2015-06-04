package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"
	"os"

	"code.google.com/p/google-api-go-client/googleapi/transport"
	"code.google.com/p/google-api-go-client/youtube/v3"
	_ "github.com/go-sql-driver/mysql"
)

var (
	query = flag.String("query", "花澤香菜", "search term")
	maxResults = flag.Int64("max-results", 50, "Max Youtube results")
)


type YoutubeMovie struct {
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
	call := service.Search.List("id,snippet").
		Q(*query).
		MaxResults(*maxResults).
		Order("date").
		PublishedAfter(yesterday.Format(time.RFC3339))
	response, err := call.Do()
	if err != nil {
		log.Fatalf("error making search API call: %v", err)
	}

	videos := make(map[string]YoutubeMovie)
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = YoutubeMovie{item.Snippet.Title, item.Snippet.Description, item.Snippet.Thumbnails.Default.Url}
		}
	}

	printIDs("videos", videos)

	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	for id, youtube := range videos {

		_, err := db.Exec("insert into youtube_movies (title, movie_id, description, disabled, created_at) values (?, ?, ?, ?, ?)", youtube.title, id, youtube.description, 0, time.Now())
		if err != nil {
			log.Fatalf("mysql connect error: %v", err)
		}
	}
}

func printIDs(sectionName string, matches map[string]YoutubeMovie) {
	fmt.Printf("%v:\n", sectionName)
	for id, youtube := range matches {
		fmt.Printf("[%v] %v : %v \n", id, youtube.title, youtube.description)
	}
	fmt.Printf("\n\n")
}
