-- 查询某个用户的博客

SELECT
  b.blog_id ID,
  b.title,
  b.hits,
  b.slug,
  IFNULL(b.body, ''),
  IFNULL(b.is_published, 0) published,
  b.group_id,
  b.create,
  b.update,
  u.user_id,
  u.user_number,
  u.name
FROM
  test_user u
JOIN
  blog b ON u.user_id=b.blog_id
WHERE
  u.name = ?
LIMIT ?, ?
