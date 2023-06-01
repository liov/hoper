pg_dumpall -U postgres -p 5432 > bak.sql
psql -U postgres -f bak.sql

kubectl exec pod_name  -n namespace -- pg_dumpall -U postgres -p 5432 > bak.sql
kubectl exec pod_name  -n namespace -- pg_dump -U postgres -p 5432 -d test > bak.sql

# 待测试
kubectl exec pod_name  -n namespace -- psql -U postgres < bak.sql
--inserts #insert语句导出
docker exec postgres-old pg_dumpall -U postgres | docker exec -i postgres-new psql -U postgres

kubectl exec -it postgres-f9b466ff-lrv2z  -n tools --  psql -U postgres;

set time zone "Asia/Shanghai";
SET TIMEZONE='Asia/Shanghai';

vim postgresql.conf
log_timezone = 'Asia/Shanghai'
timezone = 'Asia/Shanghai'

# upgrade
pg_dumpall -U postgres -h postgres.tools -p 5432 | psql -U postgres

vim /home/postgres/data/pg_hba.conf
host    all     all     0.0.0.0/0        md5