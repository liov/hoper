local function close_redis(red)
    if not red then
        return
    end
    -- 释放连接(连接池实现)，毫秒
    local pool_max_idle_time = 10000
    -- 连接池大小
    local pool_size = 100
    local ok, err = red:set_keepalive(pool_max_idle_time, pool_size)
    local log = ngx_log
    if not ok then
        log(ngx_ERR, "set redis keepalive error : ", err)
    end
end

local redis = require "resty.redis"
local red = redis:new()

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
ngx.say(IP_List..":<br><span style='width: 112px;display: inline-block;text-align: center'>IP地址</span><span>访问次数</span><br>")
local res, err = red:zrange(IP_List,0,-1,"WITHSCORES")
if not res then
    ngx.say("failed to get ip_list: ", err)
    return
end

for k, v in pairs(res) do
    if k%2 == 0 then
    	ngx.say("<span style='margin-left: 10px;color: #e96900'>",v,"<br>")
    else
    	 ngx.say("<span style='width: 120px;display: inline-block;text-align: center;color: #b854d4'>",v,"</span>")
    end
end

close_redis(red)
