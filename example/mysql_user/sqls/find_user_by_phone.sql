-- 根据phone来查找用户

SELECT
  user_id, name, phone, age, balance, text, create_date
FROM user
WHERE phone=?
