# go 编译动态库
set GOOS=linux 
set GOARCH=amd64
// go 动态库
go install -buildmode=shared -linkshared  std
go build -buildmode=shared -linkshared demo
// c abi
go build -buildmode=c-shared -o libgobblob.so

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

# go编译mod下的包
go build -o out modname/package