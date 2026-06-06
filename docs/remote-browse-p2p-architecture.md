# 远程相册浏览（NAT 内网 P2P + 中继）技术方案

本文档描述「远程相册」端到端架构：**信令、中继、P2P 数据面、缩略图/列举/读删均在 `server/rust/rb`**（`rb-daemon` + `rb` Agent + `librb`）。Flutter Viewer 经 **gRPC over HTTP/2** 访问 Agent；**Go `hoper` 主站不再注册 `/rb/*` 路由**（公网入口用 Nginx 反代至 `rb-daemon`）。

---

## 1. 目标与约束

| 维度 | 要求 |
|------|------|
| 客户端 | Flutter（`client/app`），可内嵌 H5（UniApp）做部分 UI |
| 远端 | Agent：`rb` + `RB_ROOM` 或 Flutter `rb_agent_run` |
| 网络 | NAT 后需打洞；失败走中继；双端 IPv6 时优先直连 |
| 传输 | **Protobuf + gRPC**；并发由 HTTP/2 多路复用提供 |
| 服务端 | Rust `rb-daemon`（信令+中继）+ `rb`（Agent） |
| 数据编码 | `proto/remotebrowse/` |
| 非目标 | 不做实时屏幕镜像；中继不做重画质转码 |

---

## 2. 总体架构

```
┌─────────────┐   WSS 信令 + P2P gRPC(h2)   ┌──────────────────────────────┐
│ Flutter     │◄────────────────────────────►│ rb-daemon（/rb/signal、relay）│
│  Viewer     │  直连/中继: h2c on TCP      │ rb Agent（MediaSvc + gRPC）   │
│             │  ICE: 帧隧道 + h2 gRPC(FFI)  └──────────────────────────────┘
└─────────────┘
```

- **信令面**：`rb-daemon` `/rb/signal`；可选 Go `file` 反代 WebSocket。
- **数据面**：建立字节流后跑 **HTTP/2 cleartext（h2c）+ `RemoteBrowseService` gRPC**。

---

## 3. 连接策略

1. **局域网 / 手填 IP TCP 直连**（默认 `19091`）→ `RbBoundSocketConnector` → gRPC
2. **ICE 打洞**（UDP 上 `[4B len][payload]` 可靠帧 → duplex → h2 gRPC；Viewer 经 `rb_ice_grpc_*` FFI）
3. **TCP 中继**（`RBRL` 握手后 **透明 TCP 桥**，不再长度前缀包装应用帧）→ 同上 gRPC

**STUN**：`stun:stun.l.google.com:19302`。后续可评估 ICE 上真 HTTP/3（需 tonic 升级）。

---

## 4. 传输与 RPC

| 链路 | 字节通道 | 应用层 |
|------|----------|--------|
| TCP 直连 / TCP 中继 | 裸 `TcpStream` | h2c + gRPC |
| ICE | `ice_link` 分块帧 → `duplex` | h2c + gRPC（Agent 内建；Viewer FFI） |

**RPC**（`browse.service.proto`）：`ListFiles`、`GetThumbnail` / `GetThumbnailBatch`、`ReadFile`（server stream）、`DeleteFile`、`PipeFile` 等。业务实现集中在 **`remotebrowse_svc::MediaSvc`**（Agent TCP/ICE、gRPC 调试服务共用）。

**并发**：取消旧版 wire 单队列；多缩略图/读文件可并行多 HTTP/2 stream。

---

## 5. 信令（rb-daemon）

- Room 配对、ICE candidate 交换（无 WebRTC SDP 媒体语义）。
- 下发中继 `host:port` + `sessionId`；Viewer/Agent 发 `RBRL` 后进入透明转发。

---

## 6. 中继要点

- **`daemon/relay`**：`copy_bidirectional` 透明 TCP，承载双方 h2 握手与 gRPC。
- 限流、TTL、可观测性同前。

---

## 7. Rust 模块（`server/rust/rb`）

| 模块 | 职责 |
|------|------|
| `remotebrowse_svc` | 列表/缩略图/读删/沙箱路径 |
| `p2p_grpc` | duplex/TCP 上 `serve_h2` / `connect_h2` |
| `ice_link` | ICE `Conn` ↔ duplex ↔ gRPC |
| `ice_grpc_ffi` | Viewer ICE gRPC FFI |
| `grpc_server` | 本机 axum + gRPC 调试 |
| `transport/agent` | TCP/ICE 入站接 MediaSvc |

**已移除**：自研 wire 帧协议（`wire_agent` / `wire_client` / `tcp_wire`）。

---

## 8. Flutter（`client/app`）

- **`RbViewerSession`**：直连 / ICE FFI / 中继 → **`RbGrpcSession`**
- **`rb_grpc_transport`**：已连接 `Socket` → HTTP/2 connector
- **`RbIceGrpcBridge`**：ICE 路径 Rust FFI 返回 protobuf 字节
- **`wire_codec.dart`**：仅保留 `rbRelayJoinBytes` / `rbRoleViewer|Agent`
- UI：`remote_browse_files_view` 等使用 `_grpc` 会话

**构建**：改 Rust FFI 后执行 `server/rust/rb/scripts/build_flutter_lib.sh`，产物在 `client/app/staticLibs/`（iOS 经 `rb.xcconfig` 静态链 Runner，Android 经 CMake 将 `.a` 链入 `librb.so`）。

---

## 9. 安全

- 信令 TLS；路径沙箱；审计元数据访问。

---

## 10. 分阶段

| 阶段 | 交付 |
|------|------|
| M1 | daemon + MediaSvc + proto |
| M2 | 直连/中继/ICE 全走 gRPC h2 |
| M3 | 相册 UX + 批量缩略图 |
| M4 | UniApp 嵌入 |

*与 `proto/remotebrowse/`、`server/rust/rb`、`client/app/lib/remotebrowse/` 保持一致。*
