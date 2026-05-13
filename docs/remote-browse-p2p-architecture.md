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
│ Flutter/H5  │◄──────────────────────►│ Go 信令与编排（server/go/webrtc）│
│  client     │   注册/配对/ICE 交换     │ 会话、鉴权、配额、审计         │
└──────┬──────┘                          └──────────────┬───────────────┘
       │                                              │
       │ 可选：Go↔Rust 控制面 gRPC 双向流              │
       │                                              ▼
       │                                 ┌──────────────────────────────┐
       │                                 │ Rust 媒体/解码（server/rust/rfv）│
       │                                 │ 缩略图、按需转码、校验           │
       │                                 └──────────────────────────────┘
       │
       │  P2P（首选）或 中继（兜底）
       ▼
┌─────────────┐ ◄════ 加密应用流 ════► ┌─────────────┐
│ Viewer 端   │                        │ Agent 端    │
│ (浏览方)    │                        │ (被控方/库)  │
└─────────────┘                        └─────────────┘
```

- **信令面**：仅小流量 **ICE 参数与候选（ufrag/pwd、candidate trickle）**、会话令牌；实现代码统一落在 **`server/go/webrtc`**（目录名沿用现有结构，**不引入 WebRTC 协议栈**：不使用 SDP 音视频会话、不使用 DTLS/SCTP/DataChannel）。
- **数据面**：在 **ICE 选定后的 UDP 四元组** 上跑 **QUIC（quic-go）** 或经 **自研中继** 的 QUIC/TLS；之上承载 **文件列表、缩略图请求、原文件分片** 的 **自定义二进制协议**（建议 **TLS1.3 / Noise** 或 QUIC 自带加密，中继仅见密文）。

---

## 3. 连接策略（IPv6 → ICE 打洞 → 中继）

建议实现为**显式状态机**，顺序如下（可在实现中并行探测，但对外表现仍按优先级收敛）：

1. **IPv6 可达性探测（同网段或公网 IPv6）**  
   - 若信令交换得到双方 **全局 IPv6** 且探测成功（对目标地址发 **QUIC** 或配合 **STUN** 绑定验证），可走 **直连 QUIC（quic-go）**，失败则进入步骤 2。
2. **ICE 打洞 + QUIC（本方案固定选型）**  
   - 使用 **[pion/ice](https://github.com/pion/ice)**（如 `github.com/pion/ice/v4`）**Agent**：STUN 绑定、候选收集、连通性检查、**UDP hole punching**；**不使用** `pion/webrtc`、不使用 SDP 媒体协商、不使用 DataChannel。  
   - ICE 进入 **Connected** 且拿到可用 `ICEConn`（或等价的 `net.PacketConn` 包装）后，在该 **同一 UDP socket** 上创建 **quic-go Transport / Connection**（需处理 ICE 与 QUIC 的读写互斥、路径变更与重连策略）。  
   - 应用数据全部走 **QUIC stream** + 自定义帧 / Protobuf。
3. **自研中继**  
   - 当 ICE `failed` / QUIC 握手超时 / 策略禁止 P2P 时，双方连接 **`server/go/webrtc` 实现的中继接入点**（建议 **QUIC** 或 **TCP + TLS**），按 **session id** 转发密文；中继逻辑与信令同仓库、同目录边界，便于运维。  
   - 中继实现语言可选 Go（与信令同进程或 sidecar）；**不要求** Rust 实现中继，Rust 专注 `rfv` 内媒体流水线。

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

## 5. 信令协议（公网 Go，目录 `server/go/webrtc`）

职责：

- 设备注册、用户鉴权（与主服务集成时可 HTTP 反向代理到 Cherry，或在本包内起独立端口；**代码根目录**仍为 `server/go/webrtc`）。
- **Room / Pairing**：短码或扫码绑定 Viewer ↔ Agent。
- **交换 ICE**：`ufrag`/`pwd`、**trickle candidate**、可选 `end-of-candidates`；**不交换 SDP Offer/Answer（WebRTC 语义）**。
- 下发 **中继地址**、**会话密钥材料**（ECDH 等，中继不解密）。
- 可选：带宽档位、就近 relay 列表。

实现位置（强制约束）：

- **`server/go/webrtc/`**：例如 `api/`（gin 路由）、`ice/`（Agent 封装）、`relay/`（中继接入）、`service/`（会话状态机）、沿用并扩展 `service/signal/`（WS hub）。
- Proto：`proto/remotebrowse/`（或 `proto/webrtc/ice_signal.proto` 等）定义信令消息；`protogen` 生成 Go；Flutter 用同套 proto。

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

### 9.1 Go（`server/go/webrtc`）

- **ICE 信令、会话状态机、中继接入认证**、配额与审计（与主站打通时通过内部 API / gRPC 调用 `server/go` 其他包亦可，但 **远程浏览相关新增代码默认落在本目录**）。
- 对 Rust（`rfv`）：同机部署时 **gRPC 双向流** 或 **Unix domain stream**，由 `webrtc` 包内客户端向 **`rfv` 暴露的处理接口** 下发「路径 / 任务 id / 字节范围」，回传缩略图或错误；Go 不做重 CPU 解码。

### 9.2 Rust（强制 `server/rust/rfv`）

- **不新增** `server/rust` 下并列 crate 承载本特性；在 **`rfv`** 内扩展模块即可，例如：  
  - `rfv/src/remotebrowse/` 或 `rfv/src/thumbnail_worker.rs`：列表缓存、缩略图/预览编码、哈希。  
  - 与现有 `file.rs` 复用常量（如 `MAX_SIZE`）、FFmpeg 初始化逻辑，避免重复进程级开销。
- 与 Go 的 **双向流**：便于背压与取消（用户快速滑动时取消缩略图任务）。

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
| M1 | `server/go/webrtc` 公网信令 + 同局域网 **QUIC** 调试；`rfv` 缩略图与 protobuf 列表 |
| M2 | **pion/ice** 打洞成功后在同一 UDP 上 **quic-go** P2P；失败切 **`webrtc` 包内自研中继** |
| M3 | IPv6 直连分支 + 自适应缩略图 + Flutter 相册 UX |
| M4 | **`webrtc`↔`rfv` 双向流**接缩略图流水线；UniApp 嵌入与扫码配对 |

---

## 14. 小结

- **目录约束**：Go 侧 **`server/go/webrtc`**（信令 + ICE 编排 + 中继接入）；**不用 WebRTC**；数据面为 **`pion/ice` + `quic-go`**。  
- **Rust 约束**：媒体与解码逻辑全部在 **`server/rust/rfv`** 内扩展，经双向流与 `webrtc` 包协作。  
- **连接**：IPv6 可直达时 **QUIC**；否则 **ICE 打洞 + QUIC**；再失败 **自研中继**。  
- **省带宽**：分页索引、内容寻址缩略图、Range、滑动取消、中继降档。  

若 ICE 与 quic-go 绑定 UDP 的实现成本超预期，可在 **不改变业务目录命名** 的前提下，在文档附录中单独立项「技术债 / 备选：TURN-only UDP 通道简化 QUIC 部署」或评估 **libp2p（需目录与仓库策略另议）**。

---

*文档路径：`hoper/docs/remote-browse-p2p-architecture.md`*  
*与仓库约定：`CLAUDE.md`、`proto/`、**`server/go/webrtc`**（本特性 Go 代码）、**`server/rust/rfv`**（本特性 Rust 代码）、`client/app`、`client/uniapp` 保持一致。*
