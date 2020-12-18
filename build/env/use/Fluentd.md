EFK由ElasticSearch、Fluentd和Kiabana三个开源工具组成。其中Elasticsearch是一款分布式搜索引擎，能够用于日志的检索，Fluentd是一个实时开源的数据收集器,而Kibana 是一款能够为Elasticsearch 提供分析和可视化的 Web 平台。这三款开源工具的组合为日志数据提供了分布式的实时搜集与分析的监控系统。

而在此之前，业界是采用ELK(Elasticsearch + Logstash + Kibana)来管理日志。Logstash是一个具有实时渠道能力的数据收集引擎,但和fluentd相比，它在效能上表现略逊一筹，故而逐渐被fluentd取代，ELK也随之变成EFK。

Log Source：日志来源。在微服务中，我们的日志主要来源于日志文件和Docker容器，日志文件包括服务器log，例如Nginx access log（记录了哪些用户，哪些页面以及用户浏览器、ip和其他的访问信息）, error log(记录服务器错误日志)等。
Logstash：数据收集处理引擎，可用于传输docker各个容器中的日志给EK。支持动态的从各种数据源搜集数据，并对数据进行过滤、分析、丰富、统一格式等操作，然后存储以供后续使用。
Filebeat：和Logstash一样属于日志收集处理工具，基于原先 Logstash-fowarder 的源码改造出来的。与Logstash相比，filebeat更加轻量，占用资源更少
ElasticSearch:日志搜索引擎
Kibana:用于日志展示的可视化工具
Grafana:类似Kibana，可对后端的数据进行实时展示

EFK架构
Fluentd是一个开源的数据收集器，专为处理数据流设计，使用JSON作为数据格式。它采用了插件式的架构，具有高可扩展性高可用性，同时还实现了高可靠的信息转发。

因此，我们加入Fluentd来收集日志。加入后的EFK架构如图所示。

