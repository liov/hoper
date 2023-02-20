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
	);

-- 每个班最高分
select 班级名称,分数=MAX(成绩)
from
学生表 s join 班级表 c
on s.班级号=c.班级号
join 成绩表 sc
on s.学号=sc.学号
group by 班级名称

-- 一、按name分组取val最大的值所在行的数据。
-- 方法1：
select a.* from tb a where val = (select max(val) from tb where name = a.name) order by a.name
-- 方法2：
select a.* from tb a where not exists(select 1 from tb where name = a.name and val > a.val)
-- 方法3：
select a.* from tb a,(select name,max(val) val from tb group by name) b where a.name = b.name and a.val = b.val order by a.name
-- 方法4：
select a.* from tb a inner join (select name , max(val) val from tb group by name) b on a.name = b.name and a.val = b.val order by a.name
-- 方法5
select a.* from tb a where 1 > (select count(*) from tb where name = a.name and val > a.val ) order by a.name
-- Group By 分组后保留最新一条记录
select a.* from
(select * from user order by id desc) a
group by a.id;

# 查找几个字段中的最大值 UNION

SELECT
	customer_id,
	Max( last_visit_time )
FROM
	(
	SELECT
		id AS customer_id,
		created_at AS last_visit_time
	FROM
		d_crm_customer.customer_info UNION ALL
	SELECT
		customer_id,
		visit_time AS last_visit_time
	FROM
		d_crm_customer.customer_visit UNION ALL
	SELECT
		customer_id,
		sign_time AS last_visit_time
	FROM
		d_crm_sales.sign_info
	) a
GROUP BY
	customer_id
	LIMIT 1000 OFFSET 0;


# count

