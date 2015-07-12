package dbtweet

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"../basedb"
)

type Tweet interface {
	Add(int, string, int64) bool
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

func (u *DBTweet) Add(user_id int, tweet string, tweet_id int64) bool {
	db := u.dbobject.Init()

	_, err := db.Exec("insert into tweets (user_id, tweet, tweet_id, created_at) values (?, ?, ?, ?)", user_id, tweet, tweet_id, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	return true
}
