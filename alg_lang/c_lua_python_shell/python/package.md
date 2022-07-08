# 1、pip freeze方法（不推荐）
如果你在写的项目是使用虚拟环境写的，就可以使用这个方法，因为这个方法会将你整个Python环境的包全把生成出来，如果你不是使用虚拟环境，使用这个方法，你会发现生成的文件，里面有很多你并不需要的包，这样吗使用安装的依赖包的时候会有很多不需要的包

终端使用命令：
pip freeze > requirements.txt


# 2、pipreqs第三方库（推荐）
使用 pipreqs 可以自动检索到当前项目下的所有组件及其版本，并生成 requirements.txt 文件，极大方便了项目迁移和部署的包管理。相比直接用pip freeze 命令，能直接隔离其它项目的包生成。

使用步骤：
1、先安装pipreqs库
pip install pipreqs

2、在当前目录使用生成
pipreqs ./ --encoding=utf8 --force


pip install -r requirements.txt