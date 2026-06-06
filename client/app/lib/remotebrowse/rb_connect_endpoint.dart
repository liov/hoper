import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/rb_connection_profile.dart';

/// 连接目标：信令服务器（ws/wss）或 TCP 直连（IP/域名，含公网、IPv6）。
enum RbConnectTargetKind { signal, direct }

class RbConnectTarget {
  RbConnectTarget._(this.kind, {this.signalUri, this.directHost, this.directPort = RbDirectDialer.defaultPort});

  final RbConnectTargetKind kind;
  final Uri? signalUri;
  final String? directHost;
  final int directPort;

  bool get isDirect => kind == RbConnectTargetKind.direct;
  bool get isSignal => kind == RbConnectTargetKind.signal;

  /// 存 profile.direct：默认端口可省略。
  String get directStorage {
    if (!isDirect || directHost == null) {
      return '';
    }
    if (directPort == RbDirectDialer.defaultPort) {
      return _formatDirectHost(directHost!);
    }
    return '${_formatDirectHost(directHost!)}:$directPort';
  }

  static String endpointDisplay(RbConnectionProfile p) {
    final d = p.direct.trim();
    if (d.isNotEmpty) {
      return d;
    }
    final s = p.signalUrl.trim();
    if (s.isNotEmpty) {
      return formatHostPortInput(s);
    }
    return '';
  }

  /// 存 profile 用：仅 host:port，不含 ws 路径。
  static String endpointStorage(RbConnectTarget target) {
    if (target.isDirect) {
      return target.directStorage;
    }
    final uri = target.signalUri;
    if (uri == null) {
      return '';
    }
    return formatSignalHostPort(uri);
  }

  static String defaultEndpointInput() => formatSignalHostPort(parseSignalWsUri(RemoteBrowseApi.rbDebugBaseUrl));

  static String formatHostPortInput(String raw) {
    final t = raw.trim();
    if (t.contains('://') || t.contains('/rb/')) {
      return formatSignalHostPort(parseSignalWsUri(t));
    }
    return t;
  }

  static String formatSignalHostPort(Uri uri) {
    final h = uri.host;
    final host = h.contains(':') ? '[$h]' : h;
    if (!uri.hasPort) {
      return host;
    }
    final implicit = uri.scheme == 'wss' ? 443 : 80;
    if (uri.port == implicit) {
      return host;
    }
    return '$host:${uri.port}';
  }

  /// 优先探测 `/rb/health`；成功→信令，失败→直连（缺省 19091）。
  static Future<RbConnectTarget> resolve(String raw, {Duration probeTimeout = const Duration(seconds: 3)}) async {
    final t = raw.trim();
    if (t.isEmpty) {
      throw ArgumentError('地址不能为空');
    }
    if (t.contains('://') || t.contains('/rb/')) {
      return RbConnectTarget._(RbConnectTargetKind.signal, signalUri: parseSignalWsUri(t));
    }
    final hp = parseHostPort(t);
    final signalUri = parseSignalWsUri(t);
    try {
      await probeSignalHealth(signalUri, timeout: probeTimeout);
      return RbConnectTarget._(RbConnectTargetKind.signal, signalUri: signalUri);
    } catch (_) {
      return RbConnectTarget._(
        RbConnectTargetKind.direct,
        directHost: hp.host,
        directPort: hp.port ?? RbDirectDialer.defaultPort,
      );
    }
  }

  static RbConnectTarget fromProfile(RbConnectionProfile p) {
    final d = p.direct.trim();
    if (d.isNotEmpty) {
      return parseDirect(d);
    }
    final s = p.signalUrl.trim();
    if (s.isNotEmpty) {
      return RbConnectTarget._(RbConnectTargetKind.signal, signalUri: parseSignalWsUri(s));
    }
    throw ArgumentError('连接未配置服务器地址');
  }

  /// 同步解析（仅用于已存 profile 的直连字段；用户输入请用 [resolve]）。
  static RbConnectTarget parse(String raw) => parseDirect(raw);

  static RbConnectTarget parseDirect(String raw) {
    final t = raw.trim();
    if (t.isEmpty) {
      throw ArgumentError('地址不能为空');
    }
    final hp = parseHostPort(t);
    return RbConnectTarget._(
      RbConnectTargetKind.direct,
      directHost: hp.host,
      directPort: hp.port ?? RbDirectDialer.defaultPort,
    );
  }

}

/// 解析 host[:port]，支持 IPv6 `[addr]`、`[addr]:port`。
({String host, int? port}) parseHostPort(String raw) {
  var s = raw.trim();
  if (s.isEmpty) {
    throw ArgumentError('地址不能为空');
  }
  if (s.startsWith('[')) {
    final end = s.indexOf(']');
    if (end < 1) {
      throw ArgumentError('IPv6 地址需使用方括号，例如 [2001:db8::1]:19091');
    }
    final host = s.substring(1, end);
    if (s.length == end + 1) {
      return (host: host, port: null);
    }
    if (s[end + 1] != ':') {
      throw ArgumentError('IPv6 端口格式应为 [地址]:端口');
    }
    final port = int.tryParse(s.substring(end + 2));
    if (port == null || port < 1 || port > 65535) {
      throw ArgumentError('端口无效');
    }
    return (host: host, port: port);
  }
  final colon = s.lastIndexOf(':');
  if (colon > 0 && !s.contains(']')) {
    final portStr = s.substring(colon + 1);
    final port = int.tryParse(portStr);
    if (port != null && portStr == port.toString()) {
      return (host: s.substring(0, colon), port: port);
    }
  }
  return (host: s, port: null);
}

String _formatDirectHost(String host) {
  if (host.contains(':') && !host.startsWith('[')) {
    return '[$host]';
  }
  return host;
}

/// `Socket.connect` 用：去掉 IPv6 方括号。
String rbSocketConnectHost(String host) {
  final t = host.trim();
  if (t.startsWith('[') && t.endsWith(']')) {
    return t.substring(1, t.length - 1);
  }
  return t;
}
