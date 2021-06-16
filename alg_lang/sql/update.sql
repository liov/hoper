UPDATE customer_info
SET status     = 1,
    created_at = '2020-08-06 00:00:00',
    updated_at = '2020-08-06 00:00:00',
    succeeded  = 1,
    stage      = 2
WHERE id >= (SELECT id FROM (SELECT id FROM `customer_info` WHERE customer_num = '2020') temp);
