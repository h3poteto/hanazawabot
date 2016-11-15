package dbtweet

import (
	"database/sql"
	"time"

	"../basedb"

	"github.com/pkg/errors"
)

type Tweet interface {
	Add(int, string, int64) error
}

type DBTweet struct {
	database *sql.DB
}

func (self *DBTweet) Initialize() {
	self.database = basedb.SharedInstance().Connection
}

func NewDBTweet() *DBTweet {
	dbtweet := &DBTweet{}
	dbtweet.Initialize()
	return dbtweet
}

func (u *DBTweet) Add(user_id int, tweet string, tweet_id int64) error {
	_, err := u.database.Exec("insert into tweets (user_id, tweet, tweet_id, created_at) values (?, ?, ?, ?)", user_id, tweet, tweet_id, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}

	return nil
}
