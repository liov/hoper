local m = package.loaded[ngx.var.module]
m.handle()
ngx.say(test_var_exec_every_time.."<br>")
