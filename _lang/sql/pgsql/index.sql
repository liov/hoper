# 复合主键
首先删除主键约束.您可以通过键入来获取约束的名称

\d my_table
并在索引下查找类似的内容:

"my_table_pkey" PRIMARY KEY, btree (datetime, uid)
放弃它:

alter table my_table drop constraint my_table_pkey;
然后通过执行以下操作创建新的复合主键:

alter table my_table add constraint my_table_pkey primary key (datetime, uid, action);

# 联合索引
CREATE INDEX index_name
    ON table_name (column1_name, column2_name);
# 唯一索引
CREATE UNIQUE INDEX index_name
    on table_name (column_name);

# 局部索引
CREATE INDEX index_name
    on table_name (conditional_expression);

# 删除索引
DROP INDEX index_name;