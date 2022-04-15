
-- DDL for object User.
CREATE TABLE `users` (
  `id` INT NOT NULL DEFAULT 0,
  `name` VARCHAR(200) NOT NULL DEFAULT '',
  `mailbox` VARCHAR(200) NOT NULL DEFAULT '',
  `sex` TINYINT NOT NULL DEFAULT 0,
  `longitude` DECIMAL(11, 4) NOT NULL DEFAULT '0.00',
  `latitude` DECIMAL(11, 4) NOT NULL DEFAULT '0.00',
  `description` VARCHAR(200) NOT NULL DEFAULT '',
  `password` VARCHAR(200) NOT NULL DEFAULT '',
  `head_url` VARCHAR(200) NOT NULL DEFAULT '',
  `status` INT NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object User.
CREATE INDEX `idx_users_mailbox_password` ON `users`(`mailbox`, `password`);
CREATE INDEX `idx_users_name` ON `users`(`name`);


-- DDL for object Blog.
CREATE TABLE `blogs` (
  `id` INT NOT NULL DEFAULT 0,
  `user_id` INT NOT NULL DEFAULT 0,
  `title` VARCHAR(200) NOT NULL DEFAULT '',
  `content` VARCHAR(200) NOT NULL DEFAULT '',
  `status` INT NOT NULL DEFAULT 0,
  `readed` INT NOT NULL DEFAULT 0,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object Blog.
CREATE INDEX `idx_blogs_user_id` ON `blogs`(`user_id`);

