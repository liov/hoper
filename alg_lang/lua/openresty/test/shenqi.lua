ngx.say("hello")
ngx.flush() -- 显式的向客户端刷新响应输出
ngx.sleep(3)
