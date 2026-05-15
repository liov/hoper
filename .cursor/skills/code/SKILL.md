---
name: protocol-and-dao-gen
description: 规范Protocol Proto修改和DAO/Model代码生成流程。当需要修改proto定义、新增数据库表DAO/Model、或更新protocol依赖时使用此规范。
---

# 代码生成规范

## Proto 修改流程

1. **修改位置**: 在proto目录下修改proto文件
2. **生成 go 代码**: 在 server/go 目录下执行 `protogen go -d -e -w -v -i ../../proto` 生成 go 代码
3. **生成 rust 代码**: 在 server/rust 目录下执行 `cargo run --bin rfv` 生成 rust 代码
4. **生成 dart 代码**: 在 client/app 目录下执行 `protogen dart -p ../../thirdparty/protobuf/_proto -i ../../proto -o lib/gen/pb` 生成 dart 代码
5. **生成 ts 代码**: 在 client/uniapp 目录下执行 `protogen ts -p ../../thirdparty/protobuf/_proto -i ../../proto -o gen/pb` 生成 ts 代码