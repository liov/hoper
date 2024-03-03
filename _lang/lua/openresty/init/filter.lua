local _M = { _VERSION = '0.1' }
local mt = { __index = _M}
local tinsert = table.insert
local tconcat = table.concat
local cjson = require "cjson"
function _M.new(self,config_path)
    local tab = {}
    local routeMap = require(config_path)
    local route_data = {module={},direct={}}
    local whitelist = routeMap.whitelist
    for i=1,#whitelist["module"] do
        tinsert(route_data["module"],tconcat({'^',whitelist["module"][i],'$'}))
    end
    for i=1,#whitelist["direct"] do
        tinsert(route_data["direct"],tconcat({'^',whitelist["direct"][i],'$'}))
    end
    local rewritelist = routeMap.rewritelist
    local rewrite_data = {}
    local rewrite_urls = {}
    local x = 1
    for k,v in pairs(rewritelist) do
        tinsert(rewrite_data,tconcat({'^(?<index',x,'>',k,')$'}))
        tinsert(rewrite_urls,v)
         x = x + 1
    end
    tab.rewrite_urls = rewrite_urls
    tab.rewrite_pattern = tconcat(rewrite_data,'|')
    tab.route_module = tconcat(route_data["module"],'|')
    tab.route_direct = tconcat(route_data["direct"],'|')
    return setmetatable(tab, mt)
end

function _M.route_verify(self)
    local lua_path = ngx.var.lua_path
    local m = ngx.re.match(lua_path,self.route_module)
    if not m then
        m = ngx.re.match(lua_path,self.route_direct)
    else
        return true
    end
    if not m then
        m = ngx.re.match(lua_path,self.rewrite_pattern)
        if not m then
            ngx.var.lua_path = "error"
        else
            --lua表长度的大坑
            --local locant = ngx.re.match(next(m,#m), "^index(\\d+)","Djo")
            for i=1,#self.rewrite_urls,1 do
                if m['index'..i] then
                    ngx.var.lua_path = self.rewrite_urls[i]
                    if ngx.re.match(ngx.var.lua_path,self.route_module) then
                        return true
                    end
                    break
                end
            end
        end
    end
end

return _M

--莫名其妙的bug
-- set $lua_path $1;
-- ngx.var.lua_path = nil
--ngx.var.1 error
--ngx.var[1]
--原来是windows结束进程，其实后台还有无数进程