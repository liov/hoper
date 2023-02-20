show data_directory;-- 查询配置文件所在位置
show config_file; -- 查询数据储存目录
select * from pg_database;
select * from pg_namespace;
select * from pg_tables;
select * from user;

--  查询所有表名
select
    relname as table_name,(select description from pg_description where objoid = oid and objsubid = 0) as table_comment
from pg_class
where
        relkind = 'r'
  and relname not like 'pg_%'
  and relname not like 'sql_%'
order by
    table_name;
-- 查询一个表的所有字段信息
select
    a.attname as 字段名称,
    format_type(a.atttypid,a.atttypmod) as 类型,
    (case when atttypmod-4>0 then atttypmod-4 else 0 end) as 长度,
    (case
         when (select count(*) from pg_constraint where conrelid = a.attrelid and conkey[1]=attnum and contype='p')>0 then 'PRI'
         when (select count(*) from pg_constraint where conrelid = a.attrelid and conkey[1]=attnum and contype='u')>0 then 'UNI'
         when (select count(*) from pg_constraint where conrelid = a.attrelid and conkey[1]=attnum and contype='f')>0 then 'FRI'
         else '' end) as 索引,
    (case when a.attnotnull=true then 'NO' else 'YES' end) as 允许为空,
    col_description(a.attrelid,a.attnum) as 说明
from pg_attribute a where attstattarget=-1 and attrelid = (select oid from pg_class where relname ='ttask');

SELECT relname,attname,typname,attnum FROM pg_class c,pg_attribute a,pg_type t
WHERE c.relname = 'ttask' AND c.oid = attrelid AND atttypid = t.oid AND attnum > 0;
-- 查询一个表的索引
SELECT t.relname AS table_name, c.relname,i.indnatts,i.indkey AS index_name FROM (
                                                                                     SELECT relname,indexrelid FROM pg_index i, pg_class c WHERE c.relname = 'xxl_job_qrtz_trigger_info' AND indrelid = c.oid) t,
                                                                                 pg_index i,pg_class c WHERE t.indexrelid = i.indexrelid AND i.indexrelid = c.oid;
-- 查询某Schema下的每张表的记录数：
select relname as TABLE_NAME, reltuples as rowCounts from pg_class where relkind = 'r' and relnamespace = (select oid from pg_namespace where nspname='public') order by rowCounts desc;

-- 查询某个表的表名和表注释
select relname as tabname,cast(obj_description(relfilenode,'pg_class') as varchar) as comment from pg_class c where   relname ='表名'

-- 查询表空间大小
select pg_size_pretty(pg_relation_size('表名'));
-- 统计各数据库占用的磁盘大小
SELECT d.datname AS Name,  pg_catalog.pg_get_userbyid(d.datdba) AS Owner,
       CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
                THEN pg_catalog.pg_size_pretty(pg_catalog.pg_database_size(d.datname))
            ELSE 'No Access'
           END AS SIZE
FROM pg_catalog.pg_database d
ORDER BY
    CASE WHEN pg_catalog.has_database_privilege(d.datname, 'CONNECT')
             THEN pg_catalog.pg_database_size(d.datname)
         ELSE NULL
        END DESC -- nulls first
LIMIT 20



系统表 ：
1）pg_authid表：包含有关数据库认证标识符(角色)的信息。一个角色体现"用户"和"组"的概念。一个用户实际上只是一个设置了 rolcanlogin 标志的角色。任何角色(不管设置了 rolcanlogin)标志)都可以有其它角色做为成员；因为用户标识是集群范围的，pg_authid 在一个集群里所有的数据库之间是共享的：每个集群只有一个 pg_authid 拷贝，而不是每个数据库一个。
字段描述：
rolname：角色名称
rolsuper：角色拥有超级用户权限
rolinherit：角色自动继承其所属角色的权限
rolcreaterole：角色可以创建更多角色
rolcreatedb：角色可以创建数据库
rolcatupdate：角色可以直接更新系统表。如果没有设置这个字段为真，即使超级用户也不能这么做。
rolcanlogin：角色可以登录，也就是说，这个角色可以给予会话认证标识符。
rolconnlimit：对于可以登录的角色，限制其最大并发连接数量。-1 表示没有限制。
rolpassword：口令(可能是加密的)；如果没有则为 NULL
rolvaliduntil：口令失效时间(只用于口令认证)；如果没有失效期，则为 NULL

