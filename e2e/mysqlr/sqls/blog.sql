SELECT
	Id ,
	SUM(`title`) AS title_count ,
	status
FROM
	blogs
WHERE
	id = 1
LIMIT 1,2;
