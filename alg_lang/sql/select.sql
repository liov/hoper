-- rownum

SET @rownum := 0;
SELECT
    41 AS field_id,
    sale_unit field_value_desc,
    @rownum := @rownum + 1 rownum
FROM
    ( SELECT sale_unit FROM third_sku_info GROUP BY sale_unit ) a

-- 每个班级前十名
select * from 学生信息表 a
where  (select count(*) from 学生信息表 where 班级ID = a.班级ID and 班内名次 > a.班内名次) < 10

-- 每个班第一名
SELECT
	*
FROM
	t_stu_score a
WHERE
	a.score IN (
		SELECT
			MAX(score)
		FROM
			t_stu_score b
		WHERE
			a.class_id = b.class_id
		ORDER BY
			score DESC
	)

-- 每个班最高分
select 班级名称,分数=MAX(成绩)
from
学生表 s join 班级表 c
on s.班级号=c.班级号
join 成绩表 sc
on s.学号=sc.学号
group by 班级名称

--一、按name分组取val最大的值所在行的数据。
--方法1：
select a.* from tb a where val = (select max(val) from tb where name = a.name) order by a.name
--方法2：
select a.* from tb a where not exists(select 1 from tb where name = a.name and val > a.val)
--方法3：
select a.* from tb a,(select name,max(val) val from tb group by name) b where a.name = b.name and a.val = b.val order by a.name
--方法4：
select a.* from tb a inner join (select name , max(val) val from tb group by name) b on a.name = b.name and a.val = b.val order by a.name
--方法5
select a.* from tb a where 1 > (select count(*) from tb where name = a.name and val > a.val ) order by a.name
-- Group By 分组后保留最新一条记录
select a.* from
(select * from user order by id desc) a
group by a.id