2）pg_auth_members表： 显示角色之间的成员关系。任何非闭环的关系集合都是允许的。因为用户标识是集群范围的，pg_auth_members 是在一个集群里的所有数据库之间共享的：每个集群里只有一个 pg_auth_members 拷贝，而不是每个数据库一个。
字段描述：
roleid：拥有有成员的角色的 ID【pg_authid.oid】
member：属于 roleid 角色的一个成员的角色的 ID【pg_authid.oid】
grantor：赋予此成员关系的角色的 ID【pg_authid.oid】

3）pg_database表：表存储关于可用数据库的信息。和大多数系统表不同，pg_database 是在一个集群里的所有数据库共享的：每个集群只有一份 pg_database 拷贝，而不是每个数据库一份。
字段描述：
datname：数据库名字
datdba：数据库所有人，通常为其创建者【pg_authid.oid】

4）pg_class表：表记载表和几乎所有有字段或者是那些类似表的东西。包括索引(不过还要参阅 pg_index)、序列、视图、复合类型和一些特殊关系类型
字段描述：
relname：表、索引、视图等的名字
relnamespace：包含这个关系的名字空间(模式)的 OID【pg_namespace.oid】
relpersistence：p = 永久表，u = 无日志表， t = 临时表
relkind：r = 普通表， i = 索引， S = 序列， t = TOAST表， v = 视图， m = 物化视图， c = 组合类型， f = 外部表， p = 分区表， I = 分区索引

5）pg_description表：可以给每个数据库对象存储一个可选的描述(注释)。你可以用 COMMENT 命令操作这些描述，并且可以用 psql 的 \d 命令查看。许多内置的系统对象的描述提供了 pg_description 的初始内容
字段描述：
objoid：这条描述所描述的对象的 OID 【任意 oid 属性】
classoid：这个对象出现的系统表的 OID【pg_class.oid】
objsubid：对于一个表字段的注释，它是字段号(objoid 和 classoid 指向表自身)。对于其它对象类型，它是零。
description：作为对该对象的描述的任意文本

6）pg_index表：存储索引的具体信息
字段描述：
indexrelid：索引在pg_class里记录的oid
indrelid：使用这个表在pg_class里的记录的oid

7）pg_tablespace表：存储有关可用的表空间的信息。表可以放置在特定的表空间里，以帮助管理磁盘布局。与大多数系统表不同，pg_tablespace 在一个集群中的所有数据库之间共享：每个集群只有一份 pg_tablespace 的拷贝，而不是每个数据库一个。
字段描述：
spcname：表空间名
spcowner：表空间的所有者

8）pg_namespace表：模式/存储名字空间。名字空间是 SQL 模式下层的结构：每个名字空间有独立的关系，类型等集合但并不会相互冲突
字段描述：
nspnamce：名空间名称
nspowner：名空间的所有者

9）pg_attribute表：存储关于表的字段的信息。数据库里每个表的每个字段都在 pg_attribute 里有一行。还有用于索引，以及所有在 pg_class 里有记录的对象。
字段描述：
attrelid:字段所属的表oid【pg_class.oid】
attname:表字段名
atttypid:字段数据类型oid【pg_type.oid】
attnum：字段数目。普通字段是从 1 开始计数的。系统字段(比如 oid)有(任意)正数。

10）pg_inherits：存储表的分区关系
字段描述：
inhrelid：子表的OID【pg_class.oid】
inhparent：父表的OID【pg_class.oid】
inhseqno：如果一个子表存在多个直系父表(多重继承)，这个数字表明此继承字段的排列顺序。计数从1开始

11）pg_tables:提供对数据库中每个表的信息的访问
字段描述：
schemaname：模式名称【pg_namespace.nspname】
tablename：表名【pg_class.relname】
tableowner：表拥有者【pg_authid.rolname】
tablespace：表空间名称

