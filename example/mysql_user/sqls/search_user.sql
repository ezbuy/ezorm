-- 对用户进行模糊查询

SELECT user_id, name, phone, age
FROM user
WHERE name LIKE ?
LIMIT ?, ?
