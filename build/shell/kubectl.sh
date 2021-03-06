# 扩容
kubectl  scale deployment nginx-deployment --replicas=4

kubectl get pod pod名称 -n 命名空间名称 -o yaml | kubectl replace --force -f -.

kubectl create configmap my-config-2 --from-file=/etc/resolv.conf
# key的名称是文件名称，value的值是这个文件的内容

kubectl create configmap my-config-3 --from-file=test
# 目录中的文件名为key，文件内容是value

# log
# shellcheck disable=SC2006
kubectl config use-context stage -n namespace&& pod=`kubectl get pods |grep -oE podname[a-zA-Z0-9-]*` && kubectl logs -f $pod

pods=$(kubectl get pods --selector=job-name=pi --output=jsonpath={.items..metadata.name}) && kubectl logs $pods

kubectl describe deployment data-center

pods=$(kubectl get pods --selector=app=${PWD##*/} --output=jsonpath={.items..metadata.name}) && kubectl logs -f $pods

# deploy
../deploy/main -flow all -env dev -name ${PWD##*/} -ns openmng -path . -ver v1.1.0-$(date "+%Y%m%d%H%M%S")

env=stage && git pull && make config && make deploy env=$env tag=v$(date "+%y%m%d%H%M")

# 获取当前时间

#shell 实现获取当前时间，并进行格式转换的方法：

#1）原格式输出

#2018年 09月 30日 星期日 15:55:15 CST

time1=$(date)
echo $time1


#2）时间串输出

#20180930155515

#!bin/bash
time2=$(date "+%Y%m%d%H%M%S")
echo $time2


#3）2018-09-30 15:55:15

#!bin/bash
time3=$(date "+%Y-%m-%d %H:%M:%S")
echo $time3
#4）2018.09.30

#!bin/bash
time4=$(date "+%Y.%m.%d")
echo $time4
#注意
#
#1、date后面有一个空格，shell对空格要求严格
#
#2、变量赋值前后不要有空格

#1 Y显示4位年份，如：2018；y显示2位年份，如：18。
#2 m表示月份；M表示分钟。
#3 d表示天；D则表示当前日期，如：1/18/18(也就是2018.1.18)。
#4 H表示小时，而h显示月份。
#5 s显示当前秒钟，单位为毫秒；S显示当前秒钟，单位为秒。

#获取当前目录
#!/bin/bash
${PWD##*/}