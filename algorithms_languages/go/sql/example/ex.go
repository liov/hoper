package example

import "github.com/jinzhu/gorm"

var db *gorm.DB

func example() {
	var plans []string
	var planOrderId string
	db.Table("d_kaop.t_plan_detail a").
		Select(`a.plan_order_id,
	a.employee_account,
	b.employee_name,
	b.mobile,
	b.email,
	b.department,
	b.position,
	a.amount AS distribute_score,
	b.user_id,
	a.amount - COALESCE(sum(case when c.balance_status = 0 then c.total_balance_amount ELSE 0 end),0) + COALESCE(sum(c.invalid),0) AS exchange_score,
	sum(case when c.balance_status = 0 then c.available_balance_amount ELSE 0 end) AS available_score `).
		Joins("LEFT JOIN d_kaop.t_operator_employee b ON a.employee_account = b.employee_account AND a.operator_id = b.operator_id").
		Joins(`	LEFT JOIN (SELECT
	d.total_balance_amount,
	d.available_balance_amount,
	d.balance_status,
	e.invalid,
	d.user_id,
	d.batch_num 
FROM
	d_aura_jike.user_balance d
	LEFT JOIN (
	SELECT
		sum( balance_amount ) AS invalid,
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
	AND d.user_id = e.user_id  WHERE d.STATUS = 0 AND d.batch_num IN (?)) c ON b.user_id = c.user_id`,
			db.Table("d_kaop.t_plan_score_batch").Select("batch_no").
				Where("plan_order_id = ? AND is_deleted = 0", planOrderId).QueryExpr()).
		Where(`a.plan_order_id = ? AND a.is_deleted = 0`, planOrderId).
		Group("b.user_id").Scan(&plans)
}
