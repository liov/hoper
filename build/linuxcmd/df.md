df命令的基本语法是：

df [options] [devices]

1.检查文件系统磁盘空间

“df”命令显示文件系统设备名称，磁盘块，使用的总磁盘空间，可用磁盘空间，文件系统上使用率和安装点的百分比等信息。

 -a, --all            include pseudo, duplicate, inaccessible file systems
  -B, --block-size=SIZE  scale sizes by SIZE before printing them; e.g.,
                          '-BM' prints sizes in units of 1,048,576 bytes;
                          see SIZE format below
  -h, --human-readable  print sizes in powers of 1024 (e.g., 1023M)
  -H, --si              print sizes in powers of 1000 (e.g., 1.1G)
  -i, --inodes  显示inode 信息而非块使用量
  -k   即--block-size=1K
  -l, --local  只显示本机的文件系统
      --no-sync  取得使用量数据前不进行同步动作(默认)
      --output[=FIELD_LIST]  use the output format defined by FIELD_LIST,
                              or print all fields if FIELD_LIST is omitted.
  -P, --portability    use the POSIX output format
      --sync            invoke sync before getting usage info
      --total          elide all entries insignificant to available space,
                          and produce a grand total
  -t, --type=TYPE      limit listing to file systems of type TYPE
  -T, --print-type      print file system type
  -x, --exclude-type=TYPE  limit listing to file systems not of type TYPE
  -v                    (ignored)
      --help  显示此帮助信息并退出
      --version  显示版本信息并退出

所显示的数值是来自 --block-size、DF_BLOCK_SIZE、BLOCK_SIZE 
及 BLOCKSIZE 环境变量中第一个可用的 SIZE 单位。
否则，默认单位是 1024 字节(或是 512，若设定 POSIXLY_CORRECT 的话)。


用法：du [选项]… [文件]…
或：du [选项]… --files0-from=F
计算每个文件的磁盘用量，目录则取总用量。

长选项必须使用的参数对于短选项时也是必需使用的。
-a, --all 输出所有文件的磁盘用量，不仅仅是目录
–apparent-size 显示表面用量，而并非是磁盘用量；虽然表面用量通常会
小一些，但有时它会因为稀疏文件间的"洞"、内部碎
片、非直接引用的块等原因而变大。
-B, --block-size=大小 使用指定字节数的块
-b, --bytes 等于–apparent-size --block-size=1
-c, --total 显示总计信息
-D, --dereference-args 解除命令行中列出的符号连接
–files0-from=F 计算文件F 中以NUL 结尾的文件名对应占用的磁盘空间
如果F 的值是"-"，则从标准输入读入文件名
-H 等于–dereference-args (-D)
-h, --human-readable 以可读性较好的方式显示尺寸(例如：1K 234M 2G)
–si 类似-h，但在计算时使用1000 为基底而非1024
-k 等于–block-size=1K
-l, --count-links 如果是硬连接，就多次计算其尺寸
-m 等于–block-size=1M
-L, --dereference 找出任何符号链接指示的真正目的地
-P, --no-dereference 不跟随任何符号链接(默认)
-0, --null 将每个空行视作0 字节而非换行符
-S, --separate-dirs 不包括子目录的占用量
-s, --summarize 只分别计算命令列中每个参数所占的总用量
-x, --one-file-system 跳过处于不同文件系统之上的目录
-X, --exclude-from=文件 排除与指定文件中描述的模式相符的文件
–exclude=PATTERN 排除与PATTERN 中描述的模式相符的文件
–max-depth=N 显示目录总计(与–all 一起使用计算文件)
当N 为指定数值时计算深度为N；
–max-depth=0 等于–summarize
–time 显示目录或该目录子目录下所有文件的最后修改时间
–time=WORD 显示WORD 时间，而非修改时间：
atime，access，use，ctime 或status
–time-style=样式 按照指定样式显示时间(样式解释规则同"date"命令)：
full-iso，long-iso，iso，+FORMAT
–help 显示此帮助信息并退出
–version 显示版本信息并退出

[大小]可以是以下的单位(单位前可加上整数)：
kB 1000，K 1024，MB 1000000，M 1048576，还有 G、T、P、E、Z、Y。