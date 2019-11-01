看端口：

ps -aux | grep tomcat

发现并没有8080端口的Tomcat进程。

使用命令：netstat –apn

查看所有的进程和端口使用情况。发现下面的进程列表，其中最后一栏是PID/Program name

clip_image002

发现8080端口被PID为9658的Java进程占用。

进一步使用命令：ps -aux | grep java，或者直接：ps -aux | grep pid 查看

clip_image004

就可以明确知道8080端口是被哪个程序占用了！然后判断是否使用KILL命令干掉！


方法二：直接使用 netstat   -anp   |   grep  portno
即：netstat –apn | grep 8080



查看进程：

1、ps 命令用于查看当前正在运行的进程。
grep 是搜索
例如： ps -ef | grep java
表示查看所有进程里 CMD 是 java 的进程信息
2、ps -aux | grep java
-aux 显示所有状态
ps
3. kill 命令用于终止进程
例如： kill -9 [PID]
-9 表示强迫进程立即停止
通常用 ps 查看进程 PID ，用 kill 命令终止进程