-- 根据电话号码查询某个用户的详细信息

SELECT
  u.user_id id,
  u.name,
  u.phone,
  u.age,
  u.balance,
  u.text,
  u.create_date,

  IFNULL(ud.score, 0),
  IFNULL(ud.balance, 0),
  IFNULL(ud.text, '') detail_text

FROM user u
JOIN user_detail ud ON u.user_id=ud.user_id
WHERE u.phone LIKE ?
