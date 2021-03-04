-- rownum

SET @rownum := 0;
SELECT
    41 AS field_id,
    sale_unit field_value_desc,
    @rownum := @rownum + 1 rownum
FROM
    ( SELECT sale_unit FROM third_sku_info GROUP BY sale_unit ) a