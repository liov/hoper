# You can't specify target table 'product_attr' for update in FROM clause
UPDATE customer_info
SET status     = 1,
    created_at = '2020-08-06 00:00:00',
    updated_at = '2020-08-06 00:00:00',
    succeeded  = 1,
    stage      = 2
WHERE id >= (SELECT id FROM (SELECT id FROM `customer_info` WHERE customer_num = '2020') temp);

# upsert
INSERT INTO user_role(id, role_id) VALUES (1, 1) ON DUPLICATE KEY UPDATE role_id = 1;