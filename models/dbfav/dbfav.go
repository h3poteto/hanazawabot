package dbfav

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"../basedb"
)

type Fav interface {
	Add(int, int) error
}

type DBFav struct {
	dbobject basedb.DB
}

func (self *DBFav) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBFav() *DBFav {
	dbfav := &DBFav{}
	dbfav.Initialize()
	return dbfav
}

func (u *DBFav) Add(user_id int, youtube_movie_id int) error {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into youtube_movie_favs (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}
	return nil
}
