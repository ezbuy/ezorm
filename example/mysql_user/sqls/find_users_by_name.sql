-- 查询用户的详细信息

SELECT
  -- 生成的返回结构体将会优先使用别名作为字段名称。这里将会使用Id
  u.user_id id,
  -- 不使用别名，使用"表名称+字段名称"作为生成结构体的字段名称
  u.name,
  u.phone,
  u.age,
  u.balance,
  u.text,
  u.create_date,

  -- 可以使用Func，例如IFNULL来防止nil值引起的错误
  IFNULL(ud.score, 0),
  IFNULL(ud.balance, 0),
  -- 别名会被自动转换为驼峰形式，这里会变成"DetailText"字段
  IFNULL(ud.text, '') detail_text

FROM user u
JOIN user_detail ud ON u.user_id=ud.user_id
-- 直接使用占位符表示sql语句的参数
WHERE u.name LIKE ?
LIMIT ?, ?
