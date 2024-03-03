setmetatable(_M, {__index = function(self, cmd)
    local method =
    function (self, ...)
        return _do_cmd(self, cmd, ...)
    end

    -- cache the lazily generated method in our
    -- module table
    _M[cmd] = method
    return method
end})

--这段代码的精妙之处，__index是一个函数，当搜索表找不到属性时，会执行__index函数生成函数
