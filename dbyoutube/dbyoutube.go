package dbyoutube

import (
	"fmt"
	"time"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ChimeraCoder/anaconda"

	"../database"
)

var (
	youtube_prefix = "https://www.youtube.com/watch?v="
)

type YoutubeMovie interface {
	Add(string, string, string) bool
	SelectRandom() *DBYoutubeMovie
	SelectToday() *DBYoutubeMovie
	ConvertYoutubeID() string
	revertYoutubeID(string) string
	ScanYoutubeMovie(anaconda.Tweet) *DBYoutubeMovie
	Select(int, string) *DBYoutubeMovie
}

type DBYoutubeMovie struct {
	Id int
	Title string
	MovieId string
	Description string
	Used bool
	Disabled bool
	dbobject database.DB
}

func (self *DBYoutubeMovie) Initialize() {
	myDatabase := &database.Database{}
	var myDb database.DB = myDatabase
	self.dbobject = myDb
}

func NewDBYoutubeMovie() *DBYoutubeMovie {
	dbyoutube := &DBYoutubeMovie{}
	dbyoutube.Initialize()
	return dbyoutube
}

func (u *DBYoutubeMovie) Add(title string, id string, description string) bool {
	db := u.dbobject.Init()

	_, err := db.Exec("insert into youtube_movies (title, movie_id, description, disabled, created_at) values (?, ?, ?, ?, ?)", title, id, description, 0, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	return true
}

func (u *DBYoutubeMovie) SelectRandom() *DBYoutubeMovie {
	db := u.dbobject.Init()

	rows, err := db.Query("select * from youtube_movies order by rand() limit 1;")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	id, movie_id, title, description,used, disabled, created_at, updated_at := 0, "", "", "", 0, 0 , "", ""
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled, &created_at, &updated_at)
		if err != nil{
			return nil
		}
	}

	fdisabled := false
	if disabled != 0 {
		fdisabled = true
	}

	fused := false
	if used != 0 {
		fused = true
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: fused, Disabled: fdisabled}

}

func (u *DBYoutubeMovie) SelectToday() *DBYoutubeMovie {
	db := u.dbobject.Init()

	yesterday := time.Now().Add(-24 * time.Hour)
	rows, err := db.Query("select * from youtube_movies where created_at > ? and used = 0 limit 1;", yesterday)
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	id, movie_id, title, description, used, disabled, created_at, updated_at := 0, "", "", "", 0, 0 , "", ""
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled, &created_at, &updated_at)
		if err != nil{
			return nil
		}
	}

	fdisabled := false
	if disabled != 0 {
		fdisabled = true
	}
	fused := false
	if used != 0 {
		fused = true
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: fused, Disabled: fdisabled}
}


func (u *DBYoutubeMovie) ConvertYoutubeID() string {
	if u.MovieId != "" {
		return youtube_prefix + u.MovieId
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
	db := u.dbobject.Init()
	defer db.Close()

	for _, url := range tweet.Entities.Urls {
		movie := u.Select(0, u.revertYoutubeID(url.Expanded_url))
		if movie != nil {
			return movie
		}
	}
	return nil
}

func (u *DBYoutubeMovie) Select(id int, keyword string) *DBYoutubeMovie {
	db := u.dbobject.Init()
	defer db.Close()

	rows, err := db.Query("select * from youtube_movies where movie_id like '%" + keyword + "%' or title like '%" + keyword + "%' or description like '%" + keyword + "%';")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return nil
	}

	id, movie_id, title, description, used, disabled, created_at, updated_at := 0, "", "", "", 0, 0 , "", ""
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled, &created_at, &updated_at)
		if err != nil{
			return nil
		}
	}
	fdisabled := false
	if disabled != 0 {
		fdisabled = true
	}

	fused := false
	if used != 0 {
		fused = true
	}
	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: fused, Disabled: fdisabled}
}
