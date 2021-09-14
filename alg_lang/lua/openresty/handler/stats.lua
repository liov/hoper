
local _M={_VERSION = 0.1}

function _M.handler()
    local redis = require "redis"
    local red = redis.new()
    local IP_List = "IP_List"

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
end