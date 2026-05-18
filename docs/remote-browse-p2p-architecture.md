# 远程相册浏览（NAT 内网 P2P + 中继）技术方案

本文档描述「远程相册」端到端架构：**Agent / Viewer 传输与数据面均在 `server/rust/rfv`**（`webrtc-ice`、直连 TCP、TCP 中继、ffmpeg 缩略图）；**`server/go/file` 仅 HTTP/gRPC 门面**（收请求 → 调 rfv gRPC / 反代信令 → 返回）；**不编译、不实现** `server/go/webrtc` 下的远程相册逻辑。

---

## 1. 目标与约束

| 维度 | 要求 |
|------|------|
| 客户端 | Flutter（`client/app`），可内嵌 H5（UniApp，`client/uniapp`）做部分 UI 或运营页 |
| 远端 | 同构或轻量 Agent（建议 Flutter + 原生后台 / 独立小进程均可） |
| 网络 | 双方多在 NAT 后；需打洞；失败走中继；**双端具备可用 IPv6 时优先 IPv6 直连** |
| 传输 | 打洞成功后自建应用层协议；**仅传文件元数据 + 原始/压缩字节**；缩略图与大图在客户端解码渲染 |
| 服务端 | Go `file` 负责 HTTP/gRPC 门面；Rust `rfv` 负责信令 daemon、P2P、缩略图/转码 |
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
- **数据面**：**ICE 成功后裸应用 wire 帧**（`rfv` `transport`）；失败走 **TCP 中继 `RBRL`**；之上为 protobuf 封装的列表/缩略图请求（Agent 侧 ffmpeg 生成缩略图）。

---

## 3. 连接策略（IPv6 → ICE 打洞 → 中继）

建议实现为**显式状态机**，顺序如下（可在实现中并行探测，但对外表现仍按优先级收敛）：

1. **局域网 / 手填 IP TCP 直连**（默认 `19091`）
2. **ICE 打洞**（Rust `webrtc-ice`；Flutter Viewer 经 `rb_ice_viewer_*` FFI）
3. **TCP 中继**（`rfv-daemon`；`RBRL` 握手 + 长度前缀帧）

**STUN**：`stun:stun.l.google.com:19302`（`rfv` `ice_common`）。严苛 NAT 可后续加 TURN。

---

## 4. 技术选型说明（含可替代方案）

### 4.1 为何采用 Rust `webrtc-ice` + 裸 wire

- 仅需 **NAT 打洞 + 自定义帧协议**，不必引入 WebRTC SDP/媒体栈。  
- **Agent**：主机 `rfv` + 环境变量 `RB_ROOM`（或 Flutter `rb_agent_run`）；**Viewer** 仅在客户端（Flutter `librfv`）。

### 4.2 组件

| 能力 | 实现 | 说明 |
|------|------|------|
| ICE | Rust `webrtc-ice` | `rfv` `client/ice_*` |
| 信令 + 中继 | `rfv-daemon` | `/rb/signal`、`RBRL` |
| HTTP 门面 | `server/go/file` | gRPC → rfv；`/rb/signal` 反代 daemon |
| 缩略图 | `rfv` ffmpeg | P2P wire 与 gRPC 共用 `remotebrowse/` |

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
- **主机 CLI**：`RB_ROOM=<房间> rfv`；浏览目录由 Viewer wire 指定；可选 `RB_AGENT_SANDBOX`。
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

### 9.1 Go（仅 `server/go/file`）

- **`remotebrowse` Service**：`GetHealth` / `ListFiles` / `GetThumbnail` → **`file/rfvclient` gRPC 调 rfv**（Go 不生成缩略图、不打洞）。
- **`/rb/signal`**：WebSocket **反代** `RB_SIGNAL_UPSTREAM`（默认 `ws://127.0.0.1:8080/rb/signal` → `rfv-daemon`）。

### 9.2 Rust（`server/rust/rfv`）

- **`rfv`（`host` feature）**：gRPC/HTTP；`RB_ROOM` 时并行 Agent（`wire_agent` / ICE / 中继）。
- **`client` + FFI**：`rb_ice_viewer_*`、`rb_agent_run`；Flutter `librfv`（`--features transport`）。
- **`daemon`**：`rfv-daemon` 信令 + 中继。
- **`media`**：gRPC 列表/缩略图（非 P2P 场景的 HTTP 兜底）。

### 9.3 `server/go/webrtc`

- **仅** 视频/直播 Gin 路由（`/video/*`、`/live/stream`）；**不含**远程相册代码。

---

## 10. 客户端

### 10.1 Flutter（`client/app`）

- **Viewer**：`RbViewerSession` → 直连 / **Rust ICE FFI** / 中继 → `RbWireSession`。  
- **Agent**：`RbAgentNative` → `rb_agent_run` 或子进程 `RB_ROOM=… rfv`。  
- **信令**：WSS 对接 Go `file` 的 `/rb/signal`（反代 `rfv-daemon`）。  
- **UI**：`GridView` + `PageView`（与现有 `lib/pages/image/slide_image.dart` 类似）；缩略图本地磁盘缓存 + content hash。  
- **WebView**：UniApp 仅做 UI 时，**P2P 能力仍走 Flutter 原生**；独立 H5 无法完成本方案的 UDP ICE+QUIC，需明确不支持或降级为「仅中继 + 服务端转发」（非本方案首选）。

### 10.2 UniApp（`client/uniapp`）

- 相册 UI、扫码配对等；通过 **jsBridge** 调 Flutter「连接 / 列表 / 取图」。  
- 与 Flutter 能力矩阵一致：**无原生宿主则不做 P2P 数据面**。

---

## 11. 中继服务要点

- **实现位置**：**`rfv-daemon`** `daemon/relay`（Rust）。
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
| M2 | 直连 TCP + ICE + 中继（`RB_ROOM` + rfv / Flutter viewer FFI） |
| M3 | 自适应缩略图 + Flutter 相册 UX |
| M4 | UniApp 嵌入与扫码配对 |

---

## 14. 小结

- **目录约束**：Go **`server/go/file`**（门面）；Rust **`server/rust/rfv`**（信令、P2P、媒体）；**`server/go/webrtc`** 仅视频/直播。  
- **连接**：直连 TCP → ICE → TCP 中继。  
- **省带宽**：分页索引、内容寻址缩略图、Range、滑动取消、中继降档。  

*与仓库约定：`CLAUDE.md`、`proto/`、**`server/go/file`**、**`server/rust/rfv`**、`client/app`、`client/uniapp` 保持一致。*
