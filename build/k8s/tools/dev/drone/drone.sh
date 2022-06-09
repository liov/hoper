github oauth

https://drone.hoper.xyz
https://drone.hoper.xyz/login

openssl rand -hex 16

docker pull drone/drone:latest

配置
Drone 服务器使用环境变量进行配置。本文引用了配置选项的子集，定义如下。有关配置选项的完整列表，请参阅配置。

DRONE_GITHUB_CLIENT_ID
必需的字符串值提供您在上一步中生成的 GitHub oauth 客户端 ID。
DRONE_GITHUB_CLIENT_SECRET
必需的字符串值提供您在上一步中生成的 GitHub oauth 客户端密码。
DRONE_RPC_SECRET
必需的字符串值提供在上一步中生成的共享密钥。这用于验证服务器和运行器之间的 rpc 连接。必须为服务器和运行器提供相同的秘密值。
DRONE_SERVER_HOST
必需的字符串值提供您的外部主机名或 IP 地址。如果使用 IP 地址，您可以包括端口。例如drone.company.com.
DRONE_SERVER_PROTO
必需的字符串值提供您的外部协议方案。此值应设置为 http 或 https。如果您配置 ssl 或 acme，此字段默认为 https。如果您将 Drone 部署在负载均衡器或带有 SSL 终止的反向代理后面，则应设置此值。https

docker run \
  --volume=/var/lib/drone:/data \
  --env=DRONE_GITHUB_CLIENT_ID=your-id \
  --env=DRONE_GITHUB_CLIENT_SECRET=super-duper-secret \
  --env=DRONE_RPC_SECRET=super-duper-secret \
  --env=DRONE_SERVER_HOST=drone.company.com \
  --env=DRONE_SERVER_PROTO=https \
  --publish=80:80 \
  --publish=443:443 \
  --restart=always \
  --detach=true \
  --name=drone \
  drone/drone:latest