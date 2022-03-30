
-- DDL for object User.
CREATE TABLE `user` (
  `user_id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL DEFAULT 'DefaultName',
  `phone` VARCHAR(20) NOT NULL DEFAULT '',
  `age` INT NOT NULL DEFAULT 3,
  `balance` DECIMAL(20, 4) NOT NULL DEFAULT '0.00',
  `text` VARCHAR(400) DEFAULT 'Hello, user!',
  `create_date` BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object User.
CREATE INDEX `idx_user_name_phone` ON `user`(`name`, `phone`);
CREATE INDEX `idx_user_create_date` ON `user`(`create_date`);


-- DDL for object UserDetail.
CREATE TABLE `user_detail` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL DEFAULT 0,
  `score` INT NOT NULL DEFAULT 0,
  `balance` INT NOT NULL DEFAULT 0,
  `text` VARCHAR(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '';

-- Indexes for object UserDetail.
CREATE INDEX `idx_user_detail_user_id` ON `user_detail`(`user_id`);

