

USE `test`;

CREATE TABLE `auto_blogs` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`user_id` INT(11) NOT NULL DEFAULT '0',
	`title` VARCHAR(100) NOT NULL DEFAULT '',
	`content` VARCHAR(100) NOT NULL DEFAULT '',
	`status` INT(11) NOT NULL DEFAULT '0',
	`readed` INT(11) NOT NULL DEFAULT '0',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` BIGINT(20) NOT NULL DEFAULT '0',
	PRIMARY KEY(`id`),
	KEY `status_of_auto_blog_idx` (`status`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'auto_blogs';

