local need_module = router_filter:route_verify()

local function load_module()
    local path = require "path"
    local info = debug.getinfo(1,"S")
    --获取当前路径
    local pathinfo = info.short_src
    --由于获取的路径为反斜杠(\)所以用上面的函数转为正斜杠(/)
    local filepath = string.match(path.conversion(pathinfo),"^(.*/).*/.*$")
    package.loaded[ngx.var.lua_path] = dofile(filepath .."/lua/"..ngx.var.lua_path..".lua")
end

local function get_module(uri)
    local ret =  string.match(uri,"[^/.]+")
    return ret
end
-- 这段逻辑是 lua_path是/lua/user，然后去找这个模块，加载后lua_path设为router，router作为统一的入口
if need_module then
    --模块不存在，加载模块
    if package.loaded[ngx.var.lua_path] == nil then
        load_module()
    else
        --重载参数为true，重载模块
        if ngx.var.arg_reload then
            package.loaded[ngx.var.lua_path] = nil
            load_module()
        end
    end
    ngx.var.module = get_module(string.gsub(ngx.var.lua_path,"/","."))
    ngx.var.lua_path = "router"
end
