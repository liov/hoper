-- 一个跳板机连一个主机再连另一个主机，纯手敲命令，太屌了

mysql -h ${host} -P ${port} -u ${user} -p ${password}
-- 列出所有数据库
show database;
-- 选择
use ${database};
-- 列出所有表
show tables;
