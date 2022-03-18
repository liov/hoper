# pogres表移动到另一个库
pg_dump -t table_to_copy source_db | psql target_db