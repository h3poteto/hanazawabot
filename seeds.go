package main

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	create_youtube_movies()
	create_serifs()
	create_users()
	create_tweets()
	create_youtube_movie_favs()
	create_youtube_movie_retweets()
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

func create_users() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select users.id from users;")
	if err != nil {
		_, _ = db.Exec("CREATE TABLE users (id int(11) NOT NULL AUTO_INCREMENT, twitter_id bigint(20) unsigned NOT NULL, screen_name varchar(255) DEFAULT NULL, created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	}
}

func create_tweets() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select tweets.id from tweets;")
	if err != nil {
		_, _ = db.Exec("CREATE TABLE tweets (id int(11) NOT NULL AUTO_INCREMENT, user_id int(11), tweet varchar(255) DEFAULT NULL, created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
	}
}

func create_youtube_movie_favs() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select youtube_movie_favs.id from youtube_movie_favs;")
	if err != nil {
		_, err = db.Exec("CREATE TABLE youtube_movie_favs (id int(11) NOT NULL AUTO_INCREMENT, user_id int(11), youtube_movie_id int(11), created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
		if err != nil {
			log.Fatalf("mysql error: %v ", err)
		}
	}
}


func create_youtube_movie_retweets() {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Query("select youtube_movie_retweets.id from youtube_movie_retweets;")
	if err != nil {
		_, err = db.Exec("CREATE TABLE youtube_movie_retweets (id int(11) NOT NULL AUTO_INCREMENT, user_id int(11), youtube_movie_id int(11), created_at datetime DEFAULT NULL, updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;")
		if err != nil {
			log.Fatalf("mysql error: %v ", err)
		}
	}
}
