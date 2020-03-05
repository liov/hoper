count (表达式)--分组里非空记录数
count (表达式)--分组里非空记录数
count(*)--所有记录
count(1)--所有记录
count(case job = 'CLERK' then 2 end )--CLERK 人数
count(comm)--有奖金的人数
count(distinct job)--distinct(去重），共有多少种工作

select deptno,
count(1) 总人数,
count(case when job ='SALESMAN' then '1' end) 销售人数,
count(case when job ='MANAGER' then '1' end) 主管人数
from emp
group by deptno;--如果不group，会认为所有数据是一组，返回一个数据