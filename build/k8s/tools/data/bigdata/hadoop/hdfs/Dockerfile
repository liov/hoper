FROM openjdk:11

ENV HADOOP_TGZ_URL=https://mirrors.tuna.tsinghua.edu.cn/apache/hadoop/common/hadoop-3.3.3/hadoop-3.3.3.tar.gz
ENV HADOOP_HOME=/opt/hadoop
ENV PATH=$HADOOP_HOME/bin:$PATH
RUN set -ex; \
    mkdir -p $HADOOP_HOME; \
    wget -nv -O $HADOOP_HOME/src.tgz $HADOOP_TGZ_URL; \
    tar -xf $HADOOP_HOME/src.tgz --strip-components=1 -C $HADOOP_HOME; \
    rm $HADOOP_HOME/src.tgz;

# clean trash file or dir
RUN rm -rf $HADOOP_HOME/share/doc/

COPY docker-entrypoint.sh /

RUN  chmod +x /docker-entrypoint.sh

EXPOSE 9870 9000
ENTRYPOINT ["/docker-entrypoint.sh"]