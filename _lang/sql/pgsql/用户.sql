
create user gogs password 'gogs' createdb createrole;
grant connect ON DATABASE tools to gogs;
GRANT USAGE ON SCHEMA gogs TO gogs;
grant select,insert,update,delete ON  ALL TABLES IN SCHEMA gogs to gogs;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA gogs to gogs;

1.创建不需要密码登陆的用户zjy：
CREATE ROLE zjy LOGIN;

2.创建需要密码登陆的用户zjy1：
CREATE USER zjy1 WITH PASSWORD 'zjy1';

3.创建有时间限制的用户zjy2：
CREATE ROLE zjy2 WITH LOGIN PASSWORD 'zjy2' VALID UNTIL '2019-05-30';

4.创建有创建数据库和管理角色权限的用户admin：
CREATE ROLE admin WITH CREATEDB CREATEROLE;

5.创建具有超级权限的用户：admin
CREATE ROLE admin WITH SUPERUSER LOGIN PASSWORD 'admin';

6.创建复制账号：repl
CREATE USER repl REPLICATION LOGIN ENCRYPTED PASSWORD 'repl';

7.其他说明
创建复制用户
CREATE USER abc REPLICATION LOGIN ENCRYPTED PASSWORD '';
CREATE USER abc REPLICATION LOGIN ENCRYPTED PASSWORD 'abc';
ALTER USER work WITH ENCRYPTED password '';

创建scheme 角色
CREATE ROLE abc;
CREATE DATABASE abc WITH OWNER abc ENCODING UTF8 TEMPLATE template0;
\c abc

创建schema
CREATE SCHEMA abc;
ALTER SCHEMA abc OWNER to abc;
revoke create on schema public from public;

创建用户
create user abc with ENCRYPTED password '';
GRANT abc to abc;
ALTER ROLE abc WITH abc;

##创建读写账号
CREATE ROLE abc_rw;
CREATE ROLE abc_rr;

##赋予访问数据库权限，schema权限
grant connect ON DATABASE abc to abc_rw;
GRANT USAGE ON SCHEMA abc TO abc_rw;

##赋予读写权限
grant select,insert,update,delete ON  ALL TABLES IN SCHEMA abc to abc;

赋予序列权限
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA abc to abc;

赋予默认权限
ALTER DEFAULT PRIVILEGES IN SCHEMA abc GRANT select,insert,update,delete ON TABLES TO abc;

赋予序列权限
ALTER DEFAULT PRIVILEGES IN SCHEMA abc GRANT ALL PRIVILEGES ON SEQUENCES TO abc;


#用户对db要有连接权限
grant connect ON DATABASE abc to abc;

#用户要对schema usage 权限，不然要select * from schema_name.table ,不能用搜索路径
GRANT USAGE ON SCHEMA abc TO abc;
grant select ON ALL TABLES IN SCHEMA abc to abc;
ALTER DEFAULT PRIVILEGES IN SCHEMA abc GRANT select ON TABLES TO abc;

create user abc_w with ENCRYPTED password '';
create user abc_r with ENCRYPTED password '';

GRANT abc_rw to abc_w；

GRANT abc_rr to abc_r;