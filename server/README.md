一锅端并不可取
不能为了学语言把相同的东西用不同语言实现一遍，这是不可取的，且无意义，至少一个人做不到
找到不同语言的适用点并且专注
虽然我最初设想的是不同语言构成的后端互通，但是这是一个伪技能和伪需求，这得多大的公司才能用到这么多语言

go java就应该做互联网
js ts 就应该做前端
c rust就应该在性能敏感的地方
python 专注于数据和脚本

lua openresty ,游戏热更新,嘿嘿嘿，k8s ingress上了openresty，lua目前也实在没理由用，就看有没有机会搞游戏了
cpp csharp就应该做GUI 虽然csharp如此诱人，但我可能精力不够了 如果我有机会做游戏的话+lua
dart...再说

k8s ingress controller 从openresty又切回了nginx
不过openresty主要作为网关，在上istio+k8s之后应用场景被替代了，虽说也可作为controller，
但是目前折腾意义不大，所以目前没有动态语言提供服务的场景，我能想到的就是python做数据分析，
然而这块我也不熟，先留着，目前也没必要用openresty自己实现一个ingress，有动态服务需求的可以
用js

lua 前端
go kotlin 业务
rust 基础设施