12）pg_type:存储有关数据类型的信息。基本类型和枚举类型(标量类型)是用CREATE TYPE创建的， 域是使用CREATE DOMAIN创建的。同时还为数据库中每个表自动创建一个复合类型， 以表示该表的行结构。还可以用CREATE TYPE AS创建复合类型.
字段描述：
typname:数据类型名
typtype:对于基础类型是b，对于复合类型是c(比如，一个表的行类型)。 对于域类型是d，E的枚举类型，对于伪类型是p。
typrelid:如果是复合类型(参阅typtype)那么这个字段指向pg_class中定义该表的行。对于自由存在的复合类型，pg_class记录并不表示一个表，但是总需要它来查找该类型连接的pg_attribute记录。对于非复合类型为零。[pg_class.oid]

    一 系统表总览
系统表名	用途
pg_aggregate	聚集函数
pg_am	索引访问方法
pg_amop	访问方法操作符
pg_amproc	访问方法支持过程
pg_attrdef	列默认值
pg_attribute	表列 ( “属性” )
pg_authid	认证标识符（角色）
pg_auth_members	认证标识符成员关系
pg_cast	转换（数据类型转换）
pg_class	表、索引、序列、视图 （“关系”）
pg_collation	排序规则（locale信息）
pg_constraint	检查约束、唯一约束、主键约束、外键约束
pg_conversion	编码转换信息
pg_database	本数据库集簇中的数据库
pg_db_role_setting	每角色和每数据库的设置
pg_default_acl	对象类型的默认权限
pg_depend	数据库对象间的依赖
pg_description	数据库对象上的描述或注释
pg_enum	枚举标签和值定义
pg_event_trigger	事件触发器
pg_extension	已安装扩展
pg_foreign_data_wrapper	外部数据包装器定义
pg_foreign_server	外部服务器定义
pg_foreign_table	外部表信息
pg_index	索引信息
pg_inherits	表继承层次
pg_init_privs	对象初始特权
pg_language	编写函数的语言
pg_largeobjec	t大对象的数据页
pg_largeobject	_metadata大对象的元数据
pg_namespace	模式
pg_opclass	访问方法操作符类
pg_operator	操作符
pg_opfamily	访问方法操作符族
pg_partitioned_table	有关表的分区键的信息
pg_pltemplate	过程语言的模板数据
pg_policy	行安全策略
pg_proc	函数和过程
pg_publication	逻辑复制的发布
pg_publication_rel	与发布映射的关系
pg_range	范围类型的信息
pg_replication_origin	已注册的复制源
pg_rewrite	查询重写规则
pg_seclabel	数据库对象上的安全标签
pg_sequence	有关序列的信息
pg_shdepend	共享对象上的依赖
pg_shdescription	共享对象上的注释
pg_shseclabel	共享数据库对象上的安全标签
pg_statistic	规划器统计
pg_statistic_ext	扩展的规划器统计
pg_subscription	逻辑复制订阅
pg_subscription_rel	订阅的关系状态
pg_tablespace	本数据库集簇内的表空间
pg_transform	转换（将数据类型转换为过程语言需要的形式）
pg_trigger	触发器
pg_ts_config	文本搜索配置
pg_ts_config_map	文本搜索配置的记号映射
pg_ts_dict	文本搜索字典
pg_ts_parser	文本搜索分析器
pg_ts_template	文本搜索模板
pg_type	数据类型
pg_user_mapping	将用户映射到外部服务器
1 系统表
根据PostgreSQL12.2 的系统表进行整理的
1.1 pg_class
pg_ class 是数据字典最重要的一个表
pg_class记录表和几乎所有具有列或者像表的东西。这包括索引（但还要参见pg_index）、序列（但还要参见pg_sequence）、视图、物化视图、组合类型和TOAST表，参见relkind

每一个DDL/DML操作都必须跟这个表发生联系，在进行整库操作时经常使用到pg_class里面的东西，把它们整理出来，对数据库的了解有很大帮助。

