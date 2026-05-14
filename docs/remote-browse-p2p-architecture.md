# 远程相册浏览（NAT 内网 P2P + 中继）技术方案

本文档描述在 **hoper** 仓库内实现「手机端远程浏览另一台设备上的文件」的端到端架构：信令与编排代码落在 **`server/go/webrtc`**（**仅** `pion/ice` 打洞 + **quic-go**，**不采用 WebRTC**）；媒体解码与缩略图流水线落在 **`server/rust/rfv`**；信令走公网、数据面尽量 P2P、失败走中继；强调 **省带宽**、**IPv6 优先**；与 **Cherry/gRPC、proto、Flutter、UniApp** 等现有规范对齐。

---

## 1. 目标与约束

| 维度 | 要求 |
|------|------|
| 客户端 | Flutter（`client/app`），可内嵌 H5（UniApp，`client/uniapp`）做部分 UI 或运营页 |
| 远端 | 同构或轻量 Agent（建议 Flutter + 原生后台 / 独立小进程均可） |
| 网络 | 双方多在 NAT 后；需打洞；失败走中继；**双端具备可用 IPv6 时优先 IPv6 直连** |
| 传输 | 打洞成功后自建应用层协议；**仅传文件元数据 + 原始/压缩字节**；缩略图与大图在客户端解码渲染 |
| 服务端 | 公网信令；Go 负责 HTTP/gRPC/会话；Rust 负责重解码、转码、缩略图流水线 |
| 数据编码 | 业务消息 **Protobuf**（定义于 `proto/`） |
| 体验 | 相册式缩略图网格；点击进入预览；上下滑切文件 |
| 非目标 | 不做实时屏幕镜像；不在中继上做重画质实时转码（避免成本与延迟） |

---

## 2. 总体架构

```
┌─────────────┐     HTTPS/WSS/gRPC      ┌──────────────────────────────┐
│ Flutter/H5  │◄──────────────────────►│ Go file 门面（/rb/health、list、thumb）│
│  client     │   注册/配对/ICE 交换     │ 信令可反代至 rfv-daemon               │
└──────┬──────┘                          └──────────────┬───────────────┘
       │                                              │
       │ 可选：Go↔Rust 控制面 gRPC 双向流              │
       │                                              ▼
       │                                 ┌──────────────────────────────┐
       │                                 │ Rust rfv（媒体 + 信令 daemon） │
       │                                 │ 缩略图、ICE FFI、/rb/signal    │
       │                                 └──────────────────────────────┘
       │
       │  P2P（首选）或 中继（兜底）
       ▼
┌─────────────┐ ◄════ 加密应用流 ════► ┌─────────────┐
│ Viewer 端   │                        │ Agent 端    │
│ (浏览方)    │                        │ (被控方/库)  │
└─────────────┘                        └─────────────┘
```

- **信令面**：公网 **`rfv-daemon`**（`server/rust/rfv`，`/rb/signal` + TCP 中继）；Go **`server/go/file`** 通过 `RB_SIGNAL_UPSTREAM` 反代 WebSocket，并暴露 `/rb/health`、`/rb/v1/list`、`/rb/v1/thumb`（列表/缩略图经 gRPC 调 `rfv`）。
- **数据面**：在 **ICE 选定后的 UDP 四元组** 上跑 **QUIC（quic-go）** 或经 **自研中继** 的 QUIC/TLS；之上承载 **文件列表、缩略图请求、原文件分片** 的 **自定义二进制协议**（建议 **TLS1.3 / Noise** 或 QUIC 自带加密，中继仅见密文）。

---

## 3. 连接策略（IPv6 → ICE 打洞 → 中继）

建议实现为**显式状态机**，顺序如下（可在实现中并行探测，但对外表现仍按优先级收敛）：

