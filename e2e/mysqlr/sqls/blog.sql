SELECT
	Id ,
	SUM(`title`) AS title_count ,
	status
FROM
	blogs
WHERE
	id = 1;
