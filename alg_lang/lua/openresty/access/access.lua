

local redis = require "redis"
local red = redis.new()

--imoprt_path = dofile("./init.lua")
--imoprt_path()

red:set_timeout(1000)

local ok, err = red:connect("127.0.0.1", 6379)
if not ok then
    ngx.say("failed to connect: ", err,"<br>")
    return
end

local IP_List = "IP_List"

local headers=ngx.req.get_headers()
local ip=headers["X-REAL-IP"] or headers["X_FORWARDED_FOR"] or ngx.var.remote_addr or "0.0.0.0"

local exist, err = red:zrank(IP_List,ip)
if not exist then
    red:zadd(IP_List,1,ip)
else
	red:zincrby(IP_List,1,ip)
end


close_redis(red)
