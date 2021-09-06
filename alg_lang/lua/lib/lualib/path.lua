local _M = {}

function _M.conversion(value)
	if not string.find(value, "\\",1) then
		return value
	else
		return string.gsub(value,"\\","/")
	end
end

return _M
