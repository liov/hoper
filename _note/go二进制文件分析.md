// 查看占用大小前20的包
go tool nm -sort size -size $bin |head -n 20

## go version 命令
可以先尝试使用 go version 命令获取二进制文件中包含的依赖包信息：

`go version -m $bin`

## strings 命令
也可以使用 strings 命令获取二进制文件中的信息:
`strings $bin |grep github.com`

## go tool nm 命令
go tool nm 命令也可以得到相关信息:
`go tool nm $bin |grep github.com`

## go tool objdump 命令
go tool objdump 命令也可以间接得到相关信息:
`go tool objdump $bin |grep github.com`

## redress 工具
redress 是一个专门用于分析 Go 二进制可执行文件的开源软件，通过这个工具也可以得到想要的包信息：
`redress -pkg -filepath -vendor -unknown $bin`


# 分析工具
go install https://github.com/jondot/goweight@master