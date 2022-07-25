
-- DDL for object Blog.
CREATE TABLE `blog` (
  `blog_id` INT NOT NULL DEFAULT 0,
  `title` VARCHAR(200) NOT NULL DEFAULT '',
  `hits` INT NOT NULL DEFAULT 0,
  `slug` VARCHAR(200) NOT NULL DEFAULT '',
  `body` VARCHAR(200),
  `user` INT NOT NULL DEFAULT 0,
  `is_published` TINYINT NOT NULL DEFAULT 0,
  `group_id` BIGINT NOT NULL DEFAULT 0,
  `create` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`blog_id`)
  KEY `idx_blog_user_is_published` (`user`, `is_published`)
  UNIQUE KEY `idx_blog_slug` (`slug`)
  KEY `idx_blog_user` (`user`)
  KEY `idx_blog_is_published` (`is_published`)
  KEY `idx_blog_group_id` (`group_id`)
  KEY `idx_blog_create` (`create`)
  KEY `idx_blog_update` (`update`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';


-- DDL for object User.
CREATE TABLE `test_user` (
  `user_id` INT NOT NULL DEFAULT 0,
  `user_number` INT NOT NULL DEFAULT 0,
  `name` VARCHAR(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`)
  UNIQUE KEY `idx_test_user_user_number` (`user_number`)
  KEY `idx_test_user_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