1. **手填 IP / 局域网 TCP 直连**（`RB_DIRECT_TCP`，默认 `19091`）
2. **房间码 + QUIC 房间**（`RB_QUIC_PORT`，默认 `19092`；`quic-go` 调试证书，经信令交换 `peer_endpoints`）
3. **ICE 打洞**（`pion/ice` / `webrtc-ice`；`RB_RAW_ICE=1` 时裸 wire）
4. **自研 TCP 中继**（`RBRL` 握手 + 长度前缀帧；`rfv-daemon` 与 Go 内嵌 relay 协议一致）

**STUN/TURN**：`pion/ice` 配置 `ICEServers`；严苛 NAT 可叠加 **TURN（如 coturn）** 作为 **UDP 中继**，仍可在该 UDP 通道上跑 QUIC（需注意 MTU/中间盒）；与「自研应用层中继」二选一或分级：优先 ICE 直连，其次 TURN，最后自研中继。

---

## 4. 技术选型说明（含可替代方案）

### 4.1 为何采用「pion/ice + QUIC」而非 WebRTC

- 本需求只要 **NAT 打洞 + 可靠多路传输**，**QUIC** 即可提供加密、多流、拥塞控制；无需 **SDP、DTLS-SRTP、SCTP DataChannel** 整条 WebRTC 栈。  
- **代价**：`pion/ice` 与 **quic-go** 共用 UDP 的 glue 层需自行打通（读写泵、关闭顺序、ICE restart 与 QUIC 0-RTT/会话恢复的协调），工作量高于「直接 DataChannel」，但与「不要 WebRTC」的约束一致。  
- **可选替代（若 ICE+QUIC glue 风险过高）**：在 **不改变目录约束** 的前提下，仍可将数据面改为 **ICE + 裸 UDP 帧 + 自定义 ARQ**（省依赖、费心力）；或评估 **libp2p（Rust，放在 rfv 或 sidecar）** 与 Go 信令拆分——与当前「Go 仅 webrtc 目录 + Rust 仅 rfv」强约束冲突时再单独开目录评审。

### 4.2 开源组件建议

| 能力 | 建议库 | 说明 |
|------|--------|------|
| ICE 打洞 | `github.com/pion/ice/v4`（版本以 go.mod 为准） | Agent、候选 trickle、STUN/TURN；**不用** `pion/webrtc` |
| 传输 | `github.com/quic-go/quic-go` | IPv6 直连与 ICE 成功后的 P2P 数据面；多 stream 传文件/控制 |
| 中继传输 | 自研 + `quic-go` 或 `crypto/tls` | 实现放在 `server/go/webrtc` 下子包（如 `relay/`） |
| 信令辅助 | `webrtc/service/signal` | 可演进为 **WS/HTTP** 上的 ICE JSON 交换，替代示例里的 stdin |

---

## 5. 信令协议（公网 Rust rfv-daemon + Go file 门面）

职责：

- 设备注册、用户鉴权（与主服务集成时可 HTTP 反向代理到 Cherry；**`/rb/*` 业务路由**在 `server/go/file`）。
- **Room / Pairing**：短码或扫码绑定 Viewer ↔ Agent。
- **交换 ICE**：`ufrag`/`pwd`、**trickle candidate**、可选 `end-of-candidates`；**不交换 SDP Offer/Answer（WebRTC 语义）**。
- 下发 **中继地址**、**会话密钥材料**（ECDH 等，中继不解密）。
- 可选：带宽档位、就近 relay 列表。

实现位置：

- **`server/rust/rfv`**：`daemon/signal`（Hub）、`daemon/relay`；二进制 **`rfv-daemon`**（`RB_HTTP`、`RB_RELAY_TCP`）。
- **`server/go/file/remotebrowse`**：gRPC/HTTP 门面；`RB_SIGNAL_UPSTREAM` 时 `/rb/signal` 反代至 daemon。
- **`server/go/webrtc/session`**：`rbagent` / `rbviewer` CLI 与打洞/QUIC 房间逻辑（**不再**注册 `/rb/signal`）。
- Proto：`proto/remotebrowse/`；`protogen` 生成 Go；Flutter 用同套 proto。

