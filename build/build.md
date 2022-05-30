# go 编译动态库
set GOOS=linux 
set GOARCH=amd64
## go 动态库
go install -buildmode=shared -linkshared  std
go build -buildmode=shared -linkshared demo
## c abi
go build -buildmode=c-shared -o libgobblob.so

## go编译mod下的包
go build -o out modname/package
## 带时区信息
-tags timetzdata

## windows
windows需要安装gcc编译器，我用的的MinGW包，解压，把bin目录加入环境变量。

然后执行命令之后发现会报错，windows下go不支持生成动态库。

>>go build -buildmode=c-shared -o libgobblob.dll
-buildmode=c-shared not supported on windows/amd64
这一步折腾了好久，最终在stackoverflow找到了解决方法。[[ https://stackoverflow.com/questions/40573401/building-a-dll-with-go-1-7 | building-a-dll-with-go]]

编译静态库
go build -buildmode=c-archive -o libgobblob.a
gobblob.c文件，然后把go代码中要导出的函数，在gobblob.c中全部调用一遍。
```c
#include <stdio.h>
#include "libgobblob.h"

// force gcc to link in go runtime (may be a better solution than this)
void dummy() {
// 所有在go中要导出的代码都在这里调用一次，参数随便写，只要能编译ok即可
gobblob_init(NULL,NULL,NULL);
gobblob_deinit(NULL);
gobblob(NULL,NULL,NULL,NULL,NULL,NULL);
}

int main() {

}
```
执行如下命令，生成dll
gcc -shared -pthread -o libgobblob.dll gobblob.c libgobblob.a -lWinMM -lntdll -lWS2_32 -Iinclude
# 静态编译
go build -tags netgo
CGO_ENABLED=0 go build
go build -ldflags '-s -w --extldflags "-static -fpic"'
可选参数-ldflags 是编译选项：

-s -w 去掉调试信息，可以减小构建后文件体积，
--extldflags "-static -fpic" 完全静态编译，这样编译生成的文件就可以任意放到指定平台下运行，而不需要运行环境配置。
## 显然对于带CGO的交叉编译，CGO_ENABLED必须开启。
cgo的内部连接和外部连接
internal linking
internal linking的大致意思是若用户代码中仅仅使用了net、os/user等几个标准库中的依赖cgo的包时，cmd/link默认使用internal linking，而无需启动外部external linker(如:gcc、clang等)，不过由于cmd/link功能有限，仅仅是将.o和pre-compiled的标准库的.a写到最终二进制文件中。因此如果标准库中是在CGO_ENABLED=1情况下编译的，那么编译出来的最终二进制文件依旧是动态链接的，即便在go build时传入 -ldflags '-extldflags "-static"'亦无用，因为根本没有使用external linker

这样就会出现下文中命令行带参数-ldflags '-extldflags "-static"'，编译出来的还是会显示为动态连接。

external linking
而external linking机制则是cmd/link将所有生成的.o都打到一个.o文件中，再将其交给外部的链接器，比如gcc或clang去做最终链接处理。如果此时，我们在cmd/link的参数中传入 -ldflags '-linkmode "external" -extldflags "-static"'，那么gcc/clang将会去做静态链接，将.o中undefined的符号都替换为真正的代码。我们可以通过-linkmode=external来强制cmd/link采用external linker


docker run --rm -v /mnt/d/SDK/gopath:/go -v $PWD:/work -w /work/tools/server golang go build -ldflags '-linkmode "external" -extldflags "-static"' -o /work/build/tmp/main /work/tools/server/fileserver.go

# android
arm64 aarch64-linux-android
amd64 x86_64-linux-android

set GOOS=android
set GOARCH=arm64
set CGO_ENABLED=1
set AR=D:\SDK\Android\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android-ar.exe
set CC=D:\SDK\Android\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android30-clang.cmd
set CXX=D:\SDK\Android\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android30-clang++.cmd
//set CGO_LDFLAGS=-LD:\SDK\Android\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\x86_64-linux-android\lib -g -O2
go build -buildmode=c-shared -o libhello.so hello.go
# go 编译静态库
go build -buildmode=c-archive -o libgobblob.a

# '_debugLifecycleState != _ElementLifecycle.defunct': is not true.

You can copy paste run full code below
You can move _controller.dispose(); before super.dispose();
code snippet

```dart
@override
void dispose() {
_controller.dispose();
super.dispose();
}
```
working demo
