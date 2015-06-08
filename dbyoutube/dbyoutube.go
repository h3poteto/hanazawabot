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
	SelectRandom() (tweet string, error string)
	convertYoutubeID(string) string
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

	defer db.Close()

	return true
}

func (u *DBYoutubeMovie) SelectRandom() (tweet string, error string) {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("select * from youtube_movies order by rand() limit 1;")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	id, movie_id, title, description, disabled, created_at, updated_at := 0, "", "", "", 0 , "", ""
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &disabled, &created_at, &updated_at)
		if err != nil{
			return "", err.Error()
		}
	}

	return u.convertYoutubeID(movie_id), ""
}

func (u *DBYoutubeMovie) convertYoutubeID(movie_id string) string {
	return "https://www.youtube.com/watch?v=" + movie_id
}
