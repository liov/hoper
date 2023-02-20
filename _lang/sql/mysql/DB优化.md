1   优化的哲学
注：优化有风险，涉足需谨慎

a   优化可能带来的问题？
优化不总是对一个单纯的环境进行，还很可能是一个复杂的已投产的系统；

优化手段本来就有很大的风险，只不过你没能力意识到和预见到；

任何的技术可以解决一个问题，但必然存在带来一个问题的风险；

对于优化来说解决问题而带来的问题，控制在可接受的范围内才是有成果；

保持现状或出现更差的情况都是失败！

b   优化的需求？
稳定性和业务可持续性，通常比性能更重要；

优化不可避免涉及到变更，变更就有风险；

优化使性能变好，维持和变差是等概率事件；

切记优化，应该是各部门协同，共同参与的工作，任何单一部门都不能对数据库进行优化！

所以优化工作，是由业务需要驱使的！

c   优化由谁参与？

在进行数据库优化时，应由数据库管理员、业务部门代表、应用程序架构师、应用程序设计人员、应用程序开发人员、硬件及系统管理员、存储管理员等，业务相关人员共同参与。 

2   优化思路

a  优化什么？
在数据库优化上有两个主要方面：即安全与性能。

安全->数据可持续性；

性能->数据的高性能访问。

b   优化的范围有哪些？
存储、主机和操作系统方面：

主机架构稳定性；

I/O规划及配置；

Swap交换分区；

OS内核参数和网络问题。

应用程序方面：

应用程序稳定性；

SQL语句性能；

串行访问资源；

性能欠佳会话管理；

这个应用适不适合用MySQL。

数据库优化方面：

内存；

数据库结构（物理&逻辑）；

实例配置。

说明：不管是设计系统、定位问题还是优化，都可以按照这个顺序执行。

c   优化维度？
数据库优化维度有四个：

硬件、系统配置、数据库表结构、SQL及索引。



优化选择：

优化成本：硬件>系统配置>数据库表结构>SQL及索引。

优化效果：硬件<系统配置<数据库表结构<SQL及索引。

1   优化工具有啥？
 
a   数据库层面？
检查问题常用工具：

1）MySQL

2）msyqladmin：MySQL客户端，可进行管理操作

3）mysqlshow：功能强大的查看shell命令

4）show [SESSION | GLOBAL] variables：查看数据库参数信息

5）SHOW [SESSION | GLOBAL] STATUS：查看数据库的状态信息

6）information_schema：获取元数据的方法

7）SHOW ENGINE INNODB STATUS：Innodb引擎的所有状态

8）SHOW PROCESSLIST：查看当前所有连接session状态

9）explain：获取查询语句的执行计划

10）show index：查看表的索引信息

11）slow-log：记录慢查询语句

12）mysqldumpslow：分析slowlog文件的

不常用但好用的工具：

1）Zabbix：监控主机、系统、数据库（部署zabbix监控平台）

2）pt-query-digest：分析慢日志

3）MySQL slap：分析慢日志

4）sysbench：压力测试工具

5）MySQL profiling：统计数据库整体状态工具    

6）Performance Schema：MySQL性能状态统计的数据

7）workbench：管理、备份、监控、分析、优化工具（比较费资源）

关于Zabbix参考：

http://www.cnblogs.com/clsn/p/7885990.html

b   数据库层面问题解决思路？

一般应急调优的思路：针对突然的业务办理卡顿，无法进行正常的业务处理，需要立马解决的场景。

1）show processlist；

2）explain  select id ,name from stu where name='clsn'; # ALL  id name age  sex；

select id,name from stu  where id=2-1 函数 结果集>30；show index from table；

3）通过执行计划判断，索引问题（有没有、合不合理）或者语句本身问题；

4）show status  like '%lock%';    # 查询锁状态

kill SESSION_ID;   # 杀掉有问题的session。

常规调优思路：针对业务周期性的卡顿，例如在每天10-11点业务特别慢，但是还能够使用，过了这段时间就好了。

1）查看slowlog，分析slowlog，分析出查询慢的语句；

2）按照一定优先级，一个一个排查所有慢语句；

3）分析top SQL，进行explain调试，查看语句执行时间；

4）调整索引或语句本身。

c  系统层面？

Cpu方面：

vmstat、sar top、htop、nmon、mpstat；
内存：

