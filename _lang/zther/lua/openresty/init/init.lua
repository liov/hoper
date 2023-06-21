local path = require "path"


--由于获取的路径为反斜杠(\)所以用上面的函数转为正斜杠(/)
work_dir = path.current_dir()
ngx.log(ngx.Inf, work_dir)
--添加搜索路径
package.path = string.format("%s?.lua;%s?/init.lua;%s", work_dir, work_dir, package.path)


router_filter = require("init.filter"):new("config")

test_var_exec_every_time = os.date("%c")