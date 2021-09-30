ngx.say("URL错了哦<br>")
local cjson = require "cjson"
ngx.say(cjson.encode(router_filter).."<br>")
ngx.say(cjson.encode(test_var_exec_every_time).."<br>")
