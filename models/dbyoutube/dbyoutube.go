package dbyoutube

import (
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"../basedb"
)

var (
	youtube_prefix = "https://www.youtube.com/watch?v="
)

type YoutubeMovie interface {
	Add(string, string, string) error
	SelectRandom() (*DBYoutubeMovie, error)
	SelectToday() (*DBYoutubeMovie, error)
	ConvertYoutubeID() (string, error)
	revertYoutubeID(string) (string, error)
	ScanYoutubeMovie(anaconda.Tweet) *DBYoutubeMovie
	Select(string) (*DBYoutubeMovie, error)
}

type DBYoutubeMovie struct {
	Id          int
	Title       string
	MovieId     string
	Description string
	Used        bool
	Disabled    bool
	dbobject    basedb.DB
}

func (self *DBYoutubeMovie) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBYoutubeMovie() *DBYoutubeMovie {
	dbyoutube := &DBYoutubeMovie{}
	dbyoutube.Initialize()
	return dbyoutube
}

func (u *DBYoutubeMovie) Add(title string, id string, description string) error {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into youtube_movies (title, movie_id, description, disabled, created_at) values (?, ?, ?, ?, ?)", title, id, description, 0, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}

	return nil
}

func (u *DBYoutubeMovie) SelectRandom() (*DBYoutubeMovie, error) {
	db := u.dbobject.Init()
	defer db.Close()

	rows, err := db.Query("select id, title, movie_id, description, used, disabled from youtube_movies order by rand() limit 1;")
	if err != nil {
		return nil, errors.Wrap(err, "mysql select error")
	}

	var id int
	var movie_id, title, description string
	var used, disabled bool

	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled)
		if err != nil {
			return nil, errors.Wrap(err, "mysql scan error")
		}
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: used, Disabled: disabled}, nil

}

func (u *DBYoutubeMovie) SelectToday() (*DBYoutubeMovie, error) {
	db := u.dbobject.Init()
	defer db.Close()

	yesterday := time.Now().Add(-24 * time.Hour)
	rows, err := db.Query("select id, title, movie_id, description, used, disabled from youtube_movies where created_at > ? and used = 0 limit 1;", yesterday)
	if err != nil {
		return nil, errors.Wrap(err, "mysql select error")
	}

	var id int
	var movie_id, title, description string
	var used, disabled bool

	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled)
		if err != nil {
			return nil, errors.Wrap(err, "mysql scan error")
		}
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: used, Disabled: disabled}, nil
}

func (u *DBYoutubeMovie) ConvertYoutubeID() (string, error) {
	if u.MovieId == "" {
		return "", errors.New("cannot find youtube movie")
	}
	return youtube_prefix + u.MovieId, nil
}

func (u *DBYoutubeMovie) revertYoutubeID(url string) (string, error) {
	params := strings.Replace(url, youtube_prefix, "", -1)
	movies := strings.Split(params, "?")
	for _, id := range movies {
		return id, nil
	}
	return "", errors.New("cannot find youtube movie")
}

func (u *DBYoutubeMovie) ScanYoutubeMovie(tweet anaconda.Tweet) *DBYoutubeMovie {
	db := u.dbobject.Init()
	defer db.Close()

	for _, url := range tweet.Entities.Urls {
		youtubeID, err := u.revertYoutubeID(url.Expanded_url)
		if err != nil {
			continue
		}
		movie, _ := u.Select(youtubeID)
		if movie != nil {
			return movie
		}
	}
	return nil
}

func (u *DBYoutubeMovie) Select(keyword string) (*DBYoutubeMovie, error) {
	db := u.dbobject.Init()
	defer db.Close()

	rows, err := db.Query("select id, movie_id, title, description, used, disabled from youtube_movies where movie_id like '%" + keyword + "%' or title like '%" + keyword + "%' or description like '%" + keyword + "%';")
	if err != nil {
		return nil, errors.Wrap(err, "mysql select error")
	}

	var id int
	var movie_id, title, description string
	var used, disabled bool
	for rows.Next() {
		err = rows.Scan(&id, &title, &movie_id, &description, &used, &disabled)
		if err != nil {
			return nil, errors.Wrap(err, "mysql scan error")
		}
	}

	return &DBYoutubeMovie{Id: id, Title: title, MovieId: movie_id, Description: description, Used: used, Disabled: disabled}, nil
}
