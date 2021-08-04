USE `test_sql`;

CREATE TABLE IF NOT EXISTS `user`(
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(1000) DEFAULT NULL,
	`phone` varchar(20) DEFAULT NULL,
	`password` varchar(50) DEFAULT NULL,

	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `user_detail`(
	`id` int NOT NULL AUTO_INCREMENT,
	`user_id` int NOT NULL,
	`email` varchar(100) DEFAULT NULL,
	`text` varchar(1000) DEFAULT NULL,

	PRIMARY KEY (`id`),
	UNIQUE KEY `user_detail_user_id` (`user_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `user_role`(
	`id` int NOT NULL AUTO_INCREMENT,
	`user_id` int NOT NULL,
	`role_id` int NOT NULL,

	PRIMARY KEY (`id`),
	KEY `user_role_key` (`user_id`, `role_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `role`(
	`id` int NOT NULL AUTO_INCREMENT,
	`name` varchar(500) DEFAULT NULL,

	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

