package dbuser

import (
	"fmt"
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
}

type User interface {
	Add(int64, string) bool
	Select(int64, string) sql.DB
}

type DBUser struct {
}

func (u *DBUser) Add(twitter_id int64, screen_name string) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("insert into user (twitter_id, screen_name, created_at) values (?, ?, ?)", twitter_id, screen_name, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	return true
}

func (u *DBUser) Select(twitter_id int64, screen_name string) sql.DB {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
}
