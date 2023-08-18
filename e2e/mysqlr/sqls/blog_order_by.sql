SELECT
	b.user_id,
	b.status
FROM
	blogs AS b
WHERE
	b.user_id > 0
ORDER BY b.id DESC
LIMIT 1,2;
