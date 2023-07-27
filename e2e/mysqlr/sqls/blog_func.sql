
SELECT
	UPPER(title) as u_title,
	LENGTH(title) as len_title
FROM
	blogs
WHERE
	id = 1;
