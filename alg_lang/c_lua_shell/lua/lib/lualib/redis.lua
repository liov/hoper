local _M = {}

function _M.new()
    local redis = require "resty.redis"
    local red = redis:new()

      red:set_timeout(1000)

      local ok, err = red:connect("192.168.1.204", 6379)
      if not ok then
          ngx.log(ngx.ERR, "failed to connect: ", err)
          return nil
      end

      local count
      count, err = red:get_reused_times()
      if 0 == count then
          ok, err = red:auth("123456")
          if not ok then
              ngx.log(ngx.ERR, "failed to auth: ", err)
              return
          end
      elseif err then
          ngx.log(ngx.ERR, "failed to get reused times: ", err)
          return
      end
    return red
end

function _M.close_redis(red)
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
        log(ngx.ERR, "set redis keepalive error : ", err)
    end
end

return _M