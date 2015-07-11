package dbfav

import (
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"../database"
)

type Fav interface {
	Add(int, int) bool
}

type DBFav struct {
	dbobject database.DB
}

func (self *DBFav) Initialize() {
	myDatabase := &database.Database{}
	var myDb database.DB = myDatabase
	self.dbobject = myDb
}

func NewDBFav() *DBFav {
	dbfav := &DBFav{}
	dbfav.Initialize()
	return dbfav
}

func (u *DBFav) Add(user_id int, youtube_movie_id int) bool {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into youtube_movie_favs (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}
	return true
}
