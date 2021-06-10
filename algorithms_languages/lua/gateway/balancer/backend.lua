local balancer = require "ngx.balancer"
local host = {"192.168.1.111", "192.168.1.112"}
local backend = ""
local port = ngx.var.server_port
local remote_ip = ngx.var.remote_addr
local key = remote_ip..port
local hash = ngx.crc32_long(key);
hash = (hash % 2) + 1
backend = host[hash]
ngx.log(ngx.DEBUG, "ip_hash=", ngx.var.remote_addr, " hash=", hash, " up=", backend, ":", port)
local ok, err = balancer.set_current_peer(backend, port)
if not ok then
    ngx.log(ngx.ERR, "failed to set the current peer: ", err)
    return ngx.exit(500)
end
ngx.log(ngx.DEBUG, "current peer ", backend, ":", port)