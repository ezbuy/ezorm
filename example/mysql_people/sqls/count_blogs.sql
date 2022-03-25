
SELECT
  COUNT(1)
FROM
  test_user u
JOIN
  blog b ON u.user_id=b.blog_id
WHERE
  u.name = ?
