TLS 证书生成方法
生成私钥 ssl.key：

openssl genrsa -out ssl.key 2048
如果要对私钥进行传输/备份，建议先对私钥进行密码加密：

openssl rsa -in ssl.key -des3 -out encrypted.key
利用私钥就可以生成证书了。OpenSSL使用x509命令生成证书。这里需要区分两个概念：证书(certificate)和证书请求(certificate sign request)

证书是自签名或CA签名过的凭据，用来进行身份认证
证书请求是对签名的请求，需要使用私钥进行签名
x509命令可以将证书和证书请求相互转换，不过我们这里只用到从证书请求到证书的过程

从私钥或已加密的私钥均可以生成证书请求。生成证书请求 ssl.csr 的方法为：

openssl req -config openssl.cnf -new -key ssl.key -out ssl.csr
如果私钥已加密，需要输入密码。req命令会通过命令行要求用户输入国家、地区、组织等信息（也可以直接指定命令参数提供），这些信息会附加在证书中展示给连接方。

接下来用私钥对证书请求进行签名。根据不同的场景，签名的方式也略有不同，一般有两种

自签
自签名的原理是用私钥对该私钥生成的证书请求进行签名，生成证书文件。该证书的签发者就是自己，所以验证方必须有该证书的私钥才能对签发信息进行验证，所以要么验证方信任一切证书，面临冒名顶替的风险，要么被验证方的私钥（或加密过的私钥）需要发送到验证方手中，面临私钥泄露的风险。

当然自签名也不是一无用处，比如需要内部通讯的两台电脑需要使用加密又不想用第三方证书，可以在两端使用相同的私钥和证书进行验证（当然这样就跟对称加密没有区别了）
自签名的方法为：

openssl x509 -req -in ssl.csr -signkey ssl.key -out ssl.crt
同样如果ssl.key已加密，需要输入密码。

配置简单的 nginx tls 证书可以使用这个方式
CA 签
CA证书是一种特殊的自签名证书，可以用来对其它证书进行签名。这样当验证方选择信任了CA证书，被签名的其它证书就被信任了。在验证方进行验证时，CA证书来自操作系统的信任证书库，或者指定的证书列表。K8S 集群RBAC权限认证就使用这种模式。

生成自签名证书的方法为：

openssl x509 -req -in sign.csr -extensions v3_ca -signkey sign.key -out sign.crt
使用CA证书对其它证书进行签名的方法为：

openssl x509 -req -in ssl.csr -extensions v3_usr -CA sign.crt -CAkey sign.key -CAcreateserial -out ssl.crt