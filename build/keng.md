## 用其他主机docker login登录Harbor仓库报错
```bash
Error response from daemon: Get https://192.168.30.24/v2/: dial tcp 192.168.30.24:443: connect: connection refused
 vim /etc/docker/daemon.json
{
        "registry-mirrors": ["http://hoper.xyz"],
        "insecure-registries": ["192.168.xx.xx"]
}
restart docker
```
## Error loading config file XXX.dockerconfig.json - stat /home/XXX/.docker/config.json: permission denied
```
    这是因为docker的文件夹的权限问题导致的，处理办法如下，执行：
    
    sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
    
    sudo chmod g+rwx "/home/$USER/.docker" -R
```

## Temporary failure in name resolution 错误
```bash
/etc/hosts
127.0.0.1       localhost.localdomain localhost
vim /etc/resolv.conf
nameserver   xxx
nameserver   xxx
```

## IDEA中总模块名与java中maven模块名冲突
改总模块名

## Java搞了半天缺依赖
```$xslt
pom中只有test，少了
 <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
 </dependency>
```
## springboot管理普通类
@Component，@Autowired，@PostConstruct，init()

## java调go 远程主机强迫关闭了一个现有的连接。
建channel的时候少了usePlaintext()