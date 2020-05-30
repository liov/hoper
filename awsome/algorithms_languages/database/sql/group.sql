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