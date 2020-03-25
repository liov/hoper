--INNER JOIN（内连接,或等值连接）：获取两个表中字段匹配关系的记录。
--LEFT JOIN（左连接）：获取左表所有记录，即使右表没有对应匹配的记录。
--RIGHT JOIN（右连接）： 与 LEFT JOIN 相反，用于获取右表所有记录，即使左表没有对应匹配的记录。

SELECT count(*) as t_count,
count(case when field1 - field2 <= field3 then '1' end) as w_count,
count(case when field1 = field2 then '1' end) as f_count FROM `table2`
left join table1 on table2.table2_id = table1.id
AND table1.status = 0 AND table1.field4 = 0;

SELECT count(*) as t_count,
count(case when field1 - field2 <= field3 then '1' end) as w_count,
count(case when field1 = field2 then '1' end) as f_count FROM `table2`
right join table1 on table2.table2_id = table1.id
AND table1.status = 0 AND table1.field4 = 0;

SELECT count(*) as t_count;
SELECT count(case when table1.field4 = 0 then '1' end) as t_count;