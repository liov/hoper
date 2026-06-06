/// Viewer 与 Agent 之间的 P2P 数据通道类型。
enum RbLinkKind {
  /// IP/域名 TCP 直连（默认 19091，含公网、IPv6）
  directTcp,

  /// WebRTC ICE
  ice,

  /// 经 rb-daemon 中继
  relayTcp,
}

extension RbLinkKindDisplay on RbLinkKind {
  String get label => switch (this) {
    RbLinkKind.directTcp => '直连',
    RbLinkKind.ice => '自动组网',
    RbLinkKind.relayTcp => '服务器转发',
  };

  String get detail => switch (this) {
    RbLinkKind.directTcp => 'IP/域名 TCP（默认 19091）',
    RbLinkKind.ice => 'ICE + gRPC（h2）',
    RbLinkKind.relayTcp => '经服务器转发',
  };
}
