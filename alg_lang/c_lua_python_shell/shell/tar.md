tail
grep
ps

filename.zip的解压:

unzip filename.zip


filename.tar.gz的解压:

tar -zxvf filename.tar.gz
其中zxvf含义分别如下

z: 　　gzip  　　　　　　　　    压缩格式

x: 　　extract　　　　　　　　  解压

v:　　 verbose　　　　　　　　详细信息

f: 　　file(file=archieve)　　　　文件

c: compress 压缩


关于tar的详细命令可以

tar --help

tar –xvf file.tar //解压 tar包  
tar –cvf file.tar //压缩 tar包
tar -zxvf file.tar.gz //解压tar.gz
tar -jxvf file.tar.bz2 //解压 tar.bz2
tar –Zxvf file.tar.Z //解压tar.Z
tar -Jxvf filename.tar.xz //解压tar.xz
unrar e file.rar //解压rar
unzip file.zip //解压zip

1、.tar 用 tar –xvf 解压
2、.gz 用 gzip -d或者gunzip 解压
3、.tar.gz和.tgz 用 tar –zxf 解压
4、.bz2 用 bzip2 -d或者用bunzip2 解压
5、.tar.bz2用tar –xjf 解压
6、.Z 用 uncompress 解压
7、.tar.Z 用tar –xZf 解压
8、.rar 用 unrar e解压
9、.zip 用 unzip 解压

--remove-files 删除原文件
ZIP

zip可能是目前使用得最多的文档压缩格式。它最大的优点就是在不同的操作系统平台，比如Linux， Windows以及Mac OS，上使用。缺点就是支持的压缩率不是很高，而tar.gz和tar.gz2在压缩率方面做得非常好。闲话少说，我们步入正题吧：

我们可以使用下列的命令压缩一个目录：

# zip -r archive_name.zip directory_to_compress

-m 删除原文件

下面是如果解压一个zip文档：

# unzip archive_name.zip

TAR

Tar是在Linux中使用得非常广泛的文档打包格式。它的好处就是它只消耗非常少的CPU以及时间去打包文件，他仅仅只是一个打包工具，并不负责压缩。下面是如何打包一个目录：

tar -zcvf test.tar.gz ./test/


压缩某个字符开头的 2010开头

tar -zcvf test.tar.gz 2010*/