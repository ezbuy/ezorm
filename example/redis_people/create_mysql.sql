
-- DDL for object Blog.
CREATE TABLE `blogs` (
  `id` INT NOT NULL,
  `user_id` INT NOT NULL,
  `title` VARCHAR(200) NOT NULL,
  `content` VARCHAR(200) NOT NULL,
  `status` INT NOT NULL,
  `readed` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object Blog.
CREATE INDEX `idx_blogs_user_id` ON `blogs`(`user_id`);


-- DDL for object User.
CREATE TABLE `users` (
  `id` INT NOT NULL,
  `name` VARCHAR(200) NOT NULL,
  `mailbox` VARCHAR(200) NOT NULL,
  `sex` BIT NOT NULL,
  `longitude` DECIMAL(11, 4) NOT NULL,
  `latitude` DECIMAL(11, 4) NOT NULL,
  `description` VARCHAR(200) NOT NULL,
  `password` VARCHAR(200) NOT NULL,
  `head_url` VARCHAR(200) NOT NULL,
  `status` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL,
  `updated_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object User.
CREATE INDEX `idx_users_mailbox_password` ON `users`(`mailbox`, `password`);
CREATE INDEX `idx_users_name` ON `users`(`name`);

