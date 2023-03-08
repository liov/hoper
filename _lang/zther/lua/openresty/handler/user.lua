
local _M={_VERSION = 0.1}

function _M.handler()
    if ngx.var.arg_method then
        local fun =_M[ngx.var.arg_method];
        if fun then
            fun();
        end
    else
        ngx.say("接口不存在<br>")
    end
end

function _M.test()
   ngx.say("test<br>")
end

return _M
