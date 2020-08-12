local _lua = {}

_lua.getVersion = function (callback)
    callback(_VERSION)
end

return _lua