---

## 6. P2P 应用层协议（省带宽）

在 **QUIC stream** 上使用统一 **帧格式**（示意）：

```
[1 byte version][1 byte type][4 byte payload_len][payload bytes]
type: HELLO | FILE_INDEX | THUMB_REQ | THUMB_DATA | FILE_CHUNK | ACK | WINDOW_UPDATE | ...
```

- **元数据**：`FILE_INDEX` 使用 protobuf（见第 8 节）；**列表分页** + **字段裁剪**（不要一次下发全库 exif）。
- **缩略图**：固定最大边（如 128/256），**WebP** 或 **JPEG q=动态**；Agent 侧由 **`server/rust/rfv`** 内模块生成（复用/扩展 `rfv/src/file.rs` 的 FFmpeg、image 管线），缓存磁盘 LRU + **内容寻址 hash**，相同文件不重复传。
- **大图/原图**：**Range 分片** + **客户端预读下一文件头几 MB**（滑动窗口）；支持 **取消正在传输的 stream**（切文件时立刻停旧流）。
- **压缩**：文本类可 zstd；已压缩格式（jpeg/png/mp4）不再二次压。
- **加密**：会话密钥派生后对 payload 加密；中继只转发密文。

---

## 7. 自适应与「码率」

本场景无实时音视频编码码率，**自适应对象**定义为：**缩略图质量、并发数、预读窗口、是否请求原图**。

建议策略：

| 信号 | 动作 |
|------|------|
| RTT 上升 / 丢包 | 降低缩略图分辨率档、减少并发 thumb 请求、缩小预读 |
| 通道判定为中继 | 默认降一档质量、限制并行、优先 WebP |
| Wi-Fi / 计费网络（客户端上报） | 默认不自动拉原图，需用户确认 |
| 滑动速度（快滑） | 取消队列中未显示的缩略图请求，只保留视口内 + 少量 lookahead |

Flutter / UniApp 侧实现 **简单 AIMD 或基于 RTT 的分档状态机** 即可，无需复杂 GCC。

---

## 8. Protobuf 建议（`proto/`）

新建包名示例：`proto/remotebrowse/browse.model.proto`、`browse.signal.proto`（具体文件名按仓库命名习惯拆分）。

建议消息（字段仅列核心）：

- `DeviceCapabilities`：`has_ipv6`, `network_metered`, `platform`
- `FileEntry`：`id`, `name`, `size`, `mtime`, `mime`, `width`, `height`, `thumb_hash`, `flags`
- `ListFilesRequest` / `ListFilesResponse`：`cursor`, `page_size`, `entries`
- `OpenFileRequest`：`id`, `variant`（thumb / preview / original）, `byte_range`
- `SignalEnvelope`：oneof { register, ice_candidate, ice_parameters, relay_token, ... }（**不含** SDP 媒体块）

生成：`server/go` 下走现有 `protogen` 流程；Flutter 使用 `protoc` 插件生成 Dart。

---

## 9. Go 与 Rust 分工与双向流

### 9.1 Go（`server/go/file` + `server/go/webrtc/session`）

- **`file`**：`/rb/health`、`/rb/v1/list`、`/rb/v1/thumb`；`rfvclient` gRPC 调 Rust 媒体；信令可反代 **`rfv-daemon`**。
- **`webrtc/session`**：`rbagent` / `rbviewer` 与直连、QUIC 房间、ICE、中继数据面；**不**再承载 Gin `/rb/signal`。

### 9.2 Rust（`server/rust/rfv`）

- **`media`**：gRPC/HTTP 列表与缩略图（`remotebrowse/`、`grpc_server`）。
- **`client`**：`webrtc-ice` + FFI（`rb_ice_viewer_*`），供 Flutter 加载 `librfv`。
- **`daemon`**：`rfv-daemon` 信令 Hub + TCP 中继 + `/rb/health`。

