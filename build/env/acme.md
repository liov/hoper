

#curl https://get.acme.sh | sh -s email=lby.i@qq.com
#~/.acme.sh/acme.sh --issue -d ${host} -d *.${host} --standalone --httpport 88


#~/.acme.sh/acme.sh --issue --dns -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please
#~/.acme.sh/acme.sh --renew -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please

# standalone
```bash
docker run --rm  -it  \
-v "$PWD/acme":/acme.sh  \
--net=host \
neilpang/acme.sh  --issue -d ${host}  --standalone --standalone --httpport 88
```
# dns
```bash
docker run --rm  -it  \
-v $PWD/acme:/acme.sh  \
--net=host \
neilpang/acme.sh --issue --dns -d ${host}  \
--yes-I-know-dns-manual-mode-enough-go-ahead-please
```
```bash
docker run --rm  -it  \
-v $PWD/acme:/acme.sh  \
--net=host  \
neilpang/acme.sh --renew  -d ${host} \
--yes-I-know-dns-manual-mode-enough-go-ahead-please
```

# zerossl连接失败更改letsencrypt
`acme.sh --set-default-ca --server letsencrypt`

```bash
docker run --rm  -it  \
-v $PWD/acme:/acme.sh  \
--net=host  \
neilpang/acme.sh --set-default-ca --server letsencrypt
```