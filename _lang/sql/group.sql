SELECT name, COUNT(*)
FROM employee_tbl
GROUP BY name;

SELECT name, SUM(singin) as singin_count
FROM employee_tbl
GROUP BY name
WITH ROLLUP;

SELECT a.plan_order_id,
       a.employee_account,
       a.employee_name,
       a.department,
       a.position,
       a.email,
       a.mobile,
       a.amount                               as distribute_score,
       a.amount - sum(c.total_balance_amount) as exchange_score,
       sum(available_balance_amount)          as available_score
FROM d_kaop.t_plan_detail a
         inner join d_kaop.t_operator_employee b on a.employee_account = b.employee_account
         right join d_aura_jike.user_balance c on b.user_id = c.user_id
WHERE (a.plan_order_id = '20200330163927629198' AND
       c.batch_num IN ('00062020033017021296140', '00062020033017021268100')
    AND a.is_deleted = 0 AND b.is_deleted = 0 AND c.status = 0)
GROUP BY b.user_id

-- 根据日期
-- 1. 找到当前表中的日期列，并且其转换成需要排序的年月格式便可，并且取出对应的字符长度。

-- 2. 如下，我需要将金额数据按照月度汇总，那么我需要做的就是把当前日期先转换成年月格式的日期，然后按照分组。

-- 3. 需要注意的是，需要将group后的日期字段和查询列的字段都转换为年月格式的字符。如 2019-05 。


select id,CONVERT(varchar(7),date,120) group by CONVERT(varchar(7),date,120)


--  4. 常用的日期格式有以下：


select CONVERT(varchar(12) , getdate(), 101 );
 --   05/19/2019

select CONVERT(varchar(12) , getdate(), 102 );
  --  2019.05.19

select CONVERT(varchar(12) , getdate(), 103 );
  --  19/05/2019

select CONVERT(varchar(12) , getdate(), 104 );
  --  19.05.2019

select CONVERT(varchar(12) , getdate(), 105 );
  --  19-05-2019

select CONVERT(varchar(12) , getdate(), 106 );
  --  19 05 2019
  --  ---------------------------------------------
select CONVERT(varchar(12) , getdate(), 107 );
   -- 05 19, 2019

select CONVERT(varchar(12) , getdate(), 108 );
  --  09:15:33

select CONVERT(varchar(12) , getdate(), 109 );
  --  05 19 2019

select CONVERT(varchar(12) , getdate(), 110 );
   -- 05-19-2019

select CONVERT(varchar(12) , getdate(), 111 );
   -- 2019/05/19

select CONVERT(varchar(12) , getdate(), 112 );
   -- 20190519

select CONVERT(varchar(12) , getdate(), 113 );
  --  19 05 2019 0

select CONVERT(varchar(12) , getdate(), 114 );
  --  09:16:06:747


-- 5. 如果需要去日期中相关的值，有以下方法
#
#     YEAR('2018-05-17 00:00:00'); -- 年
#     MONTH('2018-05-15 00:00:00'); -- 月
#     DAY('2008-05-15 00:00:00'); -- 日
#     DATEPART ( datepart , date );
#     DATEPART(MM,'2018-05-15 00:00:00');
#    年份 yy、yyyy
#    季度 qq、q
#    月份 mm、m
#    每年的某一日 dy、y
#    日期 dd、d
#    星期 wk、ww
#    工作日 dw
#    小时 hh
#    分钟 mi、n
#    秒 ss、s
#    毫秒 ms

--  EXTRACT (component_name, FROM {datetime | interval})
-- GROUP BY GROUP_CONCAT 拼某一列
SELECT t.sid,t.name,t.sex,GROUP_CONCAT(t.num) from distinct_concat t GROUP BY t.sid,t.name,t.sex;

-- order by 在 group by之后执行，要保留第一行要做子查询 实测mysql 要加LIMIT bignum
SELECT * FROM (SELECT * FROM tsp_settle_info WHERE effect_date <= now() ORDER BY effect_date DESC LIMIT 10000000) GROUP BY tsp_id
