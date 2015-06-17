package dbuser

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type User interface {
	Add(int64, string) bool
	Select(int64, string) DBUser
	SelectOrAdd(int64, string) DBUser
}

type DBUser struct {
	Id int
	ScreenName string
	TwitterID int64
}

func (u *DBUser) Add(twitter_id int64, screen_name string) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		return false
	}

	_, err = db.Exec("insert into users (twitter_id, screen_name, created_at) values (?, ?, ?)", twitter_id, screen_name, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}

	defer db.Close()

	return true
}

func (u *DBUser) Select(twitter_id int64, screen_name string) DBUser {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	var rows *sql.Rows
	if twitter_id > 0 {
		rows, err = db.Query("select * from users where twitter_id = ?;", twitter_id)
	} else {
		rows, err = db.Query("select * from users where screen_name = ?;", screen_name)
	}

	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
	}

	defer db.Close()

	var twitter int64
	id, twitter, screen_name, created_at, updated_at := 0, 0, "", "", ""
	for rows.Next() {
		err = rows.Scan(&id, &twitter, &screen_name, &created_at, &updated_at)
		if err != nil {
			fmt.Printf("mysql select error: %v \n", err)
		}
	}

	user := DBUser{Id: id, TwitterID: twitter, ScreenName: screen_name}
	return user
}

func (u *DBUser) SelectOrAdd(twitter_id int64, screen_name string) DBUser {
	user := u.Select(twitter_id, screen_name)
	if user.Id == 0 {
		result := u.Add(twitter_id, screen_name)
		if result {
			user = u.Select(twitter_id, screen_name)
		}
	}
	return user
}
