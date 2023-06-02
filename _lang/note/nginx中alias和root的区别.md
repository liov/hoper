root与alias区别

1.根路径与虚拟路径

如果访问站点http://location/c访问的就是/a/目录下的站点信息
```conf
location /c/ {
alias /a/
}
```
如果访问站点http://location/c访问的就是/a/c目录下的站点信息
```conf
location /c/ {
root /a/
}
```
2.根路径与虚拟路径

alias后面必须加 /

root结尾/可有可无

3.一般情况下，在location /中配置root，在location /other中配置alias是一个好习惯



root与alias的区别
 root与alias路径匹配主要区别在于nginx如何解释location后面的uri，这会使两者分别以不同的方式将请求映射到服务器文件上，alias是一个目录别名的定义，root则是最上层目录的定义。

 简而言之： root的处理结果是：location路径+root路径 alias的处理结果是：使用alias定义的路径
示例：

root配置：
vim root.conf
```conf
server {
    listen 80;
    server_name linux.root.com;
    location /download {
    root /code;
    }
}
```
使用root时，当我请求 http://linux.root.com/download/1.jpg 时，实际上是去找服务器上 /code/download/1.jpg 文件

alias配置：
vim /etc/nginx/conf.d/alias.conf
```conf
server {  
listen 80;  
    server_name linux.alias.com;
    location /download {    
    alias /code;  
    }
}
```
使用alias时，当我请求 http://linux.root.com/download/1.jpg 时，实际上是去找服务器上 /code/1.jpg 文件

/ 的重要性（不正确的配置会造成漏洞）
NGINX 是一个 Web 服务器，也可用作反向代理、负载平衡器、邮件代理和 HTTP 缓存。NGINX alias 指令定义了指定位置的替换。例如，在 /i/top.gif 的请求上使用以下配置：
```conf
location /i/ {
alias /data/w3/images/;
}
```
将发送文件 /data/w3/images/top.gif。但是，如果 /i../app/config.py 请求上的位置没有以目录分隔符（即 /）结尾：
```conf
location /i {
alias /data/w3/images/;
}
```
将发送文件 /data/w3/app/config.py。alias 的不正确配置可能会允许攻击者读取存储在目标文件夹外的文件。