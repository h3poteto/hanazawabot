package dbtweet

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"../basedb"
)

type Tweet interface {
	Add(int, string, int64) error
}

type DBTweet struct {
	dbobject basedb.DB
}

func (self *DBTweet) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBTweet() *DBTweet {
	dbtweet := &DBTweet{}
	dbtweet.Initialize()
	return dbtweet
}

func (u *DBTweet) Add(user_id int, tweet string, tweet_id int64) error {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into tweets (user_id, tweet, tweet_id, created_at) values (?, ?, ?, ?)", user_id, tweet, tweet_id, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}

	return nil
}
