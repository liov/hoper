INSERT INTO `customer_extra_info`(customer_id, last_visit_time)
SELECT customer_id,
       MAX(last_visit_time)
FROM (
         SELECT id         AS customer_id,
                created_at AS last_visit_time
         FROM `customer_info`
         UNION ALL
         SELECT customer_id,
                visit_time AS last_visit_time
         FROM `customer_visit`
         UNION ALL
         SELECT customer_id,
                sign_time AS last_visit_time
         FROM `d_crm_sales`.`sign_info`
     ) a
GROUP BY customer_id;

CREATE
    DEFINER = `web`@`%` PROCEDURE `insert`()
BEGIN
    declare i int;
    set i = 6001;
    while i < 7001
        do
            insert into customer_erptask(customer_id) values (i);
            insert into customer_extra_info(customer_id) values (i);
            set i = i + 1;
        end while;

END;

INSERT INTO `d_sales_support`.`common_dict_value`
(`business_id`, `business_name`, `business_value`, `business_desc`, `status`)
VALUES (5, '开放平台多品类集合', '21', '视频会员', 0),
       (5, '开放平台多品类集合', '24', '糕点', 0),
       (5, '开放平台多品类集合', '25', '饮品', 0),
       (5, '开放平台多品类集合', '26', '鲜花', 0),
       (5, '开放平台多品类集合', '27', '知识付费', 0)
;


INSERT INTO `d_aura_jike`.`sp_field_mapper`(`category_id`,`supplier_id`,`product_type`,`field_type`,`field_id`) VALUES(@s_id,@c_id,8,13,7);
SET @pid:=last_insert_id();
INSERT INTO `d_aura_jike`.`sp_field_mapper_value`(`field_mapper_id`,`field_value`,`field_display_value`) VALUES(@pid,19,'');
-- 插入如果不存在
INSERT INTO TABLE (field1, field2, fieldn) SELECT 'field1','field2','fieldn' WHERE NOT EXISTS (SELECT field FROM TABLE WHERE field = ?)

INSERT INTO brand_info (chinese_name, supplier_id, source) SELECT '测试牌', 11,     2
WHERE
    NOT EXISTS (
            SELECT
                chinese_name
            FROM
                brand_info
            WHERE
                    chinese_name = '测试牌')