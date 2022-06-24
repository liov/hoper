#!/bin/bash
set -e
hdfs_dir=$HADOOP_HOME/hdfs/
if [ $HADOOP_NODE_TYPE = "datanode" ]; then
  echo -e "\033[32m start datanode \033[0m"
  $HADOOP_HOME/bin/hdfs datanode -regular
fi
if [ $HADOOP_NODE_TYPE = "namenode" ]; then
  if [ -z $(ls -A ${hdfs_dir}) ]; then
    echo -e "\033[32m start hdfs namenode format \033[0m"
    $HADOOP_HOME/bin/hdfs namenode -format
  fi
  echo -e "\033[32m start hdfs namenode \033[0m"
  $HADOOP_HOME/bin/hdfs namenode
fi