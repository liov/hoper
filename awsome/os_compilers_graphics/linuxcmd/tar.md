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



filename.tar.bz2的解压:

tar -jxvf filename.tar.bz2

j: 　　bzip2　　　　　　　　　 压缩格式

其它选项和tar.gz解压含义相同



filename.tar.xz的解压:

tar -Jxvf filename.tar.xz
注意J大写



filename.tar.Z的解压:

tar -Zxvf filename.tar.Z
注意Z大写



关于tar的详细命令可以

tar --help

tar –xvf file.tar //解压 tar包
tar -xzvf file.tar.gz //解压tar.gz
tar -xjvf file.tar.bz2 //解压 tar.bz2
tar –xZvf file.tar.Z //解压tar.Z
unrar e file.rar //解压rar
unzip file.zip //解压zip
1、.tar 用 tar –xvf 解压
2、.gz 用 gzip -d或者gunzip 解压
3、.tar.gz和.tgz 用 tar –xzf 解压
4、.bz2 用 bzip2 -d或者用bunzip2 解压
5、.tar.bz2用tar –xjf 解压
6、.Z 用 uncompress 解压
7、.tar.Z 用tar –xZf 解压
8、.rar 用 unrar e解压
9、.zip 用 unzip 解压

压缩解压

XZ压缩最新压缩率之王
xz这个压缩可能很多都很陌生，不过您可知道xz是绝大数linux默认就带的一个压缩工具。
之前xz使用一直很少，所以几乎没有什么提起。
我是在下载phpmyadmin的时候看到这种压缩格式的，phpmyadmin压缩包xz格式的居然比7z还要小，这引起我的兴趣。
最新一段时间会经常听到xz被采用的声音，像是最新的archlinux某些东西就使用xz压缩。不过xz也有一个坏处就是压缩时间比较长，比7z压缩时间还长一些。不过压缩是一次性的，所以可以忽略。
xz压缩文件方法或命令
xz -z 要压缩的文件
如果要保留被压缩的文件加上参数 -k ，如果要设置压缩率加入参数 -0 到 -9调节压缩率。如果不设置，默认压缩等级是6.
xz解压文件方法或命令
xz -d 要解压的文件
同样使用 -k 参数来保留被解压缩的文件。
创建或解压tar.xz文件的方法
习惯了 tar czvf 或 tar xzvf 的人可能碰到 tar.xz也会想用单一命令搞定解压或压缩。其实不行 tar里面没有征对xz格式的参数比如 z是针对 gzip，j是针对 bzip2。
创建tar.xz文件：只要先 tar cvf xxx.tar xxx/ 这样创建xxx.tar文件先，然后使用 xz -z xxx.tar 来将 xxx.tar压缩成为 xxx.tar.xz
解压tar.xz文件：先 xz -d xxx.tar.xz 将 xxx.tar.xz解压成 xxx.tar 然后，再用 tar xvf xxx.tar来解包。

ZIP

zip可能是目前使用得最多的文档压缩格式。它最大的优点就是在不同的操作系统平台，比如Linux， Windows以及Mac OS，上使用。缺点就是支持的压缩率不是很高，而tar.gz和tar.gz2在压缩率方面做得非常好。闲话少说，我们步入正题吧：

我们可以使用下列的命令压缩一个目录：

# zip -r archive_name.zip directory_to_compress

下面是如果解压一个zip文档：

# unzip archive_name.zip

TAR

Tar是在Linux中使用得非常广泛的文档打包格式。它的好处就是它只消耗非常少的CPU以及时间去打包文件，他仅仅只是一个打包工具，并不负责压缩。下面是如何打包一个目录：
