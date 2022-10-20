-- 如果想获得所有表的行数信息，可以使用以下 SQL 语句
SELECT
    relname,
    reltuples
FROM
    pg_class
        CLS LEFT JOIN pg_namespace N ON ( N.oid = CLS.relnamespace )
WHERE
        nspname NOT IN ( 'pg_catalog', 'information_schema' )
  AND relkind = 'r'
ORDER BY
    reltuples DESC;

-- 更精确的计算方法是创建一个函数来实现统计功能：
CREATE TYPE table_count AS (table_name TEXT, num_rows INTEGER);


CREATE OR REPLACE FUNCTION count_em_all () RETURNS SETOF table_count  AS '
DECLARE
    the_count RECORD;
    t_name RECORD;
    r table_count%ROWTYPE;


BEGIN
    FOR t_name IN
        SELECT
            c.relname
        FROM
            pg_catalog.pg_class c LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE
            c.relkind = ''r''
            AND n.nspname = ''public''
        ORDER BY 1
        LOOP
            FOR the_count IN EXECUTE ''SELECT COUNT(*) AS "count" FROM '' || t_name.relname
            LOOP
            END LOOP;


            r.table_name := t_name.relname;
            r.num_rows := the_count.count;
            RETURN NEXT r;
        END LOOP;
        RETURN;
END;
' LANGUAGE plpgsql;


SELECT
    schemaname
     ,relname
     ,n_live_tup AS EstimatedCount
FROM pg_stat_user_tables
ORDER BY n_live_tup DESC;