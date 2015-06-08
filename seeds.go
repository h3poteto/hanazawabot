package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	create_youtube_movies()
	create_serifs()
}

func create_youtube_movies() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select youtube_movies.id from youtube_movies;")
	if err != nil {
		_, _ = db.Exec("CREATE TABLE youtube_movies (id int(11) NOT NULL AUTO_INCREMENT, title varchar(255) DEFAULT NULL, movie_id varchar(255) DEFAULT NULL, description text, disabled tinyint(1) NOT NULL DEFAULT 0, created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	}
}

func create_serifs() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select serifs.id from serifs;")
	if err != nil {
		_, _ = db.Exec("CREATE TABLE serifs (id int(11) NOT NULL AUTO_INCREMENT, body varchar(255) DEFAULT NULL, created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	}

	_, err = db.Exec("truncate table serifs;")
	_, err = db.Exec("insert into serifs (body, created_at) values (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now()), (?, now())",
		"うん、うん、これですよ。このふともも・・・",
		"オチ？オチなんてないよぉ！",
		"私のこの楽しい話を邪魔しないでくれるっ！",
		"あたしもう25歳だよ！お前あたしと結婚できんのか！",
		"そ、そんな事言うと脱ぐぞ",
		"ディスってんのかぁぁっ",
		"ねぇじゃねーよ！",
		"弟紹介しようか？")
	if err != nil {
		log.Fatalf("mysql connect error: %v", err)
	}
}