free、ps-aux；
IO设备（磁盘、网络）：

iostat、ss、netstat、iptraf、iftop、lsof；
vmstat命令说明：

1）Procs：r显示有多少进程正在等待CPU时间。b显示处于不可中断的休眠的进程数量。在等待I/O。

2）Memory：swpd显示被交换到磁盘的数据块的数量。未被使用的数据块，用户缓冲数据块，用于操作系统的数据块的数量。

3）Swap：操作系统每秒从磁盘上交换到内存和从内存交换到磁盘的数据块的数量。s1和s0最好是0。

4）Io：每秒从设备中读入b1的写入到设备b0的数据块的数量。反映了磁盘I/O。

5）System：显示了每秒发生中断的数量（in）和上下文交换（cs）的数量。

6）Cpu：显示用于运行用户代码，系统代码，空闲，等待I/O的Cpu时间。

iostat命令说明：

实例命令：iostat -dk 1 5

　　　　   iostat -d -k -x 5 （查看设备使用率（%util）和响应时间（await））

1）tps：该设备每秒的传输次数。“一次传输”意思是“一次I/O请求”。多个逻辑请求可能会被合并为“一次I/O请求”。

2）iops ：硬件出厂的时候，厂家定义的一个每秒最大的IO次数

3）"一次传输"请求的大小是未知的。

4）kB_read/s：每秒从设备（drive expressed）读取的数据量；

5）KB_wrtn/s：每秒向设备（drive expressed）写入的数据量；

6）kB_read：读取的总数据量；

7）kB_wrtn：写入的总数量数据量；这些单位都为Kilobytes。

d  系统层面问题解决办法？
你认为到底负载高好，还是低好呢？在实际的生产中，一般认为Cpu只要不超过90%都没什么问题。

当然不排除下面这些特殊情况：

Cpu负载高，IO负载低：

1）内存不够；

2）磁盘性能差；

3）SQL问题--->去数据库层，进一步排查SQL 问题；

4）IO出问题了（磁盘到临界了、raid设计不好、raid降级、锁、在单位时间内tps过高）；

5）tps过高：大量的小数据IO、大量的全表扫描。

IO负载高，Cpu负载低：

1）大量小的IO写操作：

autocommit，产生大量小IO；IO/PS，磁盘的一个定值，硬件出厂的时候，厂家定义的一个每秒最大的IO次数。

2）大量大的IO 写操作：SQL问题的几率比较大

IO和cpu负载都很高：

硬件不够了或SQL存在问题。
4   基础优化
 
a  优化思路？
定位问题点吮吸：硬件-->系统-->应用-->数据库-->架构（高可用、读写分离、分库分表）。

处理方向：明确优化目标、性能和安全的折中、防患未然。

b  硬件优化？
主机方面：

根据数据库类型，主机CPU选择、内存容量选择、磁盘选择：

1）平衡内存和磁盘资源；

2）随机的I/O和顺序的I/O；

3）主机 RAID卡的BBU（Battery Backup Unit）关闭。

CPU的选择：

CPU的两个关键因素：核数、主频

根据不同的业务类型进行选择：

1）CPU密集型：计算比较多，OLTP - 主频很高的cpu、核数还要多

2）IO密集型：查询比较，OLAP - 核数要多，主频不一定高的

内存的选择：

OLAP类型数据库，需要更多内存，和数据获取量级有关。

OLTP类型数据一般内存是Cpu核心数量的2倍到4倍，没有最佳实践。

存储方面：

1）根据存储数据种类的不同，选择不同的存储设备；

2）配置合理的RAID级别（raid5、raid10、热备盘）；

3）对与操作系统来讲，不需要太特殊的选择，最好做好冗余（raid1）（ssd、sas、sata）。

4）raid卡：

       主机raid卡选择：

           实现操作系统磁盘的冗余（raid1）；

           平衡内存和磁盘资源；

           随机的I/O和顺序的I/O；

           主机raid卡的BBU（Battery Backup Unit）要关闭。

网络设备方面：

使用流量支持更高的网络设备（交换机、路由器、网线、网卡、HBA卡）
注意：以上这些规划应该在初始设计系统时就应该考虑好。

c  服务器硬件优化？
1）物理状态灯

2）自带管理设备：远程控制卡（FENCE设备：ipmi ilo idarc）、开关机、硬件监控。