名称	类型	引用	描述
oid	oid		行标识符
relname	name		表、索引、视图等的名字
relnamespace	oid	pg_namespace.oid	包含该关系的名字空间的OID
reltype	oid	pg_type.oid	可能存在的表行类型所对应数据类型的OID（对索引为0，索引没有pg_type项）
reloftype	oid	pg_type.oid	对于有类型的表，为底层组合类型的OID，对于其他所有关系为0
relowner	oid	pg_authid.oid	关系的拥有者
relam	oid	pg_am.oid	如果这是一个表或者索引，表示索引使用的访问方法（堆、B树、哈希等）
relfilenode	oid		该关系的磁盘文件的名字，0表示这是一个“映射”关系，其磁盘文件名取决于低层状态
reltablespace	oid	pg_tablespace.oid	该关系所存储的表空间。如果为0，使用数据库的默认表空间。（如果关系无磁盘文件时无意义）
relpages	int4		该表磁盘表示的尺寸，以页面计（页面尺寸为BLCKSZ）。这只是一个由规划器使用的估计值。它被VACUUM、ANALYZE以及一些DDL命令（如CREATE INDEX）所更新。
reltuples	float4		表中的存活行数。这只是一个由规划器使用的估计值。它被VACUUM、ANALYZE以及一些DDL命令（如CREATE INDEX）所更新。
relallvisible	int4		在表的可见性映射表中被标记为全可见的页数。这只是一个由规划器使用的估计值。它被VACUUM、ANALYZE以及一些DDL命令（如CREATE INDEX）所更新。
reltoastrelid	oid	pg_class.oid	与该表相关联的TOAST表的OID，如果没有则为0。TOAST表将大属性“线外”存储在一个二级表中。
relhasindex	bool		如果这是一个表并且其上建有（或最近建有）索引则为真
relisshared	bool		如果该表在集簇中的所有数据库间共享则为真。只有某些系统目录（如pg_database）是共享的。
relpersistence	char		p = 永久表，u = 无日志表， t = 临时表
relkind	char		r = 普通表， i = 索引， S = 序列， t = TOAST表， v = 视图， m = 物化视图， c = 组合类型， f = 外部表， p = 分区表， I = 分区索引
relnatts	int2		关系中用户列的数目（系统列不计算在内）。在pg_attribute中必须有这么多对应的项。另请参阅pg_attribute.attnum。
relchecks	int2		表上CHECK约束的数目，参见pg_constraint目录
relhasrules	bool		如果表有（或曾有）规则则为真，参见 pg_rewrite目录
relhastriggers	bool		如果表有（或曾有）触发器则为真，参见 pg_trigger目录
relhassubclass	bool		如果表或者索引有（或曾有）任何继承子女则为真
relrowsecurity	bool		如果表上启用了行级安全性则为真，参见 pg_policy目录
relforcerowsecurity	bool		如果行级安全性（启用时）也适用于表拥有者则为真，参见 pg_policy目录
relispopulated	bool		如果表已被填充则为真（对于所有关系该列都为真，但对于某些物化视图却不是）
relreplident	char		用来为行形成“replica identity”的列： d = 默认 (主键，如果存在), n = 无, f = 所有列 i = 索引的indisreplident被设置或者为默认
relispartition	bool		如果表或索引是一个分区，则为真
relrewrite	oid	pg_class.oid	对于在要求表重写的DDL操作期间被写入的新关系，这个域包含原始关系的OID，否则为0。那种状态仅在内部可见，对于一个用户可见的关系这个域应该从不包含不是0的值。
relfrozenxid	xid		在此之前的所有事务ID在表中已经被替换为一个永久的（“冻结的”) 事务ID。这用于跟踪表是否需要被清理，以便阻止事务ID回卷或者允许pg_xact被收缩。如果该关系不是一个表则为0（InvalidTransactionId）。
relminmxid	xid		在此之前的多事务ID在表中已经被替换为一个事务ID。这被用于跟踪表是否需要被清理，以阻止 多事务ID回卷或者允许pg_multixact被收缩。如果关系不是一个表则 为0（InvalidMultiXactId）。
relacl	aclitem[]		访问权限，更多信息参见第 5.7 节
reloptions	text[]		访问方法相关的选项，以“keyword=value”字符串形式
relpartbound	pg_node_tree		如果表示一个分区（见relispartition），分区边界的内部表达
1.2 pg_attribute
pg_attribute系统表存储所有表(包括系统表，如pg_class)的字段信息。数据库中的每个表的每个字段在pg_attribute表中都有一行记录。

