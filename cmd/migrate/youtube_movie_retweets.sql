CREATE TABLE IF NOT EXISTS youtube_movie_retweets (
id int(11) NOT NULL AUTO_INCREMENT,
user_id int(11),
youtube_movie_id int(11),
created_at datetime DEFAULT NULL,
updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
