# ProtoBuf数据契约

<cite>
**本文档引用的文件**
- [proto/README.md](file://proto/README.md)
- [common.model.proto](file://proto/common/common.model.proto)
- [common.service.proto](file://proto/common/common.service.proto)
- [user.model.proto](file://proto/user/user.model.proto)
- [user.service.proto](file://proto/user/user.service.proto)
- [content.model.proto](file://proto/content/content.model.proto)
- [content.service.proto](file://proto/content/content.service.proto)
- [moment.model.proto](file://proto/content/moment.model.proto)
- [note.model.proto](file://proto/content/note.model.proto)
- [diary.model.proto](file://proto/content/diary.model.proto)
- [action.model.proto](file://proto/content/action.model.proto)
- [message.proto](file://proto/message/message.proto)
- [file.service.proto](file://proto/file/file.service.proto)
</cite>

## 更新摘要
**变更内容**
- 用户模型重构：支持国际电话号码格式，countryCallingCode和phone字段分离设计
- 认证请求/响应结构重组：LoginReq、SignupReq等消息结构优化
- 验证代码处理优化：增加更多字段验证规则和数据类型约束
- 设备信息增强：AccessDevice和Device模型完善地理位置和设备信息

## 目录
1. [简介](#简介)
2. [项目结构](#项目结构)
3. [核心组件](#核心组件)
4. [架构概览](#架构概览)
5. [详细组件分析](#详细组件分析)
6. [依赖关系分析](#依赖关系分析)
7. [性能考虑](#性能考虑)
8. [故障排除指南](#故障排除指南)
9. [结论](#结论)

## 简介

Hoper项目的ProtoBuf数据契约是一套完整的跨平台通信协议定义，采用Protocol Buffers 3.0语法构建。该系统通过精心设计的消息结构和枚举类型，为用户管理、内容管理、文件服务和消息通信提供了标准化的数据交换格式。

**更新** 本次重大更新重构了用户模型以支持国际电话号码格式，优化了认证流程的请求/响应结构，并增强了字段验证规则和数据类型约束。

本项目的核心设计理念包括：
- **向后兼容性保证**：通过字段标签管理和版本控制确保新旧版本间的无缝兼容
- **类型安全**：利用ProtoBuf的强类型系统防止运行时类型错误
- **跨平台支持**：生成Go、TypeScript等多种语言的客户端代码
- **验证集成**：内置字段验证规则和业务逻辑约束
- **国际化支持**：专门的国际电话号码处理机制

## 项目结构

项目采用按功能域划分的目录结构，每个模块都有独立的ProtoBuf定义文件：

```mermaid
graph TB
subgraph "ProtoBuf定义目录"
A[proto/] --> B[common/]
A --> C[user/]
A --> D[content/]
A --> E[file/]
A --> F[message/]
B --> B1[common.model.proto]
B --> B2[common.service.proto]
C --> C1[user.model.proto]
C --> C2[user.service.proto]
D --> D1[content.model.proto]
D --> D2[content.service.proto]
D --> D3[moment.model.proto]
D --> D4[note.model.proto]
D --> D5[diary.model.proto]
D --> D6[action.model.proto]
E --> E1[file.service.proto]
F --> F1[message.proto]
end
subgraph "生成代码目录"
G[proto/.generated/]
G --> G1[go/]
G --> G2[ts/]
end
A --> G
```

**图表来源**
- [proto/README.md](file://proto/README.md)
- [common.model.proto](file://proto/common/common.model.proto)
- [user.model.proto](file://proto/user/user.model.proto)
- [content.model.proto](file://proto/content/content.model.proto)

**章节来源**
- [proto/README.md](file://proto/README.md)
- [common.model.proto](file://proto/common/common.model.proto)

## 核心组件

### 数据模型层

系统定义了四类核心数据模型：用户模型、内容模型、文件模型和消息模型。

#### 用户模型设计原则

**更新** 用户模型经过重大重构，现在支持国际电话号码格式：

```mermaid
classDiagram
class User {
+uint64 id
+string name
+string mail
+string countryCallingCode
+string phone
+string account
+string password
+Gender gender
+time.Date birthday
+string countryCode
+string address
+string intro
+string signature
+string avatar
+string cover
+Role role
+string realName
+string idNo
+timestamp.Timestamp activatedAt
+model.ModelTime modelTime
+timestamp.Timestamp bannedAt
+UserStatus status
}
class UserExt {
+uint64 id
+uint64 score
+uint64 follow
+uint64 followed
+timestamp.Timestamp lastActivatedAt
}
class UserBase {
+uint64 id
+string name
+uint64 score
+Gender gender
+string avatar
}
User --> UserExt : "扩展信息"
User --> UserBase : "基础信息"
```

**图表来源**
- [user.model.proto](file://proto/user/user.model.proto)

#### 内容模型架构

内容模型采用多态设计，支持多种内容类型：

```mermaid
classDiagram
class Content {
+uint64 id
+ContentType type
+string title
+string content
+uint64 userId
+ViewPermission permission
+uint32 status
+model.ModelTime modelTime
}
class Moment {
+uint64 id
+string content
+repeated string images
+common.DataType type
+common.TinyTag mood
+common.Area area
+uint64 userId
+ViewPermission permission
}
class Note {
+uint64 id
+string title
+string content
+uint32 status
+model.ModelTime modelTime
}
class Diary {
+uint64 id
+string content
+uint64 bookId
+DiaryBook book
+common.TinyTag mood
+common.Area area
+uint64 userId
+ViewPermission permission
}
Content <|-- Moment : "继承"
Content <|-- Note : "继承"
Content <|-- Diary : "继承"
```

**图表来源**
- [content.model.proto](file://proto/content/content.model.proto)
- [moment.model.proto](file://proto/content/moment.model.proto)
- [note.model.proto](file://proto/content/note.model.proto)
- [diary.model.proto](file://proto/content/diary.model.proto)

**章节来源**
- [user.model.proto](file://proto/user/user.model.proto)
- [content.model.proto](file://proto/content/content.model.proto)

### 服务接口层

服务接口采用gRPC和HTTP/JSON双栈设计，提供RESTful API和高性能RPC两种访问方式。

#### 服务设计模式

```mermaid
sequenceDiagram
participant Client as "客户端"
participant Gateway as "API网关"
participant Service as "业务服务"
participant DB as "数据库"
Client->>Gateway : HTTP/JSON请求
Gateway->>Service : gRPC调用
Service->>DB : 数据查询/更新
DB-->>Service : 查询结果
Service-->>Gateway : 业务响应
Gateway-->>Client : HTTP/JSON响应
Note over Client,Service : 同时支持gRPC直连
```

**图表来源**
- [common.service.proto](file://proto/common/common.service.proto)
- [user.service.proto](file://proto/user/user.service.proto)
- [content.service.proto](file://proto/content/content.service.proto)

**章节来源**
- [common.service.proto](file://proto/common/common.service.proto)
- [user.service.proto](file://proto/user/user.service.proto)
- [content.service.proto](file://proto/content/content.service.proto)

## 架构概览

### 数据流架构

```mermaid
graph TB
subgraph "客户端层"
A[Web前端]
B[移动端Flutter]
C[桌面应用]
D[第三方集成]
end
subgraph "网关层"
E[HTTP/JSON网关]
F[gRPC网关]
G[认证授权]
end
subgraph "服务层"
H[用户服务]
I[内容服务]
J[文件服务]
K[消息服务]
end
subgraph "数据层"
L[PostgreSQL]
M[Redis缓存]
N[对象存储]
end
A --> E
B --> E
C --> E
D --> E
E --> G
G --> H
G --> I
G --> J
G --> K
H --> L
I --> L
J --> N
K --> M
L --> H
L --> I
N --> J
M --> K
```

**图表来源**
- [user.service.proto](file://proto/user/user.service.proto)
- [content.service.proto](file://proto/content/content.service.proto)
- [file.service.proto](file://proto/file/file.service.proto)
- [message.proto](file://proto/message/message.proto)

### 字段标签和版本管理

ProtoBuf通过字段标签实现向后兼容性：

```mermaid
flowchart TD
A[新增字段] --> B{是否向后兼容?}
B --> |是| C[分配新标签]
B --> |否| D[创建新版本]
C --> E[保持现有字段标签不变]
D --> F[版本号递增]
E --> G[客户端可安全解析]
F --> H[向后兼容处理]
G --> I[数据完整性]
H --> I
```

**图表来源**
- [common.model.proto](file://proto/common/common.model.proto)
- [user.model.proto](file://proto/user/user.model.proto)

## 详细组件分析

### 用户管理数据契约

#### 用户状态枚举设计

用户状态通过专用枚举类型实现严格的业务状态控制：

```mermaid
classDiagram
class UserStatus {
<<enumeration>>
UserStatusPlaceholder
UserStatusInActive
UserStatusActivated
UserStatusFrozen
UserStatusDeleted
}
class Role {
<<enumeration>>
PlaceholderRole
RoleNormal
RoleAdmin
RoleSuperAdmin
}
class Gender {
<<enumeration>>
GenderPlaceholder
GenderUnfilled
GenderMale
GenderFemale
}
class DeviceType {
<<enumeration>>
DeviceTypePlaceholder
DeviceTypePhone
DeviceTypePC
DeviceTypePad
DeviceTypeWatch
}
UserStatus --> Role : "关联"
UserStatus --> Gender : "关联"
UserStatus --> DeviceType : "关联"
```

**图表来源**
- [user.model.proto](file://proto/user/user.model.proto)

#### 用户操作日志追踪

系统提供完整的用户行为追踪机制：

```mermaid
sequenceDiagram
participant U as "用户"
participant S as "系统"
participant L as "日志存储"
U->>S : 执行操作
S->>S : 记录操作详情
S->>L : 写入ActionLog
L-->>S : 确认存储
S-->>U : 返回操作结果
Note over U,L : 完整的操作审计轨迹
```

**图表来源**
- [user.model.proto](file://proto/user/user.model.proto)

**章节来源**
- [user.model.proto](file://proto/user/user.model.proto)
- [user.service.proto](file://proto/user/user.service.proto)

### 内容管理数据契约

#### 内容类型体系

内容系统支持多种内容类型，每种类型都有特定的字段组合：

```mermaid
classDiagram
class ContentType {
<<enumeration>>
ContentMoment
ContentNote
ContentDairy
ContentDairyBook
ContentFavorites
ContentCollection
ContentComment
}
class ViewPermission {
<<enumeration>>
ViewPermissionAll
ViewPermissionSelf
ViewPermissionHomePage
ViewPermissionStranger
ViewPermissionShield
ViewPermissionOpen
}
class ContainerType {
<<enumeration>>
ContainerTypeFavorites
ContainerTypeDiaryBook
ContainerTypeAlbum
ContainerTypeCollection
}
ContentType --> ViewPermission : "权限控制"
ContainerType --> ViewPermission : "权限控制"
```

**图表来源**
- [content.model.proto](file://proto/content/content.model.proto)

#### 内容统计和分析

系统提供丰富的内容统计指标：

```mermaid
erDiagram
CONTENT {
uint64 id PK
ContentType type
string title
string content
uint64 userId
uint32 status
}
STATISTICS {
uint64 id PK
ContentType type
uint64 refId FK
uint64 like
uint64 browse
uint64 comment
uint64 collect
uint64 share
}
USER_ACTION {
uint64 likeId
uint64 unlikeId
uint64[] collectIds
}
CONTENT ||--o{ STATISTICS : "统计"
CONTENT ||--|| USER_ACTION : "用户动作"
```

**图表来源**
- [action.model.proto](file://proto/content/action.model.proto)

**章节来源**
- [content.model.proto](file://proto/content/content.model.proto)
- [action.model.proto](file://proto/content/action.model.proto)

### 文件服务数据契约

#### 文件上传流程

文件服务提供完整的文件生命周期管理：

```mermaid
flowchart TD
A[文件上传请求] --> B{预上传检查}
B --> |存在| C[直接返回URL]
B --> |不存在| D[分片上传]
D --> E[获取上传凭证]
E --> F[分片上传]
F --> G[完成合并]
C --> H[文件可用]
G --> H
H --> I[文件访问]
I --> J[文件删除]
```

**图表来源**
- [file.service.proto](file://proto/file/file.service.proto)

**章节来源**
- [file.service.proto](file://proto/file/file.service.proto)

### 消息通信数据契约

#### 实时消息架构

消息系统支持多种消息类型和传输模式：

```mermaid
classDiagram
class MQMessage {
+uint64 id
+ServerCmd command
+timestamp.Timestamp created_at
+uint64 user_id
+uint64 recv_user_id
+uint64 recv_group_id
+bytes payload
+uint32 type
}
class ClientMessage {
+uint64 id
+ClientCmd command
+timestamp.Timestamp created_at
+uint64 recv_user_id
+uint64 recv_group_id
+bytes payload
+uint32 type
+timestamp.Timestamp read_at
}
class ServerMessage {
+uint64 id
+ServerCmd command
+timestamp.Timestamp time
+bytes payload
+uint32 type
}
class MessageType {
<<enumeration>>
TEXT
BINARY
IMAGE
FILE
VIDEO
AUDIO
}
MQMessage --> MessageType : "消息类型"
ClientMessage --> MessageType : "消息类型"
ServerMessage --> MessageType : "消息类型"
```

**图表来源**
- [message.proto](file://proto/message/message.proto)

**章节来源**
- [message.proto](file://proto/message/message.proto)

## 依赖关系分析

### ProtoBuf导入依赖

```mermaid
graph TB
subgraph "核心依赖"
A[hopeio/utils/enum/enum.proto]
B[hopeio/time/time.proto]
C[hopeio/model/basic.proto]
D[hopeio/time/deletedAt/deletedAt.proto]
end
subgraph "项目特定依赖"
E[hopeio/utils/patch/go.proto]
F[hopeio/time/timestamp/timestamp.proto]
G[protoc-gen-openapiv2/options/annotations.proto]
end
subgraph "业务模块"
H[user/user.model.proto]
I[content/content.model.proto]
J[common/common.model.proto]
K[file/file.service.proto]
L[message/message.proto]
end
A --> H
A --> I
A --> J
B --> H
B --> I
C --> H
C --> I
D --> I
E --> H
E --> I
F --> H
F --> I
G --> H
G --> I
G --> J
G --> K
G --> L
```

**图表来源**
- [user.model.proto](file://proto/user/user.model.proto)
- [content.model.proto](file://proto/content/content.model.proto)
- [common.model.proto](file://proto/common/common.model.proto)

### 版本兼容性策略

系统采用渐进式版本升级策略：

```mermaid
flowchart LR
A[ProtoBuf 1.0] --> B[ProtoBuf 1.1]
B --> C[ProtoBuf 2.0]
A --> D[向后兼容]
B --> D
C --> D
D --> E[字段标签不重复使用]
D --> F[新增字段使用新标签]
D --> G[保留字段标记]
```

**图表来源**
- [proto/README.md](file://proto/README.md)

**章节来源**
- [proto/README.md](file://proto/README.md)

## 性能考虑

### 序列化性能优化

ProtoBuf相比JSON具有显著的性能优势：

- **二进制序列化**：比JSON更紧凑，网络传输更高效
- **零拷贝支持**：Go语言中支持零拷贝字节切片
- **流式处理**：支持大数据量的流式序列化

### 缓存策略

```mermaid
graph TB
A[客户端缓存] --> B[内存缓存]
A --> C[本地存储]
D[服务器缓存] --> E[Redis缓存]
D --> F[数据库查询缓存]
B --> G[热点数据]
C --> H[离线数据]
E --> I[会话数据]
F --> J[查询结果缓存]
```

## 故障排除指南

### 常见问题诊断

#### 字段验证失败

当字段验证失败时，系统会返回明确的错误信息：

```mermaid
flowchart TD
A[请求到达] --> B{字段验证}
B --> |通过| C[业务处理]
B --> |失败| D[验证错误]
D --> E[检查字段类型]
E --> F[检查字段范围]
F --> G[检查必填字段]
C --> H[成功响应]
G --> I[错误响应]
```

#### 兼容性问题

版本升级可能导致的兼容性问题：

1. **字段标签冲突**：避免重新使用已删除字段的标签
2. **数据类型变更**：确保新旧版本间的数据类型兼容
3. **枚举值扩展**：新增枚举值不影响现有客户端

**章节来源**
- [user.service.proto](file://proto/user/user.service.proto)
- [content.service.proto](file://proto/content/content.service.proto)

## 结论

Hoper项目的ProtoBuf数据契约设计体现了现代微服务架构的最佳实践。通过精心设计的数据模型、完善的枚举体系和严格的服务接口定义，系统实现了：

- **强类型安全**：通过ProtoBuf的静态类型检查防止运行时错误
- **向后兼容**：通过字段标签管理和版本控制确保系统演进的平滑性
- **跨平台支持**：生成多种语言的客户端代码，支持广泛的开发场景
- **性能优化**：采用二进制序列化和缓存策略，确保高并发场景下的性能表现
- **国际化支持**：专门的国际电话号码处理机制，支持全球化业务需求

**更新** 本次重大更新重构了用户模型以支持国际电话号码格式，优化了认证流程的请求/响应结构，并增强了字段验证规则和数据类型约束。这套数据契约不仅满足了当前业务需求，更为未来的功能扩展和技术演进奠定了坚实的基础。通过遵循ProtoBuf的设计原则和最佳实践，开发者可以快速理解和使用这些数据契约，提高开发效率和系统可靠性。