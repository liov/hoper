server_name与host匹配优先级如下：

1、完全匹配

2、通配符在前的，如*.test.com

3、在后的，如www.test.*

4、正则匹配，如~^\.www\.test\.com$

如果都不匹配

1、优先选择listen配置项后有default或default_server的

2、找到匹配listen端口的第一个server块