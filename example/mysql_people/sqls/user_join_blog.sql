SELECT
	u.user_id,
	b.blog_id
FROM
	test_user u INNER JOIN blog b ON u.user_id = b.user
WHERE
	u.name = "me"
