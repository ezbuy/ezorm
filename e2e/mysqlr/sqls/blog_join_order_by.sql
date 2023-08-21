SELECT
	u.id AS uid,
	b.id AS bid
FROM
	users  u
INNER JOIN
	blogs  b
ON
	u.id = b.user_id
WHERE
	u.id > 0
ORDER BY u.id DESC
