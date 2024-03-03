  1.创建文件 ：fallocate -l 2G /data/swapfile

        2.创建交换分区： mkswap /data/swapfile

        3.开启swap：swapon /data/swapfile，开启后可通过free验证

关闭swap
        swapoff -a

swap使用统计
        使用sar -S 1  每秒统计一次swap使用情况

        或使用for file in /proc/*/status ; do awk '/VmSwap|Name|^Pid/{printf $2 " " $3}END{ print ""}' $file; done | sort -k 3 -n -r | head查看各进程使用swap情况
