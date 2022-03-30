-- 根据电话号码统计用户数量

SELECT
  -- 对于统计字段，如果用了别名，将优先使用别名作为字段名称
  COUNT(0) user_count

FROM user
WHERE phone LIKE ?