-- count (表达式); -- 分组里非空记录数
-- count (表达式); -- 分组里非空记录数
-- count(*); -- 所有记录
-- count(1); -- 所有记录
-- count(case job = 'CLERK' then 2 end ); -- CLERK 人数
-- count(comm); -- 有奖金的人数
-- count(distinct job); -- distinct(去重），共有多少种工作

SELECT `level`,count(*) FROM `customer` GROUP BY `level` HAVING `level`=1 OR `level` =2 OR `level`=3;

SELECT `level`,count(id) FROM `customer` WHERE `level`=1 OR `level` =2 OR `level`=3 GROUP BY `level`=1 OR `level` =2 ,`level`=3 ORDER BY `level`;

SELECT COUNT(case when (a.level=1 OR a.level=2) then level end) as focusNum,COUNT(case when a.level=3 then level end) as totalNum FROM customer a;

SELECT a.id,sum(b.contract_number),DATE_FORMAT(create_time,'%Y-%m-%d') AS dateTime FROM `trade` a,`trade_contract` b WHERE a.id = b.trade_id GROUP BY DATE_FORMAT(create_time,'%Y-%m-%d');

SELECT
	COUNT( customerNum ) / ( SELECT count( id ) FROM `customer` WHERE `level` = 4 AND create_time BETWEEN "2017-11-30T16:00:00.000Z" AND "2018-12-13T16:00:00.000Z" ) AS rate
FROM
	(
	SELECT
		count( customer_id ) AS customerNum
	FROM
		`follow` a LEFT JOIN `customer` b
	WHERE
		a.customer_id IN ( SELECT id FROM `customer` d WHERE d.`level` = 4 AND d.create_time BETWEEN "2017-11-30T16:00:00.000Z" AND "2018-12-13T16:00:00.000Z" )
	AND a.customer_level < 4
	AND DATEDIFF(b.create_time,a.create_time) <=15
	GROUP BY a.customer_id
	) c;

# case when

select deptno,
count(1) '总人数',
count(case when job ='SALESMAN' then '1' end) '销售人数',
count(case when job ='MANAGER' then '1' end) '主管人数'
from emp
group by deptno;-- 如果不group，会认为所有数据是一组，返回一个数据

SELECT count(case when order_status =8 then '1' end),
count(case when order_status =9 then '1' end),
count(case when service_status IN (1,2) then '1' end)
FROM `order_info`;

-- 简单Case函数
SELECT
    s.s_id,
    s.s_name,
    s.s_sex,
    CASE
WHEN s.s_sex = '1' THEN '男'
WHEN s.s_sex = '2' THEN '女'
ELSE '其他'
END as sex,
 s.s_age,
 s.class_id
FROM
    t_b_student s
WHERE
= 1;
-- Case搜索函数
SELECT
    s.s_id,
    s.s_name,
    s.s_sex,
    CASE s.s_sex
WHEN '1' THEN '男'
WHEN '2' THEN '女'
ELSE '其他'
END as sex,
 s.s_age,
 s.class_id
FROM
    t_b_student s
WHERE
1 = 1;

-- 这两种方式，可以实现相同的功能。简单Case函数的写法相对比较简洁，但是和Case搜索函数相比，功能方面会有些限制，比如写判断式。
-- 还有一个需要注意的问题，Case函数只返回第一个符合条件的值，剩下的Case部分将会被自动忽略。
-- 比如说，下面这段SQL，你永远无法得到“第二类”这个结果
CASE WHEN col_1 IN ( 'a', 'b') THEN '第一类'
WHEN col_1 IN ('a')       THEN '第二类'
ELSE'其他' END

# GROUP

SELECT name, COUNT(*) FROM   employee_tbl GROUP BY name;

SELECT name, SUM(singin) as singin_count FROM  employee_tbl GROUP BY name WITH ROLLUP;

SELECT a.plan_order_id,a.employee_account,a.employee_name,a.department,a.position,a.email,a.mobile,
a.amount as distribute_score,a.amount - sum(c.total_balance_amount) as exchange_score,sum(available_balance_amount) as available_score FROM d_kaop.t_plan_detail a inner join d_kaop.t_operator_employee b on a.employee_account = b.employee_account right join d_aura_jike.user_balance c on b.user_id = c.user_id WHERE (a.plan_order_id = '20200330163927629198' AND c.batch_num IN ('00062020033017021296140','00062020033017021268100')
AND a.is_deleted = 0 AND b.is_deleted = 0 AND c.status = 0) GROUP BY b.user_id;

# join

-- INNER JOIN（内连接,或等值连接）：获取两个表中字段匹配关系的记录。
-- LEFT JOIN（左连接）：获取左表所有记录，即使右表没有对应匹配的记录。
-- RIGHT JOIN（右连接）： 与 LEFT JOIN 相反，用于获取右表所有记录，即使左表没有对应匹配的记录。
-- full join:外连接，返回两个表中的行：left join + right join。
-- cross join:结果是笛卡尔积，就是第一个表的行数乘以第二个表的行数。

-- 两个表在，join时，首先做一个笛卡尔积，on后面的条件是对这个笛卡尔积做一个过滤形成一张临时表，如果没有where就直接返回结果，如果有where就对上一步的临时表再进行过滤。
--
-- 在使用left  jion时，on和where条件的区别如下：
--
-- 1、on条件是在生成临时表时使用的条件，它不管on中的条件是否为真，都会返回左边表中的记录。
--
-- 2、where条件是在临时表生成好后，再对临时表进行过滤的条件。这时已经没有left  join的含义（必须返回左边表的记录）了，条件不为真的就全部过滤掉


-- join on:on后边写条件，以后边的条件为准生成一个临时表存储数据。

-- on和where的区别：
--
-- 1. 对于left join，不管on后面跟什么条件，左表的数据全部查出来，因此要想过滤需把条件放到where后面
--
-- 2. 对于inner join，满足on后面的条件表的数据才能查出，可以起到过滤作用。也可以把条件放到where后面。


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

--这是三张表的左连接查询；
SELECT
	a.plan_order_id,
	a.employee_account,
	a.employee_name,
	a.amount AS distribute_score,
	b.user_id,
	a.amount - sum( c.total_balance_amount ) + sum( c.invilad ) AS exchange_score,
	sum( c.available_balance_amount ) AS available_score
FROM
	d_kaop.t_plan_detail a
	LEFT JOIN d_kaop.t_operator_employee b ON a.employee_account = b.employee_account
	LEFT JOIN (
	SELECT
		d.total_balance_amount,
		d.available_balance_amount,
		e.invilad,
		d.user_id,
		d.batch_num
	FROM
		d_aura_jike.user_balance d
		LEFT JOIN (
		SELECT
			sum( balance_amount ) AS invilad,
			user_id,
			batch_num
		FROM
			d_aura_jike.user_balance_history
		WHERE
			type IN ( 6, 7 )
		GROUP BY
			user_id,
			batch_num
		) e ON d.batch_num = e.batch_num
		AND d.user_id = e.user_id
	WHERE
		d.STATUS = 0
		AND d.batch_num IN (
		SELECT
			batch_no
		FROM
			`d_kaop`.`t_plan_score_batch`
		WHERE
		( plan_order_id = '20200330163927629198' AND is_deleted = 0 ))
	) c ON b.user_id = c.user_id
WHERE
	( a.plan_order_id = '20200330163927629198' AND a.is_deleted = 0 )
GROUP BY
	a.employee_account;

# FORCE INDEX
SELECT * FROM `product_attr` FORCE INDEX(`idx_product_attr_id`) WHERE  product_id IN (SELECT product_id FROM (SELECT product_id FROM `product_attr` WHERE attr_id = 1103 GROUP BY product_id HAVING COUNT(1)>1) temp) AND attr_id = 1103 AND  created_at > 0;

# template
-- template<T:DATE_FORMAT(create_time,'%Y-%m-%d')|WEEK(create_time)|MONTH(create_time)>
SELECT count(*),T AS dateTime
FROM `customer` a
WHERE a.level = 3 OR a.`level` = 2
GROUP BY T;
