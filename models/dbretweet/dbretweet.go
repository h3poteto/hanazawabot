package dbretweet

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"../basedb"
)

type Retweet interface {
	Add(int, int) bool
}

type DBRetweet struct {
	dbobject basedb.DB
}

func (self *DBRetweet) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBRetweet() *DBRetweet {
	dbretweet := &DBRetweet{}
	dbretweet.Initialize()
	return dbretweet
}

func (u *DBRetweet) Add(user_id int, youtube_movie_id int) bool {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into youtube_movie_retweets (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}
	return true
}
