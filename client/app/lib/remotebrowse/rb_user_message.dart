import 'package:app/remotebrowse/rb_grpc_session.dart';

/// 预览页 [Get.back]：文件连接断开，由列表页自动重连。
const rbPreviewResultWireLost = 'wire_lost';

bool rbIsReadSuperseded(Object e) => e is StateError && e.message == rbReadSuperseded;

/// ExoPlayer / loopback 代理常见播放失败文案（不等于数据面已断，须结合 [RbGrpcSession.isOpen]）。
bool rbIsPlaybackDisconnect(String message) {
  if (message.trim().isEmpty) {
    return false;
  }
  if (rbIsWireLost(StateError(message))) {
    return true;
  }
  final s = message.toLowerCase();
  return s.contains('source error') ||
      s.contains('playbackexception') ||
      s.contains('unexpected end of stream') ||
      s.contains('failed to connect') ||
      s.contains('connection timed out') ||
      s.contains('internal server error') ||
      s.contains('service unavailable') ||
      s.contains('stream closed') ||
      s.contains('connectexception');
}

bool rbIsWireLost(Object e) {
  final s = e.toString();
  if (s.contains('Channel shutting down') ||
      s.contains('P2P gRPC 连接已关闭') ||
      s.contains('Socket has been closed') ||
      s.contains('socket has been closed')) {
    return true;
  }
  return s.contains('连接已断开') ||
      s.contains('连接已关闭') ||
      s.contains('Connection closed') ||
      s.contains('Connection reset') ||
      s.contains('Broken pipe') ||
      s.contains('ice read failed') ||
      s.contains('ice closed') ||
      s.contains('ice grpc failed') ||
      s.contains('unexpected wire type') ||
      s.contains('删除超时') ||
      s.contains('agent wire version mismatch') ||
      s.contains('StreamSink is bound to a stream') ||
      s.contains('UNAVAILABLE') ||
      s.contains('PROTOCOL_ERROR') ||
      s.contains('unexpected eof') ||
      s.contains('Unexpected EOF');
}

/// 播放失败时是否应走「整路重连」（仅数据面真断或 gRPC/h2 级错误）。
bool rbPlaybackShouldReconnect(String message, {required bool wireOpen}) {
  if (!wireOpen) {
    return true;
  }
  return rbIsWireLost(StateError(message));
}

/// 将异常转成用户可读文案（不依赖 dart:io，Web/原生通用）。
String rbUserMessage(Object error) {
  final s = error.toString();
  if (s.contains('RbConnectCancelled') || s.contains('已取消')) {
    return '已取消连接';
  }
  if (s.contains('Operation timed out') || s.contains('TimeoutException') || s.contains('timed out')) {
    if (s.contains('删除超时')) {
      return '删除超时，请确认对方设备上的共享已开启并重试';
    }
    return '连接超时，请检查网络，并确认对方已在本机「共享」页开启相册';
  }
  if (s.contains('unexpected wire type') || s.contains('agent wire version mismatch')) {
    return '对方版本过旧，请更新对方设备上的 App 或 Agent 后重试';
  }
  if (s.contains('Connection refused')) {
    return '无法连上服务器，请检查网络与服务器地址';
  }
  if (s.contains('Failed host lookup') || s.contains('No address associated')) {
    return '无法解析服务器地址，请检查连接设置';
  }
  if (s.contains('Connection reset') || s.contains('Connection closed')) {
    return '连接已断开，正在尝试恢复…';
  }
  if (s.contains('WebSocketChannelException')) {
    final inner = _unwrapMessage(s);
    if (inner != s) {
      return rbUserMessage(inner);
    }
    return '无法连上信令服务，请检查网络与服务器地址';
  }
  if (s.contains('channel-error') || s.contains('Unable to establish connection on channel')) {
    return '预览组件未就绪，请完全退出 App 后重试';
  }
  if (s.contains('StreamSink is bound to a stream')) {
    return '文件连接已断开';
  }
  if (s.contains('signal closed')) {
    return '文件连接已断开';
  }
  if (s.contains('Invalid port')) {
    return '服务器端口无效，请检查连接设置';
  }
  if (error is ArgumentError) {
    return error.message?.toString() ?? s;
  }
  if (s.contains('信令注册超时')) {
    return '对方未响应，请确认对方已开启「共享」且房间码一致';
  }
  if (s.contains('信令无注册响应') || s.contains('信令连接超时')) {
    return '无法连上对方，请确认对方已开启共享且房间码一致';
  }
  if (s.contains('无法访问') && s.contains('/rb/health')) {
    final m = s.contains('Bad state:') ? s.replaceFirst('Bad state: ', '') : s;
    return m.contains('health') ? '远程服务未就绪，请稍后再试' : m;
  }
  if (s.contains('path outside RB_AGENT_SANDBOX') || s.contains('path outside')) {
    return '路径超出 Agent 限制：请清空 Agent「路径限制」并重启共享；MSYS 终端执行 unset RB_AGENT_SANDBOX 后重启 rb';
  }
  if (s.contains('Windows 路径请使用')) {
    return s.contains('Bad state:') ? s.replaceFirst('Bad state: ', '') : s;
  }
  if (s.contains('HTTP 400')) {
    final m = RegExp(r'HTTP 400:\s*(.*)').firstMatch(s);
    final detail = m?.group(1)?.trim() ?? '';
    if (detail.isNotEmpty) {
      return detail.length > 120 ? '${detail.substring(0, 120)}…' : detail;
    }
    return '请求被拒绝 (400)，请确认 PC Agent 已更新并重启';
  }
  if (error is StateError) {
    final m = error.message;
    if (m.isNotEmpty && !m.startsWith('Bad state')) {
      return m;
    }
  }
  return s.length > 200 ? '${s.substring(0, 200)}…' : s;
}

String _unwrapMessage(String s) {
  final m = RegExp(r'WebSocketChannelException[^(]*\(([^)]+)\)').firstMatch(s);
  if (m != null) {
    return m.group(1)?.trim() ?? s;
  }
  return s;
}
