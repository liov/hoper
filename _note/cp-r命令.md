`cp -r dir1 dir2`
dir1目录下的内容有时候会被复制到dir2下，有时候会被复制到dir2/dir1下,取决于执行命令时,dir2是否已经存在,如果已经存在会被复制到dir2/dir1下
此时可以用 `cp -r dir1/* dir2` 解决