

local redis = require "redis"

local red = redis.new()

local IP_List = "IP_List"

local headers=ngx.req.get_headers()
local ip=headers["X-REAL-IP"] or headers["X_FORWARDED_FOR"] or ngx.var.remote_addr or "0.0.0.0"

local exist, err = red:zrank(IP_List,ip)
if not exist then
    red:zadd(IP_List,1,ip)
else
	red:zincrby(IP_List,1,ip)
end


redis.close_redis(red)

ngx.log(ngx.ERR,ngx.var.lua_path)
dofile(work_dir.."access/load_module.lua")
