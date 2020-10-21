USE `test`;

CREATE TABLE IF NOT EXISTS `blog` (
  `blog_id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(1000) COLLATE utf8mb4_bin DEFAULT NULL,
  `hits` int DEFAULT NULL,
  `slug` varchar(20) COLLATE utf8mb4_bin DEFAULT NULL,
  `body` text COLLATE utf8mb4_bin,
  `user` int DEFAULT NULL,
  `is_published` tinyint(1) DEFAULT NULL,
  `create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update` datetime DEFAULT NULL,
  `column_10` timestamp NULL DEFAULT NULL,
  `time_stamp` bigint DEFAULT NULL,
  PRIMARY KEY (`blog_id`),
  UNIQUE KEY `Blog_Slug_uindex` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `test_user` (
  `user_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL,
  `user_number` int DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `test_user_user_number_uindex` (`user_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

