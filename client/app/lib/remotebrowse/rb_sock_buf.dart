import 'dart:async';
import 'dart:io';
import 'dart:typed_data';

/// Socket 读缓冲；对端断开时让 [take] 抛错，避免 readFrame 永久挂起。
class RbSockBuf {
  RbSockBuf(Socket s) {
    _sub = s.listen(
      (chunk) {
        _pending.addAll(chunk);
        _wakeWaiters();
      },
      onDone: () => _fail(StateError('连接已断开')),
      onError: (_, _) => _fail(StateError('连接已断开')),
      cancelOnError: true,
    );
  }

  final _pending = <int>[];
  final _waiters = <(int, Completer<void>)>[];
  late final StreamSubscription<Uint8List> _sub;
  Object? _failErr;

  void _wakeWaiters() {
    for (final w in List<(int, Completer<void>)>.from(_waiters)) {
      if (_pending.length >= w.$1 && !w.$2.isCompleted) {
        w.$2.complete();
      }
    }
  }

  void _fail(Object e) {
    _failErr ??= e;
    for (final w in List<(int, Completer<void>)>.from(_waiters)) {
      if (!w.$2.isCompleted) {
        w.$2.completeError(e);
      }
    }
    _waiters.clear();
  }

  Future<void> dispose() async {
    _fail(StateError('连接已关闭'));
    await _sub.cancel();
  }

  Future<Uint8List> take(int n) async {
    final err = _failErr;
    if (err != null) {
      throw err;
    }
    while (_pending.length < n) {
      final err2 = _failErr;
      if (err2 != null) {
        throw err2;
      }
      final c = Completer<void>();
      _waiters.add((n, c));
      try {
        await c.future;
      } catch (e) {
        _waiters.removeWhere((w) => w.$2 == c);
        rethrow;
      }
      _waiters.removeWhere((w) => w.$2 == c);
    }
    final out = Uint8List.fromList(_pending.sublist(0, n));
    _pending.removeRange(0, n);
    return out;
  }
}
