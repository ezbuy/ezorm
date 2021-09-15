
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object Blog.
CREATE INDEX `idx_blog_user_is_published` ON `blog`(`user`, `is_published`);
CREATE UNIQUE INDEX `idx_blog_slug` ON `blog`(`slug`);
CREATE INDEX `idx_blog_user` ON `blog`(`user`);
CREATE INDEX `idx_blog_is_published` ON `blog`(`is_published`);
CREATE INDEX `idx_blog_group_id` ON `blog`(`group_id`);
CREATE INDEX `idx_blog_create` ON `blog`(`create`);
CREATE INDEX `idx_blog_update` ON `blog`(`update`);


-- DDL for object User.
CREATE TABLE `test_user` (
  `user_id` INT NOT NULL DEFAULT 0,
  `user_number` INT NOT NULL DEFAULT 0,
  `name` VARCHAR(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object User.
CREATE UNIQUE INDEX `idx_test_user_user_number` ON `test_user`(`user_number`);
CREATE INDEX `idx_test_user_name` ON `test_user`(`name`);

