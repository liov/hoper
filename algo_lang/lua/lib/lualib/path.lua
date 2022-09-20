local _M = {}

function _M.conversion(value)
	if not string.find(value, "\\",1) then
		return value
	else
		return string.gsub(value,"\\","/")
	end
end

function _M.current_dir()
    local info = debug.getinfo(2,"S")
    --获取当前路径
    local pathinfo = info.source
    --由于获取的路径为反斜杠(\)所以用上面的函数转为正斜杠(/)
    local filepath = string.match(_M.conversion(pathinfo),"^(.*/).*/.*$")
    return string.sub(filepath, 2, -1)
end

return _M
