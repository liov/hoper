package.path = "D:/hoper/server/dynamic/lua/lualib/?.lua;;"
package.cpath = "D:/hoper/server/dynamic/lua/clib/?.so;;"
local socket = require("socket")

host = host or "*"
port = port or 9000
if arg then
    host = arg[1] or host
    port = arg[2] or port
end
print("Binding to host '" ..host.. "' and port " ..port.. "...")
io.flush ()
s = assert(socket.bind(host, port))
i, p   = s:getsockname()
assert(i, p)
print("Waiting connection from talker on " .. i .. ":" .. p .. "...")
c = assert(s:accept())
print("Connected. Here is the stuff:")
io.flush ()
l, e = c:receive()
while not e do
    print(l)
    io.flush ()
    l, e = c:receive()
end
print(e)