名称	类型	引用	描述
attrelid	oid	pg_class.oid	列所属的表
attname	name		列名
atttypid	oid	pg_type.oid	列的数据类型
attstattarget	int4		attstattarget控制由ANALYZE对此列收集的统计信息的细节层次。0值表示不会收集任何统计信息。一个负值则说明直接使用系统默认的目标。正值的确切含义取决于数据类型。对于标量数据类型，attstattarget既是要收集的“最常见值”的目标号，也是要创建的柱状图容器的目标号。
attlen	int2		本列类型的pg_type.typlen一个拷贝
attnum	int2		列的编号。一般列从1开始向上编号。系统列（如ctid）则拥有（任意）负值编号。
attndims	int4		如果该列是一个数组类型，这里就是其维度数；否则为0。（在目前一个数组的维度数并不被强制，因此任何非零值都能有效地表明“这是一个数组”。）
attcacheoff	int4		在存储中总是为-1，但是当被载入到一个内存中的行描述符后，这里可能会被更新为属性在行内的偏移
atttypmod	int4		atttypmod记录了在表创建时提供的类型相关数据（例如一个varchar列的最大长度）。它会被传递给类型相关的输入函数和长度强制函数。对于那些不需要atttypmod的类型，这个值通常总是为-1。
attbyval	bool		该列类型的pg_type.typbyval的一个拷贝
attstorage	char		通常是该列类型的pg_type.typstorage的一个拷贝。对于可TOAST的数据类型，这可以在列创建后被修改以控制存储策略。
attalign	char		该列类型的pg_type.typalign的一个拷贝
attnotnull	bool		这表示一个非空约束。
atthasdef	bool		该列有一个默认表达式或生成的表达式，在此情况下在pg_attrdef目录中会有一个对应项来真正定义该表达式。（检查attgenerated以确定是默认还是生成的表达式。）
atthasmissing	bool		该列在行中完全缺失时会用到这个列的值，如果在行创建之后增加一个有非易失DEFAULT值的列，就会发生这种情况。实际使用的值被存放在attmissingval列中。
attidentity	char		如果是一个零字节（''），则不是一个标识列。否则，a = 总是生成，d = 默认生成。
attgenerated	char		如果为零字节('')，则不是生成的列。否则，s = stored。（将来可能会添加其他值。）
attisdropped	bool		该列被删除且不再有效。一个删除的列仍然物理存在于表中，但是会被分析器忽略并因此无法通过SQL访问。
attislocal	bool		该列是由关系本地定义的。注意一个列可以同时是本地定义和继承的。
attinhcount	int4		该列的直接祖先的编号。一个具有非零编号祖先的列不能被删除或者重命名。
attcollation	oid	pg_collation.oid	该列被定义的排序规则，如果该列不是一个可排序数据类型则为0。
attacl	aclitem[]		列级访问权限
attoptions	text[]		属性级选项，以“keyword=value”形式的字符串
attfdwoptions	text[]		属性级的外部数据包装器选项，以“keyword=value”形式的字符串
attmissingval	anyarray		这个列中是一个含有一个元素的数组，其中的值被用于该列在行中完全缺失时，如果在行创建之后增加一个有非易失DEFAULT值的列，就会发生这种情况。只有当atthasmissing为真时才使用这个值。如果没有值则该列为空。
1.3 pg_index
pg_index系统表存储关于索引的一部分信息。其它的信息大多数存储在pg_class

