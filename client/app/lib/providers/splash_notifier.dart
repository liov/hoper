import 'dart:async';

import 'package:app/global/state.dart';
import 'package:app/pages/route.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/scheduler.dart';
import 'package:flutter/widgets.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'splash_notifier.g.dart';

class SplashState {
  const SplashState({this.countdown = 0});

  final int countdown;

  SplashState copyWith({int? countdown}) => SplashState(countdown: countdown ?? this.countdown);
}

@Riverpod(keepAlive: true)
class Splash extends _$Splash {
  Timer? _timer;
  Duration time = const Duration(seconds: 3);
  var _navigated = false;
  var _disposed = false;
  Completer<void>? _adWait;
  DateTime? pausedTime;

  @override
  SplashState build() {
    ref.onDispose(() {
      _disposed = true;
      _timer?.cancel();
    });
    Future.microtask(_prepare);
    return const SplashState();
  }

  void skip() {
    _timer?.cancel();
    final wait = _adWait;
    if (wait != null && !wait.isCompleted) {
      wait.complete();
    }
    _adWait = null;
    if (_navigated) return;
    _navigated = true;
    SchedulerBinding.instance.addPostFrameCallback((_) => AppNavigator.offAllNamed(Routes.HOME));
  }

  void advertising(Widget splash) {
    if (pausedTime == null) return;
    if (DateTime.now().difference(pausedTime!) < const Duration(minutes: 10)) return;
    globalService.logger.fine('advertising');
    AppNavigator.showOverlay(
      loadingWidget: splash,
      asyncFunction: () {
        _timer?.cancel();
        _navigated = false;
        _adWait = Completer<void>();
        _startCountdown();
        return _adWait!.future;
      },
    );
  }

  Future<void> _prepare() async {
    state = state.copyWith(countdown: 0);
    try {
      if (!globalState.initialized) {
        await globalState.init();
      }
    } catch (e, st) {
      globalService.logger.warning('splash init failed: $e\n$st');
    }
    if (!_disposed) {
      _startCountdown();
    }
  }

  void _startCountdown() {
    _timer?.cancel();
    state = state.copyWith(countdown: (time.inMilliseconds / 1000).round().clamp(1, 99));
    _timer = Timer.periodic(const Duration(seconds: 1), (t) {
      if (_disposed) {
        t.cancel();
        return;
      }
      final next = state.countdown - 1;
      state = state.copyWith(countdown: next);
      if (next <= 0) {
        t.cancel();
        skip();
      }
    });
  }
}
