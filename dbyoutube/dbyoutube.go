package dbyoutube

import (
	"fmt"
	"time"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ChimeraCoder/anaconda"
)

var (
	youtube_prefix = "https://www.youtube.com/watch?v="
)

type YoutubeMovie interface {
	Add(string, string, string) bool
	SelectRandom() *DBYoutubeMovie
	convertYoutubeID(string) string
	revertYoutubeID(string) string
	ScanYoutubeMovie(anaconda.Tweet) *DBYoutubeMovie
	Select(int, string) *DBYoutubeMovie
	GetRandomMovieURL() string
}

type DBYoutubeMovie struct {
	Id int
	Title string
	MovieId string
	Description string
	Disabled bool
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

func (u *DBYoutubeMovie) SelectRandom() *DBYoutubeMovie {
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
			return nil
		}
	}

	fdisabled := false
	if disabled != 0 {
		fdisabled = true
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Disabled: fdisabled}

}

func (u *DBYoutubeMovie) GetRandomMovieURL() string {
	youtube_movie := u.SelectRandom()
	if youtube_movie == nil {
		return ""
	} else {
		return u.convertYoutubeID(youtube_movie.MovieId)
	}
}

func (u *DBYoutubeMovie) convertYoutubeID(movie_id string) string {
	if movie_id != "" {
		return youtube_prefix + movie_id
	} else {
		fmt.Printf("cannot found youtue movie in db")
		return ""
	}
}

func (u *DBYoutubeMovie) revertYoutubeID(url string) string {
	params := strings.Replace(url, youtube_prefix, "", -1)
	movies := strings.Split(params, "?")
	for _, id := range movies {
		return id
	}
	return ""
}

func (u *DBYoutubeMovie) ScanYoutubeMovie(tweet anaconda.Tweet) *DBYoutubeMovie {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	for _, url := range tweet.Entities.Urls {
		movie := u.Select(0, u.revertYoutubeID(url.Expanded_url))
		if movie != nil {
			return movie
		}
	}
	return nil
}

func (u *DBYoutubeMovie) Select(id int, keyword string) *DBYoutubeMovie {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return nil
	}

	rows, err := db.Query("select * from youtube_movies where movie_id like '%" + keyword + "%' or title like '%" + keyword + "%' or description like '%" + keyword + "%';")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return nil
	}

	id, movie_id, title, description, disabled, created_at, updated_at := 0, "", "", "", 0 , "", ""
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &disabled, &created_at, &updated_at)
		if err != nil{
			return nil
		}
	}
	fdisabled := false
	if disabled != 0 {
		fdisabled = true
	}
	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Disabled: fdisabled}
}
