-- 简单Case函数
SELECT s.s_id,
       s.s_name,
       s.s_sex,
       CASE
           WHEN s.s_sex = '1' THEN '男'
           WHEN s.s_sex = '2' THEN '女'
           ELSE '其他'
           END as sex,
       s.s_age,
       s.class_id
FROM t_b_student s;
-- Case搜索函数
SELECT s.s_id,
       s.s_name,
       s.s_sex,
       CASE s.s_sex
           WHEN '1' THEN '男'
           WHEN '2' THEN '女'
           ELSE '其他'
           END as sex,
       s.s_age,
       s.class_id
FROM t_b_student s;

-- 这两种方式，可以实现相同的功能。简单Case函数的写法相对比较简洁，但是和Case搜索函数相比，功能方面会有些限制，比如写判断式。
-- 还有一个需要注意的问题，Case函数只返回第一个符合条件的值，剩下的Case部分将会被自动忽略。
-- 比如说，下面这段SQL，你永远无法得到“第二类”这个结果
CASE WHEN col_1 IN ( 'a', 'b') THEN '第一类'
WHEN col_1 IN ('a')       THEN '第二类'
ELSE'其他' END

-- 性能极差，拆分
-- CASE WHEN 语句中的条件过滤不会用索引，所以在WHERE中要先过滤一遍，否则会全表扫描
SELECT  count(case when order_status = 8 then '1' end) as wait_delivery_count,
count(case when order_status = 18 and supplier_id <> 19 then '1' end) as abnormal_order,
count(case when order_status = 18 and supplier_id = 19 then '1' end) as abnormal_booking_order,
count(case when service_status = 1 then '1' end) as wait_refund_count FROM order_info
WHERE operator_id = 10001 AND status = 0;

SELECT COUNT(*) FROM order_info WHERE order_status = 8 AND operator_id = 100001 AND status = 0;

SELECT COUNT(*) FROM order_info
WHERE order_status = 18 AND supplier_id <> 19 AND operator_id = 100001 AND status = 0;

SELECT COUNT(*) FROM order_info
WHERE order_status = 18 AND supplier_id = 19 AND operator_id = 100001 AND status = 0;

SELECT COUNT(*) FROM service_info WHERE service_status = 3 AND operator_id = 100001;