3）第三方的监控软件、设备（snmp、agent）对物理设施进行监控。

4）存储设备：自带的监控平台。EMC2（hp收购了）、 日立（hds）、IBM低端OEM hds、高端存储是自己技术，华为存储。

d  系统优化？
Cpu：

基本不需要调整，在硬件选择方面下功夫即可。
内存：

基本不需要调整，在硬件选择方面下功夫即可。
SWAP：

MySQL尽量避免使用swap。

阿里云的服务器中默认swap为0。

IO ：

raid、no lvm、ext4或xfs、ssd、IO调度策略。

Swap调整(不使用swap分区)

/proc/sys/vm/swappiness的内容改成0（临时），/etc/sysctl. conf上添加vm.swappiness=0（永久）

这个参数决定了Linux是倾向于使用swap，还是倾向于释放文件系统cache。在内存紧张的情况下，数值越低越倾向于释放文件系统cache。

当然，这个参数只能减少使用swap的概率，并不能避免Linux使用swap。

修改MySQL的配置参数innodb_flush_ method，开启O_DIRECT模式：

这种情况下，InnoDB的buffer pool会直接绕过文件系统cache来访问磁盘，但是redo log依旧会使用文件系统cache。

值得注意的是，Redo log是覆写模式的，即使使用了文件系统的cache，也不会占用太多。

IO调度策略：

#echo deadline>/sys/block/sda/queue/scheduler   临时修改为deadline

永久修改

vi /boot/grub/grub.conf

更改到如下内容:

kernel /boot/vmlinuz-2.6.18-8.el5 ro root=LABEL=/ elevator=deadline rhgb quiet

e   系统参数调整？

Linux系统内核参数优化：

vim/etc/sysctl.conf

net.ipv4.ip_local_port_range = 1024 65535：# 用户端口范围

net.ipv4.tcp_max_syn_backlog = 4096 

net.ipv4.tcp_fin_timeout = 30 

fs.file-max=65535：# 系统最大文件句柄，控制的是能打开文件最大数量  

用户限制参数（MySQL可以不设置以下配置）：

vim/etc/security/limits.conf 

* soft nproc 65535

* hard nproc 65535

* soft nofile 65535

* hard nofile 65535

f   应用优化？

业务应用和数据库应用独立；

防火墙：iptables、selinux等其他无用服务（关闭）：

   chkconfig --level 23456 acpid off

    chkconfig --level 23456 anacron off

    chkconfig --level 23456 autofs off

    chkconfig --level 23456 avahi-daemon off

    chkconfig --level 23456 bluetooth off

    chkconfig --level 23456 cups off

    chkconfig --level 23456 firstboot off

    chkconfig --level 23456 haldaemon off

    chkconfig --level 23456 hplip off

    chkconfig --level 23456 ip6tables off

    chkconfig --level 23456 iptables  off

    chkconfig --level 23456 isdn off

    chkconfig --level 23456 pcscd off

    chkconfig --level 23456 sendmail  off

    chkconfig --level 23456 yum-updatesd  off

安装图形界面的服务器不要启动图形界面runlevel 3。 

另外，思考将来我们的业务是否真的需要MySQL，还是使用其他种类的数据库。用数据库的最高境界就是不用数据库。

5  数据库优化

SQL优化方向：执行计划、索引、SQL改写。

架构优化方向：高可用架构、高性能架构、分库分表。

a  数据库参数优化？
调整

实例整体（高级优化，扩展）：

thread_concurrency：# 并发线程数量个数

sort_buffer_size：# 排序缓存

read_buffer_size：# 顺序读取缓存

read_rnd_buffer_size：# 随机读取缓存

key_buffer_size：# 索引缓存

thread_cache_size：# (1G—>8, 2G—>16, 3G—>32, >3G—>64)

连接层（基础优化）

设置合理的连接客户和连接方式：

max_connections           # 最大连接数，看交易笔数设置    

max_connect_errors        # 最大错误连接数，能大则大

connect_timeout           # 连接超时

max_user_connections      # 最大用户连接数

skip-name-resolve         # 跳过域名解析

wait_timeout              # 等待超时

back_log                  # 可以在堆栈中的连接数量

SQL层（基础优化）

query_cache_size： 查询缓存  >>>  OLAP类型数据库,需要重点加大此内存缓存，但是一般不会超过GB。

