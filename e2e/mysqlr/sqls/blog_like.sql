SELECT
	b.id,
	b.status
FROM
	blogs AS b
WHERE
	b.title LIKE 'ezorm%'
