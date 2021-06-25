count (表达式)-- 分组里非空记录数
count(*)-- 所有记录
count(1)-- 所有记录
count(case job = 'CLERK' then 2 end ) -- CLERK 人数
count(comm) -- 有奖金的人数
count(distinct job) -- distinct(去重），共有多少种工作

select deptno,
count(1) 总人数,
count(case when job ='SALESMAN' then '1' end) 销售人数,
count(case when job ='MANAGER' then '1' end) 主管人数
from emp
group by deptno;--如果不group，会认为所有数据是一组，返回一个数据

SELECT count(case when order_status =8 then '1' end),
count(case when order_status =9 then '1' end),
count(case when service_status IN (1,2) then '1' end)
FROM `order_info`;

template<T:DATE_FORMAT(create_time,'%Y-%m-%d')|WEEK(create_time)|MONTH(create_time)>
SELECT count(*),T AS dateTime
FROM `customer` a
WHERE a.level = 3 OR a.`level` = 2
GROUP BY T;

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
