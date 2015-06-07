package dbyoutube

import (
	"fmt"
	_ "log"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type YoutubeMovie interface {
	Add(string, string, string) bool
}

type DBYoutubeMovie struct {
}


func (u *DBYoutubeMovie) Add(title string, id string, description string) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("insert into youtube_movies (title, movie_id, description, disabled, created_at) values (?, ?, ?, ?, ?)", title, id, description, 0, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	db.Close()

	return true
}
