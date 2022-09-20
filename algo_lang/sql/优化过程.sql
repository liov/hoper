SELECT SUM(b.sku_count) AS count
FROM order_info a
         LEFT JOIN order_sku b ON a.order_id = b.order_id AND b.status = 0
         LEFT JOIN order_sku_epay f ON f.order_id = a.order_id AND f.status = 0
         LEFT JOIN product_category c ON b.product_id = c.product_id AND c.status = 0
         INNER JOIN (SELECT d1.id AS parent_id, d1.category_name, d2.id
                     FROM category_info d1
                              LEFT JOIN category_info d2
                                        ON d1.id = d2.parent_id AND d2.status = 0
                     Where d1.id IN (9, 11, 12, 19, 23, 28, 92, 136, 334, 357, 7, 36, 504, 443, 261, 190, 115)
                       AND d1.status = 0) d ON c.category_id = d.id
WHERE (a.status = 0 AND a.operator_id = 0 AND a.order_status = 13 AND f.epay_type = 2);

-- 自连接再左连接改为双左连接
-- 6s -- 98ms
SELECT SUM(b.sku_count) AS count
FROM order_info a
         LEFT JOIN order_sku b ON a.order_id = b.order_id AND b.status = 0
         LEFT JOIN order_sku_epay f ON f.order_id = a.order_id AND f.status = 0
         LEFT JOIN product_category c ON b.product_id = c.product_id AND c.status = 0
         LEFT JOIN category_info d ON c.category_id = d.id AND d.`status` = 0
WHERE d.parent_id IN (9, 11, 12, 19, 23, 28, 92, 136, 334, 357, 7, 36, 504, 443, 261, 190, 115)
  AND (a.status = 0 AND a.operator_id = 0 AND a.order_status = 13 AND f.epay_type = 2);


-- 8s -- 96ms
SELECT d2.category_name AS name, SUM(b.sku_count) AS count
FROM order_info a
         LEFT JOIN order_sku b ON a.order_id = b.order_id AND b.status = 0
         LEFT JOIN order_sku_epay f ON f.order_id = a.order_id AND f.status = 0
         LEFT JOIN product_category c ON b.product_id = c.product_id AND c.status = 0
         LEFT JOIN category_info d ON c.category_id = d.id AND d.status = 0
         LEFT JOIN category_info d2 ON d.parent_id = d2.id AND d2.status = 0
WHERE (a.status = 0 AND a.operator_id = 0 AND a.order_status = 13 AND f.epay_type = 2)
  AND (d.parent_id IN (9, 11, 12, 19, 23, 28, 92, 136, 334, 357, 7, 36, 504, 443, 261, 190, 115))
GROUP BY d2.id
ORDER BY count DESC
LIMIT 10;