名称	类型	引用	描述
indexrelid	oid	pg_class.oid	此索引的pg_class项的OID
indrelid	oid	pg_class.oid	此索引的基表的pg_class项的OID
indnatts	int2		索引中的总列数（与pg_class.relnatts重复），这个数目包括键和被包括的属性
indnkeyatts	int2		索引中键列的编号，不计入任何的内含列，它们只是被存储但不参与索引的语义
indisunique	bool		表示是否为唯一索引
indisprimary	bool		表示索引是否表示表的主键（如果此列为真，indisunique也总是为真）
indisexclusion	bool		表示索引是否支持一个排他约束
indimmediate	bool		表示唯一性检查是否在插入时立即被执行（如果indisunique为假，此列无关）
indisclustered	bool		如果为真，表示表最后以此索引进行了聚簇
indisvalid	bool		如果为真，此索引当前可以用于查询。为假表示此索引可能不完整：它肯定还在被INSERT/UPDATE操作所修改，但它不能安全地被用于查询。如果索引是唯一索引，唯一性属性也不能被保证。
indcheckxmin	bool		如果为真，直到此pg_index行的xmin低于查询的TransactionXmin视界之前，查询都不能使用此索引，因为表可能包含具有它们可见的不相容行的损坏HOT链
indisready	bool		如果为真，表示此索引当前可以用于插入。为假表示索引必须被INSERT/UPDATE操作忽略。
indislive	bool		如果为假，索引正处于被删除过程中，并且必须被所有处理忽略（包括HOT安全的决策）
indisreplident	bool		如果为真，这个索引被选择为使用ALTER TABLE ... REPLICA IDENTITY USING INDEX ...的“replica identity”
indkey	int2vector	pg_attribute.attnum	这是一个indnatts值的数组，它表示了此索引索引的表列。例如一个1 3值可能表示表的第一和第三列组成了索引项。键列出现在非键（内含）列前面。数组中的一个0表示对应的索引属性是一个在表列上的表达式，而不是一个简单的列引用。
indcollation	oidvector	pg_collation.oid	对于索引键（indnkeyatts值）中的每一列，这包含要用于该索引的排序规则的OID，如果该列不是一种可排序数据类型则为零。
indclass	oidvector	pg_opclass.oid	对于索引键中的每一列（indnkeyatts值），这里包含了要使用的操作符类的OID。详见pg_opclass。
indoption	int2vector		这是一个indnkeyatts值的数组，用于存储每列的标志位。位的意义由索引的访问方法定义。
indexprs	pg_node_tree		非简单列引用索引属性的表达式树（以nodeToString()形式）。对于indkey中每一个为0的项，这个列表中都有一个元素。如果所有的索引属性都是简单引用，此列为空。
indpred	pg_node_tree		部分索引谓词的表达式树（以nodeToString()形式）。如果不是部分索引，此列为空。
1.4 pg_attrdef
pg_attrdef系统表主要存储字段缺省值，字段中的主要信息存放在pg_attribute系统表中。注意：只有明确声明了缺省值的字段在该表中才会
有记录。

名字	类型	引用	描述
oid	oid		行标识符
adrelid	oid	pg_class.oid	这个字段所属的表
adnum	int2	pg_attribute.attnum	字段编号，其规则等同于pg_attribute.attnum
adbin	text		字段缺省值的内部表现形式。
1.5 pg_constraint
pg_constraint系统表存储PostgreSQL中表对象的检查约束、主键、唯一约束和外键约束。