### 9.3 与现有代码的衔接

- `server/go/webrtc/` 内历史 **pion/webrtc 示例** 可作为「UDP 媒体」参考，**生产路径应迁移为 `pion/ice` + `quic-go`**，避免继续依赖 webrtc 主包。  
- `server/rust/rfv/src/file.rs`：目录列举、缩略图生成的既有 HTTP 处理器可逐步抽成 **库函数**，供 gRPC/UDS 任务入口复用。

---

## 10. 客户端

### 10.1 Flutter（`client/app`）

- **连接层**：数据面为 **ICE + QUIC**，浏览器无通用 API；Flutter 侧推荐 **FFI 调用与 `rfv` 同构的 Rust 库**（可后续抽 `rfv` 中 ICE/QUIC 客户端为 `cdylib`）或 **Platform Channel + 原生 Go/Rust 小模块**。Dart 纯实现 `pion/ice` 不现实；**不使用 flutter_webrtc**（该栈面向 WebRTC）。  
- **信令**：HTTPS/WSS 对接 `server/go/webrtc` 暴露的 API。  
- **UI**：`GridView` + `PageView`（与现有 `lib/pages/image/slide_image.dart` 类似）；缩略图本地磁盘缓存 + content hash。  
- **WebView**：UniApp 仅做 UI 时，**P2P 能力仍走 Flutter 原生**；独立 H5 无法完成本方案的 UDP ICE+QUIC，需明确不支持或降级为「仅中继 + 服务端转发」（非本方案首选）。

### 10.2 UniApp（`client/uniapp`）

- 相册 UI、扫码配对等；通过 **jsBridge** 调 Flutter「连接 / 列表 / 取图」。  
- 与 Flutter 能力矩阵一致：**无原生宿主则不做 P2P 数据面**。

---

## 11. 中继服务要点

- **实现位置**：与信令同属 **`server/go/webrtc`** 下子包（如 `relay/`），便于共享会话表与 TLS/QUIC 配置。
- **多租户隔离**：session id 强随机 + 短期 TTL + 绑定双方 device id。
- **限流**：每会话带宽、并发 stream、最大空闲时间。
- **可观测**：连接数、字节量、P2P 成功率、IPv6 使用率（不含用户文件内容日志）。
- **部署**：与信令同域或就近；可选边缘多实例 + **Anycast IPv6** 降低延迟。

---

## 12. 安全与合规

- 端到端加密默认开启；信令 TLS；中继 **零知识**（仅密文）。
- Agent 侧 **路径沙箱**：仅允许用户授权目录；拒绝符号链接逃逸。
- 审计：谁在何时访问了哪台设备的元数据（不存原文件）。

---

## 13. 分阶段落地

| 阶段 | 交付 |
|------|------|
| M1 | `rfv-daemon` + Go `file` 门面；`rfv` 缩略图与 protobuf 列表 |
| M2 | 直连 TCP + QUIC 房间 + ICE + 自研中继（Go CLI / Flutter / Rust FFI） |
| M3 | 自适应缩略图 + Flutter 相册 UX |
| M4 | UniApp 嵌入与扫码配对 |

---

## 14. 小结

- **目录约束**：Go **`server/go/file`**（HTTP/gRPC 门面）；Rust **`server/rust/rfv`**（媒体 + 信令 daemon + ICE FFI）；**`server/go/webrtc`** 仅保留视频/直播与 `session` CLI。  
- **连接**：手填 IP / 局域网 TCP → QUIC 房间 → ICE → 自研 TCP 中继。  
- **省带宽**：分页索引、内容寻址缩略图、Range、滑动取消、中继降档。  

*与仓库约定：`CLAUDE.md`、`proto/`、**`server/go/file`**、**`server/rust/rfv`**、`client/app`、`client/uniapp` 保持一致。*
