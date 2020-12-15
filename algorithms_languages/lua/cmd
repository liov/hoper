链接库编译
使用Visual Studio创建一个VC++项目，项目命名为Lua53，项目类型为静态库、不设置预编译头。
删除Visual Studio自动创建的.cpp文件及其对应的.h文件。
将下载的Lua代码解压，将src目录下的全部文件拷贝到项目中，然后删除lua.c、luac.c和lua.hpp这三个文件。
编译项目会得到一个Lua53.lib的文件，这就是我们编译得到的Lua链接库。
Lua解释器
Lua解释器是一个可以直接运行Lua代码的可执行文件，因此

在同一个解决方案下继续创建VC++项目，项目命名为Lua，项目类型为控制台应用程序、需设置预编译头。
删除Visual Studio自动创建的.cpp文件及其对应的.h文件。
将下载的Lua代码解压，将src目录下的全部文件拷贝到项目中，然后删除luac.c这个文件。
设置当前项目依赖于Lua53项目
编译项目会得到一个Lua.exe文件，这就是我们编译得到的Lua解释器。
Lua编译器
在同一个解决方案下继续创建VC++项目，项目命名为Lua，项目类型为控制台应用程序、需设置预编译头。
删除Visual Studio自动创建的.cpp文件及其对应的.h文件。
将下载的Lua代码解压，将src目录下的全部文件拷贝到项目中，然后删除lua.c这个文件。
设置当前项目依赖于Lua53项目
编译项目会得到一个Luac.exe文件，这就是我们编译得到的Lua编译器。

cl /MD /O2 /c /DLUA_BUILD_AS_DLL *.c
ren lua.obj lua.o
ren luac.obj luac.o
link /DLL /IMPLIB:lua5.3.0.lib /OUT:lua5.3.0.dll *.obj
link /OUT:lua.exe lua.o lua5.3.0.lib
lib /OUT:lua5.3.0-static.lib *.obj
link /OUT:luac.exe luac.o lua5.3.0-static.lib

1.打开VS开发人员命令提示符
2.cd到源码的src目录
3.依次执行以下命令
cl /MD /O2 /c /DLUA_BUILD_AS_DLL *.c
ren lua.obj lua.o
ren luac.obj luac.o
link /DLL /IMPLIB:lua5.3.5.lib /OUT:lua5.3.5.dll *.obj
link /OUT:lua.exe lua.o lua5.3.5.lib
lib /OUT:lua5.3.5-static.lib *.obj
link /OUT:luac.exe luac.o lua5.3.5-static.lib
4.新建文件夹命名为lua，并在该文件夹中新建bin，include，lib三个文件夹
5.将src目录中生成的 lua.exe，lua5.3.5.dll，luac.exe 复制到lua/bin中
6.将src目录中生成的 lua5.3.5.exp，lua5.3.5.lib，lua5.3.5-static.lib，luac.exp，luac.lib 复制到lua/lib中
7.将 lauxlib.h，lua.h，lua.hpp，luacnf.h，lualib.h 复制到lua/include中

lua jit
cd src
msvcbuild

编译luajit和luasocket的大坑
去掉compat.h里的
```
#if LUA_VERSION_NUM==501

#ifndef _WIN32
#pragma GCC visibility push(hidden)
#endif

void luasocket_setfuncs (lua_State *L, const luaL_Reg *l, int nup);
void *luasocket_testudata ( lua_State *L, int arg, const char *tname);

#ifndef _WIN32
#pragma GCC visibility pop
#endif

#define luaL_setfuncs luasocket_setfuncs
#define luaL_testudata luasocket_testudata

#endif
```
这是一个兼容文件的头定义，jit确实是5.1的，然而据说5.2才有luaL_setfuncs和luaL_testudata
于是有人提了这个pr，然而我编译一直报错，lnkerr 找不到luasocket_setfuncs，luasocket_testudata
尝试了许久才灵机一动，去了这段，编译成功了,我只能认为jit2.1是有这俩函数的
我日