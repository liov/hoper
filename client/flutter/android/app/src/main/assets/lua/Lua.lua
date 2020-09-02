local _lua = {}

_lua.getVersion = function (callback)
    callback(_VERSION)
end

_lua.updatePath = function (path)
    package.path = string.format("%s?.lua;%s?/init.lua;%s", path, path, package.path)
end

return _lua