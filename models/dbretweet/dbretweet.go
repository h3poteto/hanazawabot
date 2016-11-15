package dbretweet

import (
	"database/sql"
	"time"

	"../basedb"

	"github.com/pkg/errors"
)

type Retweet interface {
	Add(int, int) error
}

type DBRetweet struct {
	database *sql.DB
}

func (self *DBRetweet) Initialize() {
	self.database = basedb.SharedInstance().Connection
}

func NewDBRetweet() *DBRetweet {
	dbretweet := &DBRetweet{}
	dbretweet.Initialize()
	return dbretweet
}

func (u *DBRetweet) Add(user_id int, youtube_movie_id int) error {
	_, err := u.database.Exec("insert into youtube_movie_retweets (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}
	return nil
}
