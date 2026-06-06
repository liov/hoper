/// 串行化 WebSocket 上的 send/close，避免 StreamSink bound 竞态。
class RbWsOpQueue {
  Future<void> _tail = Future<void>.value();

  Future<T> run<T>(Future<T> Function() op) {
    final next = _tail.then((_) => op());
    _tail = next.then((_) {}, onError: (_) {});
    return next;
  }

  Future<void> drain() => _tail.catchError((_) {});
}

bool rbWsSinkBoundError(Object e) {
  final s = e.toString();
  return s.contains('StreamSink is bound to a stream');
}
