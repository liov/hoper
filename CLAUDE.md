# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

**hoper** is a multi-platform social media application with:
- **Backend**: Go primary service + Rust modules (video processing, WebSocket messaging)
- **Frontend**: UniApp+Vue3 cross-platform (iOS/Android/Web/mini-programs) + Flutter native app
- **Shared**: Protocol Buffer definitions in `proto/` consumed by all clients and servers
- **Libraries**: 9 Git submodules under `thirdparty/` (all internal hopeio libraries)


在 proto 目录下定义, 服务端通过 grpc 交互,客户端通过 http 调用, flutter 客户端通过 protobuf 编码与服务端交互, web 客户端通过 json 编码与服务端交互

## Submodule Setup

```bash
git submodule update --init --recursive --remote
```

Key submodules: `cherry` (microservice framework), `gox` (stdlib extensions), `initialize` (service init), `pick` (web utilities), `scaffold` (OpenTelemetry), `protobuf` (protogen tool), `diamond` (Vue3 components).

## Go Backend (`server/go/`)

```bash
cd server/go

# Install protobuf generation tools (first time)
go run $(go list -m -f {{.Dir}} github.com/hopeio/protobuf)/tools/install_tools.go

# Generate code from proto files
protogen go -d -e -w -v -i ../../proto -o protobuf

# Run server
go run main.go -c config/config.toml

# Tests
go test ./...

# Docker build
protogen go -d -e -w -v -i ../../proto -o protobuf
GOOS=linux go build -trimpath -o ../../build/hoper main.go
cd ../../build && docker build -t jybl/hoper . && docker push jybl/hoper
```

Go 1.25+. Framework: Cherry (gRPC + HTTP gateway + GraphQL). ORM: GORM. Observability: OpenTelemetry + Prometheus. Config is TOML-based; `config/local.toml` overrides `config/config.toml`.

## Rust Backend (`server/rust/`)

Two services:
- `rfv/` — video/image processing (ffmpeg, axum, tokio); exposes staticlib + cdylib for FFI
- `message/` — WebSocket messaging (axum-ws, dashmap)

```bash
cd server/rust/rfv && cargo build --release
cd server/rust/message && cargo build --release
```

Rust edition 2024.

## Frontend (`client/`)

### UniApp (cross-platform app + mini-programs)
```bash
cd client/uniapp
pnpm install         # requires Node >=18, pnpm >=9

pnpm dev             # H5 dev server
pnpm dev:app         # Native app (iOS/Android)
pnpm dev:mp-weixin   # WeChat mini-program
pnpm build:h5        # H5 production build
pnpm build:app       # App production build
```

### Web (`client/web/`) — Vue3 + Vite + TypeScript + WASM
```bash
cd client/web
yarn install
yarn serve           # dev server
yarn build           # production build
yarn test:unit       # unit tests
yarn lint            # ESLint
```

## Proto Definitions (`proto/`)

Services: `user/`, `content/`, `message/`, `file/`, `common/`. After editing `.proto` files, regenerate Go code with the `protogen` command above. The `protogen` flags:
- `-d` OpenAPI docs
- `-e` enum extensions
- `-w` gin HTTP gateway
- `-v` validation

## Architecture

**API layer**: gRPC primary; HTTP via gRPC-Gateway; GraphQL experimental.
**Deployment**: Kubernetes + APISIX API gateway. Config in `deploy/`; k8s manifests in `server/go/svc.yaml`; routing rules in `server/go/route.yaml`.
**CI/CD**: Drone, configured via `.drone.js` (multi-host build and deploy targets: localhost, mint, tx).

## Thirdparty Submodules Location

All internal libraries live under `thirdparty/` and are imported in `server/go/go.mod` via `replace` directives pointing to local paths.
