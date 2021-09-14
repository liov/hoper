local function new()
    local redis = require "resty.redis"
    local red = redis:new()

    red:set_timeout(1000)

    local ok, err = red:connect("127.0.0.1", 6379)
    if not ok then
        ngx.say("failed to connect: ", err,"<br>")
        return nil
    end
    return red
end

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