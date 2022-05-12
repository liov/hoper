# Docker Entrypoint & Cmd
先回顾下CMD指令的含义，CMD指令可以指定容器启动时要执行的命令，但它可以被docker run命令的参数覆盖掉。

ENTRYPOINT 指令和CMD类似，它也可用户指定容器启动时要执行的命令，但如果dockerfile中也有CMD指令，CMD中的参数会被附加到ENTRYPOINT 指令的后面。 如果这时docker run命令带了参数，这个参数会覆盖掉CMD指令的参数，并也会附加到ENTRYPOINT 指令的后面。这样当容器启动后，会执行ENTRYPOINT 指令的参数部分。

可以看出，相对来说ENTRYPOINT指令优先级更高。我们来看个例子，下面是Dockerfile的内容：

```dockerfile
#test
FROM ubuntu
MAINTAINER hello
RUN echo hello1 > test1.txt
RUN echo hello2 > /test2.txt
EXPOSE 80
ENTRYPOINT ["echo"]
CMD ["defaultvalue"]
```
假设通过该Dockerfile构建的镜像名为 myimage。

当运行 docker run myimage 输出的内容是 defaultvalue，可以看出CMD指令的参数得确是被添加到ENTRYPOINT指令的后面，然后被执行。
当运行docker run myimage hello world 输出的内容是 hello world ，可以看出docker run命令的参数得确是被添加到ENTRYPOINT指令的后面，然后被执行，这时CMD指令被覆盖了。
另外我们可以在docker run命令中通过 --entrypoint 覆盖dockerfile文件中的ENTRYPOINT设置，如：

`docker run --entrypoint="echo" myimage good  结果输出good`

注意，不管是哪种方式，创建容器后，通过 docker ps --no-trunc查看容器信息时，COMMAND列会显示最终生效的启动命令。

此外，很多的数据库软件的docker镜像，一般在entrypoint的位置会设置一个docker-entrypoint.sh文件，此文件位于/usr/local/bin位置，用于在容器初次启动的时候进行数据库的初始化操作。

# Kubernetes Command & args
下表总结了Docker和Kubernetes使用的字段名称：

Description	Docker field name	Kubernetes field name
The command run by the container	Entrypoint	command
The arguments passed to the command	Cmd	args
当你覆盖默认的Entrypoint和Cmd时，将应用以下规则：

如果不为容器提供command或args参数，则使用Docker镜像中定义的默认值。
如果提供command但没有提供args参数，则仅使用提供的command。Docker镜像中定义的默认EntryPoint和默认Cmd将被忽略。
如果仅为容器提供args，则Docker镜像中定义的默认Entrypoint将与您提供的args一起运行。
如果提供command和args，则将忽略Docker镜像中定义的默认Entrypoint和默认Cmd。 您的command与 args一起运行。
可以看到，k8s利用了Dockerfile的覆盖机制，使用command和args参数有选择性的覆盖了Docker镜像中的Entrypoint和Cmd启动参数，下面是一些例子：

Image Entrypoint	Image Cmd	Container command	Container args	Command run
[/ep-1]	[foo bar]	<not set>	<not set>	[ep-1 foo bar]
[/ep-1]	[foo bar]	[/ep-2]	<not set>	[ep-2]
[/ep-1]	[foo bar]	<not set>	[zoo boo]	[ep-1 zoo boo]
[/ep-1]	[foo bar]	[/ep-2]	[zoo boo]	[ep-2 zoo boo]

使用command和args的例子：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: command-demo
  labels:
    purpose: demonstrate-command
spec:
  containers:
  - name: command-demo-container
    image: debian
    command: ["printenv"]
    args: ["HOSTNAME", "KUBERNETES_PORT"]
  restartPolicy: OnFailure
```