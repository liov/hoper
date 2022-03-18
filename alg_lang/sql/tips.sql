-- 字符串加前缀
update mimvp set name=concat('米扑科技 - ', name) where name like '米扑%';
-- left(str, len) 返回字符串str的最左面len个字符。
select left('mimvp.com',5);
-- 删除前缀
update mimvp set name = substring(name, 8) where name like '米扑科技 - %';
UPDATE post SET path = substring(path,27) WHERE t_id > 461000 AND path like '\\F\\%'