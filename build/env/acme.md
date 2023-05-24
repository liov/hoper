

#curl https://get.acme.sh | sh -s email=liovx@qq.com
#~/.acme.sh/acme.sh --issue -d ${host} --standalone
#~/.acme.sh/acme.sh --issue --dns dns_cf -d ${host}

#~/.acme.sh/acme.sh --issue --dns -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please
#~/.acme.sh/acme.sh --renew -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please

# standalone
```bash
docker run --rm  -it  \
-v "$PWD/.acme.sh":/acme.sh  \
--net=host \
neilpang/acme.sh  --issue -d ${host}  --standalone
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