create role replica login replication encrypted password 'replica';

pg_basebackup -h 192.168.10.88 -p 5432 -U repl -R -F p -P -D /var/lib/pgsql/14/data