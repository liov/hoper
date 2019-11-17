https://www.postgresql.org/download/
```bash
#[binary]
https://www.postgresql.org/download/linux/ubuntu/
ubuntu 18.04
vim /etc/apt/sources.list.d/pgdg.list
deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update

#[source]
wget https://ftp.postgresql.org/pub/source/v11.3/postgresql-11.3.tar.gz
tar -xzvf
apt install libreadline-dev
apt-get install zlib1g-dev
./configure --prefix=/usr/local/postgresql
make&&make install
useradd  -m -s /bin/bash postgres

sudo passwd postgres

su postgres

mkdir /home/postgres/data

cd /home/jyb/

chown -R postgres:postgres postgre

cd /usr/local/postgresql/bin
./initdb -E UTF-8 -D /home/postgres/data --locale=en_US.UTF-8 -U postgres -W

./pg_ctl -D /home/postgres/data start

vim /home/postgres/data/pg_hba.conf

host    all     all     0.0.0.0/0                       md5

vim /home/postgres/data/postgresql.conf

listen_addresses = '*'

./pg_ctl -D /home/postgre/data/ reload

./psql

ALTER USER postgres WITH PASSWORD '123456';



sudo  passwd -d postgres

sudo -u postgres passwd
```
