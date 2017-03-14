CREATE TABLE test.Blog
(
    blog_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    title VARCHAR(1000),
    hits INT(11),
    slug VARCHAR(20),
    body TEXT,
    user INT(11),
    is_published TINYINT(1),
    `create` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `update` DATETIME,
    column_10 TIMESTAMP,
    time_stamp BIGINT(20)
);
CREATE UNIQUE INDEX Blog_Slug_uindex ON test.Blog (slug);

CREATE TABLE test.test_user
(
    user_id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    user_number INT(11)
);
CREATE UNIQUE INDEX test_user_user_number_uindex ON test.test_user (user_number);
