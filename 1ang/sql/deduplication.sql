-- 去重
UPDATE report_receivers
SET is_deleted = 1
WHERE report_id IN (SELECT report_id FROM (SELECT report_id FROM `report_receivers` GROUP BY report_id, emp_id HAVING COUNT(*) > 1) a)
AND id NOT IN (SELECT id FROM (SELECT id FROM `report_receivers` GROUP BY report_id, emp_id) b)