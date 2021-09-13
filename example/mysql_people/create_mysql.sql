
-- DDL for object Blog.
CREATE TABLE `blog` (
  `blog_id` INT NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `hits` INT NOT NULL,
  `slug` VARCHAR(200) NOT NULL,
  `body` VARCHAR(200),
  `user` INT NOT NULL,
  `is_published` BIT NOT NULL,
  `group_id` BIGINT NOT NULL,
  `create` TIMESTAMP NOT NULL,
  `update` DATETIME NOT NULL,
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
  `user_id` INT NOT NULL,
  `user_number` INT NOT NULL,
  `name` VARCHAR(200) NOT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object User.
CREATE UNIQUE INDEX `idx_test_user_user_number` ON `test_user`(`user_number`);
CREATE INDEX `idx_test_user_name` ON `test_user`(`name`);

