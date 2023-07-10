

USE `test`;

CREATE TABLE `blogs` (
	`id` BIGINT(20) NOT NULL DEFAULT '0',
	`user_id` INT(11) NOT NULL DEFAULT '0',
	`title` VARCHAR(100) NOT NULL DEFAULT '',
	`content` VARCHAR(100) NOT NULL DEFAULT '',
	`status` INT(11) NOT NULL DEFAULT '0',
	`readed` INT(11) NOT NULL DEFAULT '0',
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` BIGINT(20) NOT NULL DEFAULT '0',
	PRIMARY KEY(`id`,`user_id`),
	UNIQUE KEY `uniq_title_of_blog_uk` (`title`),
	KEY `idx_status` (`status`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT 'blogs';