pg_constraint存储表上的检查、主键、唯一、外键和排他约束（列约束也不会被特殊对待。每一个列约束都等同于某种表约束。）。
非空约束不在这里，而是在pg_attribute目录中表示。
用户定义的约束触发器（使用CREATE CONSTRAINT TRIGGER创建）也会在这个表中产生一项。
域上的检查约束也存储在这里。
名称	类型	引用	描述
oid	oid		行标识符（隐藏属性，必须被显式选择才会显示）
conname	name		约束名字（不需要唯一！）
connamespace	oid	pg_namespace.oid	包含此约束的名字空间的OID
contype	char		c = 检查约束， f = 外键约束， p = 主键约束， u = 唯一约束， t = 约束触发器， x = 排他约束
condeferrable	bool		该约束是否能被延迟？
condeferred	bool		该约束是否默认被延迟？
convalidated	bool		此约束是否被验证过？当前对于外键和检查约束只能是假
conrelid	oid	pg_class.oid	该约束所在的表，如果不是表约束则为0
contypid	oid	pg_type.oid	该约束所在的域，如果不是域约束则为0
conindid	oid	pg_class.oid	如果该约束是唯一、主键、外键或排他约束，此列表示支持此约束的索引，否则为0
confrelid	oid	pg_class.oid	如果此约束是一个外键约束，此列为被引用的表，否则为0
confupdtype	char		外键更新动作代码： a = 无动作， r = 限制， c = 级联， n = 置空， d = 置为默认值
confdeltype	char		外键删除动作代码： a = 无动作， r = 限制， c = 级联， n = 置空， d = 置为默认值
confmatchtype	char		外键匹配类型： f = 完全， p = 部分， s = 简单
conislocal	bool		此约束是定义在关系本地。注意一个约束可以同时是本地定义和继承。
coninhcount	int4		此约束的直接继承祖先数目。一个此列非零的约束不能被删除或重命名。
connoinherit	bool		为真表示此约束被定义在关系本地。它是一个不可继承约束。
conkey	int2[]	pg_attribute.attnum	如果是一个表约束（包括外键但不包括约束触发器），此列是被约束列的列表
confkey	int2[]	pg_attribute.attnum	如果是一个外键，此列是被引用列的列表
conpfeqop	oid[]	pg_operator.oid	如果是一个外键，此列是用于PK = FK比较的等值操作符的列表
conppeqop	oid[]	pg_operator.oid	如果是一个外键，此列是用于PK = PK比较的等值操作符的列表
conffeqop	oid[]	pg_operator.oid	如果是一个外键，此列是用于FK = FK比较的等值操作符的列表
conexclop	oid[]	pg_operator.oid	如果是一个排他约束，此列是没列排他操作符的列表
conbin	pg_node_tree		如果是一个检查约束，此列是表达式的一个内部表示
consrc	text		如果是一个检查约束，此列是表达式的一个人类可读的表示
1.6 pg_tablespace
该系统表存储表空间的信息。注意：表可以放在特定的表空间里，以帮助管理磁盘布局和解决IO瓶颈。

名字	类型	引用	描述
oid	oid		行标识符
spcname	name		表空间名称。
spcowner	oid	pg_authid.oid	表空间的所有者，通常是创建它的角色。
spclocation	text		表空间的位置(目录路径)。
spcacl	aclitem[]		访问权限。
1.7 pg_namespace:
该系统表存储名字空间(模式)。

名字	类型	引用	描述
oid	oid		行标识符
nspname	name		名字空间(模式)的名称。
nspowner	oid	pg_authid.oid	名字空间(模式)的所有者
nspacl	aclitem[]		访问权限。
1.8 pg_database
pg_database系统表存储数据库的信息。和大多数系统表不同的是，在一个集群里该表是所有数据库共享的，即每个集群只有一份pg_database拷贝，而不是每个数据库一份

名称	类型	引用	描述
oid	oid		行标识符（隐藏属性，必须被显式选择才会显示）
datname	name		数据库名字
datdba	oid	pg_authid.oid	数据库的拥有者，通常是创建它的用户
encoding	int4		此数据库的字符编码的编号（pg_encoding_to_char()可将此编号转换成编码的名字）
datcollate	name		此数据库的LC_COLLATE
datctype	name		此数据库的LC_CTYPE
datistemplate	bool		如果为真，则此数据库可用在CREATE DATABASE的TEMPLATE子句中，该语句将从此数据库中克隆出一个新的数据库
datallowconn	bool		如果为假则没有人能连接到这个数据库。这可以用来保护template0数据库不被修改。
datconnlimit	int4		设置能够连接到这个数据库的最大并发连接数。-1表示没有限制。
datlastsysoid	oid		数据库中最后一个系统OID，对pg_dump特别有用
datfrozenxid	xid		在此之前的所有事务ID在数据库中已经被替换为一个永久的（“冻结的”) 事务ID。这用于跟踪数据库是否需要被清理，以便组织事务ID回环或者允许pg_clog被收缩。它是此数据库中所有表的pg_class.relfrozenxid值的最小值。
datminmxid	xid		在此之前的所有多事务ID在数据库中已经被替换为一个事务ID。这用于跟踪数据库是否需要被清理，以便组织事务ID回环或者允许pg_clog被收缩。它是此数据库中所有表的pg_class.relminmxid值的最小值。
dattablespace	oid	pg_tablespace.oid	此数据库的默认表空间。在此数据库中，所有pg_class.reltablespace为0的表都将被存储在这个表空间中，尤其是非共享系统目录都会在其中。
datacl	aclitem[]		访问权限，更多信息参见 GRANT和 REVOKE
