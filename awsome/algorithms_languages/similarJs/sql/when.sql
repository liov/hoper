--简单Case函数
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
--Case搜索函数
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
= 1;

--这两种方式，可以实现相同的功能。简单Case函数的写法相对比较简洁，但是和Case搜索函数相比，功能方面会有些限制，比如写判断式。
--还有一个需要注意的问题，Case函数只返回第一个符合条件的值，剩下的Case部分将会被自动忽略。
--比如说，下面这段SQL，你永远无法得到“第二类”这个结果
CASE WHEN col_1 IN ( 'a', 'b') THEN '第一类'
WHEN col_1 IN ('a')       THEN '第二类'
ELSE'其他' END