对于经常被修改的数据，缓存会立马失效。

我们可以实用内存数据库（redis、memecache），替代他的功能。

b  存储引擎层（innodb基础优化参数）？
default-storage-engine

innodb_buffer_pool_size       # 没有固定大小，50%测试值，看看情况再微调。但是尽量设置不要超过物理内存70%

innodb_file_per_table=(1,0)

innodb_flush_log_at_trx_commit=(0,1,2) # 1是最安全的，0是性能最高，2折中

binlog_sync

Innodb_flush_method=(O_DIRECT, fdatasync)

innodb_log_buffer_size           # 100M以下

innodb_log_file_size               # 100M 以下

innodb_log_files_in_group       # 5个成员以下,一般2-3个够用（iblogfile0-N）

innodb_max_dirty_pages_pct   # 达到百分之75的时候刷写 内存脏页到磁盘。

log_bin

max_binlog_cache_size                     # 可以不设置

max_binlog_size                               # 可以不设置

innodb_additional_mem_pool_size     #小于2G内存的机器，推荐值是20M。32G内存以上100M

谈谈项目中常用的MySQL优化方法，共19条，具体如下：

1、EXPLAIN

做MySQL优化，我们要善用EXPLAIN查看SQL执行计划。

下面来个简单的示例，标注（1、2、3、4、5）我们要重点关注的数据：



type列，连接类型。一个好的SQL语句至少要达到range级别。杜绝出现all级别。

key列，使用到的索引名。如果没有选择索引，值是NULL。可以采取强制索引方式。

key_len列，索引长度。

rows列，扫描行数。该值是个预估值。

extra列，详细说明。注意，常见的不太友好的值，如下：Using filesort，Using temporary。

2、SQL语句中IN包含的值不应过多

MySQL对于IN做了相应的优化，即将IN中的常量全部存储在一个数组里面，而且这个数组是排好序的。但是如果数值较多，产生的消耗也是比较大的。再例如：select id from t where num in(1,2,3) 对于连续的数值，能用between就不要用in了；再或者使用连接来替换。

3、SELECT语句务必指明字段名称

SELECT*增加很多不必要的消耗（CPU、IO、内存、网络带宽）；增加了使用覆盖索引的可能性；当表结构发生改变时，前断也需要更新。所以要求直接在select后面接上字段名。

4、当只需要一条数据的时候，使用limit 1

这是为了使EXPLAIN中type列达到const类型

5、如果排序字段没有用到索引，就尽量少排序

6、如果限制条件中其他字段没有索引，尽量少用or

or两边的字段中，如果有一个不是索引字段，而其他条件也不是索引字段，会造成该查询不走索引的情况。很多时候使用union all或者是union（必要的时候）的方式来代替“or”会得到更好的效果。

7、尽量用union all代替union

union和union all的差异主要是前者需要将结果集合并后再进行唯一性过滤操作，这就会涉及到排序，增加大量的CPU运算，加大资源消耗及延迟。当然，union all的前提条件是两个结果集没有重复数据。

8、不使用ORDER BY RAND()

select id from `dynamic` order by rand() limit 1000;

上面的SQL语句，可优化为：

select id from `dynamic` t1 join (select rand() * (select max(id) from `dynamic`) as nid) t2 on t1.id > t2.nidlimit 1000;



9、区分in和exists、not in和not exists



select * from 表A where id in (select id from 表B)


上面SQL语句相当于

select * from 表A where exists(select * from 表B where 表B.id=表A.id)

区分in和exists主要是造成了驱动顺序的改变（这是性能变化的关键），如果是exists，那么以外层表为驱动表，先被访问，如果是IN，那么先执行子查询。所以IN适合于外表大而内表小的情况；EXISTS适合于外表小而内表大的情况。

关于not in和not exists，推荐使用not exists，不仅仅是效率问题，not in可能存在逻辑问题。如何高效的写出一个替代not exists的SQL语句？

原SQL语句：

select colname … from A表 where a.id not in (select b.id from B表)
高效的SQL语句：

select colname … from A表 Left join B表 on where a.id = b.id where b.id is null

取出的结果集如下图表示，A表不在B表中的数据：





10、使用合理的分页方式以提高分页的效率

select id,name from product limit 866613, 20



