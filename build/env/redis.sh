wget http://download.redis.io/releases/redis-5.0.5.tar.gz

tar -xzvf redis-5.0.5.tar.gz

sudo: unable to resolve host abc 虽然sudo 还是可以正常执行, 但是看到这样的通知还是会觉得烦，怎么去除这个警告呢？这个警告是因为系统找不到一个叫做 abc的hostname 通过 修改 /etc/hosts 设定, 可以解决在127.0.0.1 localhost 后面加上主机名称(hostname) 即可:127.0.0.1 localhost abc

vim /etc/hostname localhost

sudo apt install tcl

vim redis.conf

grep -n requirepass redis.conf
/requirepass
requirepass ******
/bind
这是外网可连
取消绑定局域网
#bind 127.0.0.1
取消保护模式
protected-mode no
daemonize no 为 yes 并保存

./src/redis-server ./redis.conf
