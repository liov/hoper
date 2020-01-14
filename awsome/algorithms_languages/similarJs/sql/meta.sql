SELECT TABLE_NAME,DATA_LENGTH,INDEX_LENGTH,(DATA_LENGTH+INDEX_LENGTH) as length,TABLE_ROWS,concat(round((DATA_LENGTH+INDEX_LENGTH)/1024/1024,3), 'MB') as total_size
FROM information_schema.TABLES
WHERE TABLE_SCHEMA='erp'
ORDER BY TABLE_ROWS ASC;

SELECT * FROM dpay.t_pay LIMIT 1 OFFSET 0;
show databases;