使用上述SQL语句做分页的时候，可能有人会发现，随着表数据量的增加，直接使用limit分页查询会越来越慢。

优化的方法如下：可以取前一页的最大行数的id，然后根据这个最大的id来限制下一页的起点。比如此列中，上一页最大的id是866612。SQL可以采用如下的写法：

select id,name from product where id> 866612 limit 20

11、分段查询

一些用户选择页面中，可能一些用户选择的时间范围过大，造成查询缓慢。主要的原因是扫描行数过多。这个时候可以通过程序，分段进行查询，循环遍历，将结果合并处理进行展示。

如下图这个SQL语句，扫描的行数成百万级以上的时候就可以使用分段查询：



12、避免在where子句中对字段进行null值判断

对于null的判断会导致引擎放弃使用索引而进行全表扫描。

13、不建议使用%前缀模糊查询

例如LIKE“%name”或者LIKE“%name%”，这种查询会导致索引失效而进行全表扫描。但是可以使用LIKE “name%”。

那如何查询%name%？

如下图所示，虽然给secret字段添加了索引，但在explain结果并没有使用：



那么如何解决这个问题呢，答案：使用全文索引。

在我们查询中经常会用到select id,fnum,fdst from dynamic_201606 where user_name like '%zhangsan%'; 。这样的语句，普通索引是无法满足查询需求的。庆幸的是在MySQL中，有全文索引来帮助我们。

创建全文索引的SQL语法是：

ALTER TABLE `dynamic_201606` ADD FULLTEXT INDEX `idx_user_name` (`user_name`);

使用全文索引的SQL语句是：

select id,fnum,fdst from dynamic_201606 where match(user_name) against('zhangsan' in boolean mode);

注意：在需要创建全文索引之前，请联系DBA确定能否创建。同时需要注意的是查询语句的写法与普通索引的区别。

14、避免在where子句中对字段进行表达式操作

比如：

select user_id,user_project from user_base where age*2=36;

中对字段就行了算术运算，这会造成引擎放弃使用索引，建议改成：

select user_id,user_project from user_base where age=36/2;



15、避免隐式类型转换

where子句中出现column字段的类型和传入的参数类型不一致的时候发生的类型转换，建议先确定where中的参数类型。

16、对于联合索引来说，要遵守最左前缀法则

举列来说索引含有字段id、name、school，可以直接用id字段，也可以id、name这样的顺序，但是name;school都无法使用这个索引。所以在创建联合索引的时候一定要注意索引字段顺序，常用的查询字段放在最前面。

17、必要时可以使用force index来强制查询走某个索引

有的时候MySQL优化器采取它认为合适的索引来检索SQL语句，但是可能它所采用的索引并不是我们想要的。这时就可以采用forceindex来强制优化器使用我们制定的索引。

18、注意范围查询语句

对于联合索引来说，如果存在范围查询，比如between、>、<等条件时，会造成后面的索引字段失效。

19、关于JOIN优化







LEFT JOIN A表为驱动表，INNER JOIN MySQL会自动找出那个数据少的表作用驱动表，RIGHT JOIN B表为驱动表。

注意：

1）MySQL中没有full join，可以用以下方式来解决：

select * from A left join B on B.name = A.name where B.name is null union all select * from B;

2）尽量使用inner join，避免left join：

参与联合查询的表至少为2张表，一般都存在大小之分。如果连接方式是inner join，在没有其他过滤条件的情况下MySQL会自动选择小表作为驱动表，但是left join在驱动表的选择上遵循的是左边驱动右边的原则，即left join左边的表名为驱动表。

3）合理利用索引：

被驱动表的索引字段作为on的限制字段。

4）利用小表去驱动大表：





从原理图能够直观的看出如果能够减少驱动表的话，减少嵌套循环中的循环次数，以减少 IO总量及CPU运算的次数。

5）巧用STRAIGHT_JOIN：

inner join是由MySQL选择驱动表，但是有些特殊情况需要选择另个表作为驱动表，比如有group by、order by等「Using filesort」、「Using temporary」时。STRAIGHT_JOIN来强制连接顺序，在STRAIGHT_JOIN左边的表名就是驱动表，右边则是被驱动表。在使用STRAIGHT_JOIN有个前提条件是该查询是内连接，也就是inner join。其他链接不推荐使用STRAIGHT_JOIN，否则可能造成查询结果不准确。