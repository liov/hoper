select * from information_schema.innodb_trx;
kill trx_requested_lock_id;

select t.trx_mysql_thread_id from information_schema.innodb_trx t;
kill trx_mysql_thread_id;