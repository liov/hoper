--白名单列表
local whitelist = {
    module={
        'user',
        'user/login',
        'user/register'
    },
    direct = {
        'test',
        'log',
        'error'
    }
}
--路由重写列表
local rewritelist = {
    ['user/([-_a-zA-Z0-9]+)/login'] = 'user',
    ['user/([a-zA-Z0-9]+)/register'] = 'user/register',
    ['user/([a-zA-Z0-9]+)/logout'] = 'user/logout'
}
return {
    whitelist = whitelist,
    rewritelist = rewritelist
}