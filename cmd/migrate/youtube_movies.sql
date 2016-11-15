CREATE TABLE IF NOT EXISTS youtube_movies (
id int(11) NOT NULL AUTO_INCREMENT,
title varchar(255) DEFAULT NULL,
movie_id varchar(255) UNIQUE DEFAULT NULL,
description text,
used tinyint(1) NOT NULL DEFAULT 0,
disabled tinyint(1) NOT NULL DEFAULT 0,
created_at datetime DEFAULT NULL,
updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
