

USE `test`;

CREATE TABLE `users` (
	`id` BIGINT(20) NOT NULL DEFAULT '0',
	`user_id` INT(11) NOT NULL DEFAULT '0',
	`name` VARCHAR(100) NOT NULL DEFAULT '',
	`created_at` BIGINT(20) NOT NULL DEFAULT '0',
	`updated_at` BIGINT(20) NOT NULL DEFAULT '0',
	PRIMARY KEY(`id`,`user_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'users';

