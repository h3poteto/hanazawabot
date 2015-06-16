package dbtweet

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Tweet interface {
	Add(int, string, int64) bool
}

type DBTweet struct {
}

func (u *DBTweet) Add(user_id int, tweet string, tweet_id int64) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("insert into tweets (user_id, tweet, tweet_id, created_at) values (?, ?, ?, ?)", user_id, tweet, tweet_id, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	return true
}
