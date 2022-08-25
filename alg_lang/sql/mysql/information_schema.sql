use information_schema;

-- 不准确
select table_name,table_rows from tables where TABLE_SCHEMA = 'testdb' order by table_rows desc;