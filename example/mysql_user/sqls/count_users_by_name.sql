-- 根据名称统计用户数量

SELECT
  -- 对于统计字段，如果没有设置别名，将会使用"Count0"这样的名称
  COUNT(1)

FROM user
WHERE name LIKE ?
