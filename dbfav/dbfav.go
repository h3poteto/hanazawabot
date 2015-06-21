package dbfav

import (
	"fmt"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Fav interface {
	Add(int, int) bool
}

type DBFav struct {
}

func (u *DBFav) Add(user_id int, youtube_movie_id int) bool {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	defer db.Close()
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}

	_, err = db.Exec("insert into youtube_movie_favs (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}
	return true
}
