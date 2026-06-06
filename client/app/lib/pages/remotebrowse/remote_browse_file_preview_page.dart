import 'dart:async';
import 'dart:convert';
import 'dart:io' show Platform;

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/direct_dialer.dart';
import 'package:app/remotebrowse/log.dart';
import 'package:app/remotebrowse/rb_zoomable_image.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_memory_picture.dart';
import 'package:app/remotebrowse/rb_motion_photo.dart';
import 'package:app/remotebrowse/rb_pdf_preview.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_preview_edge_swipe.dart';
import 'package:app/remotebrowse/rb_preview_nav.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:app/remotebrowse/rb_preview_store.dart';
import 'package:app/remotebrowse/rb_remote_video_player.dart';
import 'package:app/remotebrowse/rb_transcode.dart';
import 'package:app/remotebrowse/rb_video_timeline.dart';
import 'package:app/remotebrowse/rb_video_stream_server.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:app/pages/remotebrowse/rb_file_info_panel.dart';
import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/rb_file_download.dart';
import 'package:app/pages/remotebrowse/rb_preview_folder_strip.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_markdown_plus/flutter_markdown_plus.dart';
import 'package:app/util/nav.dart';
import 'package:video_player/video_player.dart';
import 'package:webview_flutter/webview_flutter.dart';

class RemoteBrowseFilePreviewPage extends StatefulWidget {
  const RemoteBrowseFilePreviewPage({
    super.key,
    required this.nav,
    this.onListChanged,
    this.wireResolver,
    this.thumbsResolver,
    this.previewStore,
    this.wireEpoch,
    this.onWireLost,
  });

  final RbPreviewNav nav;
  final VoidCallback? onListChanged;
  /// 列表页重连后返回最新 [RbGrpcSession]，切文件前会刷新。
  final RbGrpcSession? Function()? wireResolver;
  final RbThumbLoader? Function()? thumbsResolver;
  final RbPreviewContentStore? previewStore;
  /// 列表页数据面重连成功时递增，预览页自动重试。
  final ValueListenable<int>? wireEpoch;
  final VoidCallback? onWireLost;

  @override
  State<RemoteBrowseFilePreviewPage> createState() => _RemoteBrowseFilePreviewPageState();
}

class _RemoteBrowseFilePreviewPageState extends State<RemoteBrowseFilePreviewPage> {
  late final RbPreviewNav _nav;
  var _loading = true;
  var _deleting = false;
  var _downloading = false;
  var _listChanged = false;
  var _error = '';
  RbPreviewKind _kind = RbPreviewKind.unsupported;
  Uint8List? _bytes;
  String _text = '';
  VideoPlayerController? _videoCtrl;
  RbVideoStreamServer? _videoStream;
  VoidCallback? _videoErrListener;
  VideoPlayerController? _videoErrCtrl;
  String _videoLoadLabel = '';
  var _videoTranscodePreset = RbTranscodePreset.source.id;
  var _videoTranscodeVcodec = RbTranscodeVcodec.hevc;
  var _videoAgentHttpAvailable = false;
  var _videoPresetSwitching = false;
  var _videoTranscodeSeeking = false;
  var _videoPendingSeekMs = 0;
  /// HLS 裁剪列表时播放器 position 为相对时间，需加此偏移；连续流为全片 PTS 时常为 0。
  var _videoHlsOffsetMs = 0;
  /// 全片时长（毫秒），来自目录 `durationMs`，用于转码切档后修正进度条总长。
  var _videoTimelineTotalMs = 0;
  var _videoTranscodeUseHls = false;
  Uri? _transcodeHlsIndexUri;
  var _lastHlsWarmChunk = -1;
  var _videoStreamAnchorMs = 0;
  DateTime? _videoStreamPlayAt;
  int? _videoStreamPausedAbsMs;
  var _videoWasPlayingFlag = false;
  VoidCallback? _videoTimelineListener;
  var _videoPlaybackRetry = 0;
  var _motionPlaying = false;
  var _loadEpoch = 0;
  var _showImmersiveChrome = false;
  var _fileSwitching = false;
  var _wireLostHandling = false;
  var _disposing = false;
  var _lastWireEpoch = 0;
  var _scrollAtTop = true;
  var _scrollAtBottom = true;
  var _renderMarkdown = true;
  var _renderHtml = true;
  var _imageShowingOriginal = true;
  var _imageCanShowOriginalButton = false;
  var _imageLoadingOriginal = false;
  String _imagePreviewMime = '';
  Uint8List? _imageThumbBytes;
  var _previewFetchActive = false;
  var _previewLoadReceived = 0;
  var _previewLoadTotal = 0;
  DateTime? _previewProgressThrottle;
  ScrollController? _textScroll;
  Uri? _pdfMediaUri;
  RbVideoStreamServer? _pdfStream;
  RbPreviewContentStore? _ownedPreviewStore;

  RbPreviewContentStore? get _previewStore => widget.previewStore ?? _ownedPreviewStore;

  bool get _docScrollKind =>
      _kind == RbPreviewKind.text || (_kind == RbPreviewKind.markdown && !_renderMarkdown) || (_kind == RbPreviewKind.html && !_renderHtml);

  bool get _docWebKind => _kind == RbPreviewKind.pdf || (_kind == RbPreviewKind.html && _renderHtml);

  bool get _docImmersiveKind =>
      _kind == RbPreviewKind.text ||
      _kind == RbPreviewKind.markdown ||
      _kind == RbPreviewKind.html ||
      _kind == RbPreviewKind.pdf;

  bool get _imagePreviewKind => _kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto;

  bool get _hasImageThumbPlaceholder => _imageThumbBytes != null && _imageThumbBytes!.isNotEmpty;

  bool get _canShowImageWhileLoading =>
      _imagePreviewKind && (_hasImageThumbPlaceholder || (_bytes != null && _bytes!.isNotEmpty));

  bool get _showPreviewLoadBar => _previewFetchActive && (_loading || _imageLoadingOriginal);

  double? get _previewLoadProgressValue {
    final total = _previewLoadTotal;
    if (total <= 0) {
      return null;
    }
    return (_previewLoadReceived / total).clamp(0.0, 1.0);
  }

  @override
  void initState() {
    super.initState();
    _nav = widget.nav;
    _lastWireEpoch = widget.wireEpoch?.value ?? 0;
    widget.wireEpoch?.addListener(_onWireEpoch);
    _nav.thumbs?.addListener(_onThumbLoaderChanged);
    _textScroll = ScrollController()..addListener(_syncTextScrollEdges);
    if (widget.previewStore == null && !kIsWeb) {
      _ownedPreviewStore = RbPreviewContentStore(hostKey: widget.nav.wire.thumbCacheHostKey());
      unawaited(_ownedPreviewStore!.purgeExpired());
    }
    _loadCurrent();
  }

  @override
  void dispose() {
    _disposing = true;
    widget.wireEpoch?.removeListener(_onWireEpoch);
    _nav.thumbs?.removeListener(_onThumbLoaderChanged);
    _loadEpoch++;
    final scroll = _textScroll;
    if (scroll != null) {
      scroll.removeListener(_syncTextScrollEdges);
      scroll.dispose();
    }
    _textScroll = null;
    final vc = _videoCtrl;
    final vs = _videoStream;
    _videoCtrl = null;
    _videoStream = null;
    _clearVideoErrListener(vc);
    unawaited(_releaseVideo(vc, vs));
    unawaited(_disposePdf());
    super.dispose();
  }

  void _resetScrollEdges() {
    _scrollAtTop = true;
    _scrollAtBottom = true;
    if (!mounted) {
      return;
    }
    final c = _textScroll;
    if (c != null && c.hasClients) {
      try {
        c.jumpTo(0);
      } catch (_) {}
    }
  }

  void _syncTextScrollEdges() {
    if (!mounted) {
      return;
    }
    final c = _textScroll;
    if (c == null || !c.hasClients) {
      return;
    }
    final m = c.position;
    final top = m.pixels <= m.minScrollExtent + 0.5;
    final bottom = m.pixels >= m.maxScrollExtent - 0.5;
    if (top == _scrollAtTop && bottom == _scrollAtBottom) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }
      setState(() {
        _scrollAtTop = top;
        _scrollAtBottom = bottom;
      });
    });
  }

  bool _onScrollNotification(ScrollNotification n) {
    if (n.metrics.axis != Axis.vertical) {
      return false;
    }
    final top = n.metrics.pixels <= n.metrics.minScrollExtent + 0.5;
    final bottom = n.metrics.pixels >= n.metrics.maxScrollExtent - 0.5;
    if (top == _scrollAtTop && bottom == _scrollAtBottom) {
      return false;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) {
        return;
      }
      setState(() {
        _scrollAtTop = top;
        _scrollAtBottom = bottom;
      });
    });
    return false;
  }

  bool _canEdgeVerticalSwipe(bool next) {
    if (_docWebKind) {
      return false;
    }
    if (!_docScrollKind) {
      return true;
    }
    return next ? _scrollAtBottom : _scrollAtTop;
  }

  /// 先从 Widget 树摘掉 VideoPlayer，再释放原生播放器。[blocking] 仅退出预览时用，切文件须 false 以免 dispose 卡死 ANR。
  Future<void> _disposeVideo({bool blocking = false}) async {
    final c = _videoCtrl;
    final stream = _videoStream;
    _videoCtrl = null;
    _videoStream = null;
    _videoHlsOffsetMs = 0;
    _transcodeHlsIndexUri = null;
    _lastHlsWarmChunk = -1;
    _videoTimelineTotalMs = 0;
    _videoStreamAnchorMs = 0;
    _videoStreamPlayAt = null;
    _videoStreamPausedAbsMs = null;
    if (c == null && stream == null) {
      return;
    }
    if (mounted && !_disposing) {
      setState(() {});
      await WidgetsBinding.instance.endOfFrame;
    }
    if (!mounted || _disposing) {
      await _releaseVideo(c, stream);
      return;
    }
    final job = _releaseVideo(c, stream);
    if (blocking) {
      try {
        await job.timeout(const Duration(seconds: 2));
      } catch (_) {}
    } else {
      unawaited(job);
    }
  }

  Future<void> _disposePdf() async {
    final stream = _pdfStream;
    _pdfStream = null;
    _pdfMediaUri = null;
    if (stream != null) {
      try {
        await stream.close();
      } catch (_) {}
    }
  }

  void _clearVideoTimelineListener(VideoPlayerController? c) {
    final listener = _videoTimelineListener;
    if (listener != null && c != null) {
      try {
        c.removeListener(listener);
      } catch (_) {}
    }
    _videoTimelineListener = null;
  }

  void _clearVideoErrListener(VideoPlayerController? c) {
    final listener = _videoErrListener;
    if (c != null && listener != null && _videoErrCtrl == c) {
      try {
        c.removeListener(listener);
      } catch (_) {}
    }
    _videoErrListener = null;
    _videoErrCtrl = null;
  }

  /// 先从树摘掉 [VideoPlayer]，再 dispose，避免切换码率时平台侧仍持有已释放的 player id。
  Future<void> _detachVideoFromTree() async {
    final c = _videoCtrl;
    final stream = _videoStream;
    _clearVideoErrListener(c);
    _clearVideoTimelineListener(c);
    _videoCtrl = null;
    _videoStream = null;
    if (c != null) {
      try {
        if (c.value.isInitialized) {
          await c.pause();
        }
      } catch (_) {}
    }
    if (mounted && !_disposing) {
      setState(() {});
      await WidgetsBinding.instance.endOfFrame;
      await WidgetsBinding.instance.endOfFrame;
    }
    if (!kIsWeb && (Platform.isAndroid || Platform.isIOS)) {
      await Future<void>.delayed(const Duration(milliseconds: 280));
    }
    await _releaseVideo(c, stream);
  }

  Duration _absoluteVideoPosition() {
    final c = _videoCtrl;
    if (c == null) {
      return Duration.zero;
    }
    try {
      final v = c.value;
      if (!v.isInitialized) {
        return Duration.zero;
      }
      return Duration(
        milliseconds: rbVideoAbsPosMs(
          playerPosMs: v.position.inMilliseconds,
          timelineOffsetMs: _videoHlsOffsetMs,
          playbackAnchorMs: _videoStreamAnchorMs,
          isPlaying: v.isPlaying,
          playStartedAt: _videoStreamPlayAt,
          pausedAbsPosMs: _videoStreamPausedAbsMs,
          playbackSpeed: v.playbackSpeed,
        ),
      );
    } catch (_) {
      return Duration.zero;
    }
  }

  void _onVideoTimelineTick() {
    final c = _videoCtrl;
    if (c == null || !mounted) {
      return;
    }
    try {
      final v = c.value;
      if (!v.isInitialized) {
        return;
      }
      if ((_usesAgentTranscode && !_videoTranscodeUseHls) || _videoHlsOffsetMs > 0) {
        final playing = v.isPlaying;
        if (playing && !_videoWasPlayingFlag) {
          _videoStreamPlayAt = DateTime.now();
          _videoStreamPausedAbsMs = null;
        } else if (!playing && _videoWasPlayingFlag) {
          _videoStreamPausedAbsMs = rbVideoAbsPosMs(
            playerPosMs: v.position.inMilliseconds,
            timelineOffsetMs: _videoHlsOffsetMs,
            playbackAnchorMs: _videoStreamAnchorMs,
            isPlaying: true,
            playStartedAt: _videoStreamPlayAt,
            pausedAbsPosMs: null,
            playbackSpeed: v.playbackSpeed,
          );
          _videoStreamPlayAt = null;
        }
        _videoWasPlayingFlag = playing;
      }
      if (v.isPlaying) {
        _warmHlsFragmentsIfNeeded();
      }
    } catch (_) {
      return;
    }
    _safeSetState(() {});
  }

  void _warmHlsFragmentsIfNeeded() {
    final uri = _transcodeHlsIndexUri;
    if (uri == null || !_videoTranscodeUseHls || !_usesAgentTranscode) {
      return;
    }
    final absMs = _absoluteVideoPosition().inMilliseconds;
    final chunk = absMs ~/ rbTranscodeFragmentMs;
    if (chunk <= _lastHlsWarmChunk) {
      return;
    }
    _lastHlsWarmChunk = chunk;
    final fromMs = (chunk + 3) * rbTranscodeFragmentMs;
    rbWarmTranscodeAhead(uri, fromMs, fragmentCount: 8);
  }

  void _syncTranscodeTimeline(VideoPlayerController player, {required int beginMs}) {
    if (_videoTranscodePreset == RbTranscodePreset.source.id) {
      _videoHlsOffsetMs = 0;
      return;
    }
    if (_videoTranscodeUseHls) {
      _videoHlsOffsetMs = beginMs;
      return;
    }
    // 连续流 MPEG-TS 为全片绝对 PTS；initialize 时 position 常为 0，不能据此加 beginMs。
    _videoHlsOffsetMs = 0;
  }

  bool _videoWasPlaying() {
    final c = _videoCtrl;
    if (c == null) {
      return true;
    }
    try {
      return c.value.isPlaying;
    } catch (_) {
      return true;
    }
  }

  bool get _usesAgentTranscode => _videoTranscodePreset != RbTranscodePreset.source.id && _videoAgentHttpAvailable;

  void _setVideoLoadLabel(String label) {
    if (!mounted || _videoLoadLabel == label) {
      return;
    }
    setState(() => _videoLoadLabel = label);
  }

  void _attachVideoListeners(VideoPlayerController player, int epoch) {
    _clearVideoErrListener(_videoErrCtrl);
    _clearVideoTimelineListener(_videoErrCtrl);
    void onErr() => _onVideoCtrlError(player, epoch);
    _videoErrListener = onErr;
    _videoErrCtrl = player;
    player.addListener(onErr);
    void onTimeline() => _onVideoTimelineTick();
    _videoTimelineListener = onTimeline;
    player.addListener(onTimeline);
    _videoWasPlayingFlag = player.value.isPlaying;
  }

  Future<void> _seekTranscodeHls(Duration target, {bool playAfter = false}) async {
    if (!_usesAgentTranscode || _videoCtrl == null || _videoTranscodeSeeking || _videoPresetSwitching) {
      return;
    }
    _videoTranscodeSeeking = true;
    _videoPendingSeekMs = target.inMilliseconds;
    try {
      await _swapTranscodePlayerAt(
        resumeAt: target,
        playing: playAfter || _videoWasPlaying(),
        hint: '正在跳转…',
        timeoutMsg: '跳转超时，请稍后重试',
        presetSwitch: false,
      );
    } catch (e, st) {
      rbLog.warning('transcode seek failed', e, st);
      if (mounted) {
        _notice('跳转', rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 2));
      }
    } finally {
      _videoTranscodeSeeking = false;
    }
  }

  Future<VideoPlayerController> _initTranscodePlayerAt({
    required int epoch,
    required int beginMs,
    required Uri? streamUri,
    required Uri hlsUri,
    required bool presetSwitch,
    required String timeoutMsg,
  }) async {
    var player = await _prepareTranscodePlayer(
      epoch: epoch,
      streamUri: streamUri,
      hlsUri: hlsUri,
      transcodeReload: true,
      presetSwitch: presetSwitch,
      transcodeSeek: !presetSwitch,
      transcodeResume: true,
      beginMs: beginMs,
    );
    try {
      await player.initialize().timeout(
        Duration(seconds: presetSwitch ? 120 : 90),
        onTimeout: () => throw StateError(timeoutMsg),
      );
    } catch (e) {
      if (streamUri != null && rbPreferTranscodeStream) {
        rbLog.warning('transcode stream init failed, fallback hls: $e');
        await player.dispose();
        player = await _prepareTranscodePlayer(
          epoch: epoch,
          streamUri: streamUri,
          hlsUri: hlsUri,
          transcodeReload: true,
          presetSwitch: presetSwitch,
          transcodeSeek: !presetSwitch,
          transcodeResume: true,
          beginMs: beginMs,
          forceHls: true,
        );
        await player.initialize().timeout(
          Duration(seconds: presetSwitch ? 120 : 90),
          onTimeout: () => throw StateError(timeoutMsg),
        );
      } else {
        rethrow;
      }
    }
    return player;
  }

  Future<void> _swapTranscodePlayerAt({
    required Duration resumeAt,
    required bool playing,
    required String hint,
    required String timeoutMsg,
    required bool presetSwitch,
  }) async {
    if (!_usesAgentTranscode || _videoCtrl == null) {
      return;
    }
    final epoch = _loadEpoch;
    final beginMs = (resumeAt.inMilliseconds ~/ rbTranscodeFragmentMs) * rbTranscodeFragmentMs;
    final oldCtrl = _videoCtrl;
    final oldStream = _videoStream;
    _clearVideoErrListener(oldCtrl);
    _clearVideoTimelineListener(oldCtrl);
    _videoStreamAnchorMs = beginMs;
    _videoStreamPausedAbsMs = resumeAt.inMilliseconds;
    _videoStreamPlayAt = null;
    _setVideoLoadLabel(hint);
    try {
    if (oldCtrl != null) {
      try {
        if (oldCtrl.value.isInitialized) {
          await oldCtrl.pause();
        }
      } catch (_) {}
    }
    final slot = _nav.current;
    final streamUri = _nav.wire.agentTranscodeStreamUri(
      slot.relPath,
      presetId: _videoTranscodePreset,
      vcodec: _videoTranscodeVcodec,
      startMs: beginMs,
    );
    final hlsUri = _nav.wire.agentTranscodeHlsUri(
      slot.relPath,
      presetId: _videoTranscodePreset,
      vcodec: _videoTranscodeVcodec,
      beginMs: beginMs,
    );
    if (hlsUri == null) {
      throw StateError('转码地址无效');
    }
    if (!mounted || epoch != _loadEpoch) {
      return;
    }
    final player = await _initTranscodePlayerAt(
      epoch: epoch,
      beginMs: beginMs,
      streamUri: streamUri,
      hlsUri: hlsUri,
      presetSwitch: presetSwitch,
      timeoutMsg: timeoutMsg,
    );
    if (!mounted || epoch != _loadEpoch) {
      await player.dispose();
      return;
    }
    _syncTranscodeTimeline(player, beginMs: beginMs);
    _videoStreamAnchorMs = _videoTranscodeUseHls ? 0 : beginMs;
    _videoStream = null;
    _videoCtrl = player;
    _videoPlaybackRetry = 0;
    _nav.thumbs?.setPaused(true);
    _attachVideoListeners(player, epoch);
    _videoStreamPausedAbsMs = null;
    _setVideoLoadLabel('');
    if (playing) {
      try {
        await player.play();
        if (!_videoTranscodeUseHls) {
          _videoStreamPlayAt = DateTime.now();
          _videoWasPlayingFlag = true;
        }
      } catch (_) {}
    }
    unawaited(_releaseVideo(oldCtrl, oldStream));
    } finally {
      if (mounted && _videoLoadLabel == hint) {
        _setVideoLoadLabel('');
      }
    }
  }

  Future<VideoPlayerController> _prepareTranscodePlayer({
    required int epoch,
    required Uri? streamUri,
    required Uri hlsUri,
    required bool transcodeReload,
    required bool presetSwitch,
    required bool transcodeSeek,
    required bool transcodeResume,
    required int beginMs,
    bool forceHls = false,
  }) async {
    if (!forceHls && rbPreferTranscodeStream && streamUri != null) {
      final playUri = transcodeReload ? rbTranscodeUriCacheBust(streamUri) : streamUri;
      if (transcodeReload) {
        final warmTimeout = presetSwitch ? const Duration(seconds: 180) : const Duration(seconds: 90);
        await rbWarmTranscodeStream(playUri, timeout: warmTimeout);
        if (epoch != _loadEpoch || !mounted) {
          throw StateError('cancelled');
        }
      }
      rbLog.info('video play transcode stream $playUri');
      _videoTranscodeUseHls = false;
      _transcodeHlsIndexUri = null;
      _lastHlsWarmChunk = -1;
      return VideoPlayerController.networkUrl(playUri);
    }
    _videoTranscodeUseHls = true;
    _transcodeHlsIndexUri = hlsUri;
    _lastHlsWarmChunk = -1;
    final playUri = transcodeReload ? rbTranscodeUriCacheBust(hlsUri) : hlsUri;
    rbLog.info('video play transcode hls $playUri');
    if (transcodeReload) {
      final warmFrags = presetSwitch ? 12 : (transcodeSeek ? 8 : 10);
      final warmTimeout = presetSwitch ? const Duration(seconds: 180) : const Duration(seconds: 90);
      final ok = await rbWarmTranscodePlaylist(
        playUri,
        timeout: warmTimeout,
        fragmentCount: warmFrags,
        retries: 2,
        warmStartMs: beginMs,
      );
      if (epoch != _loadEpoch || !mounted) {
        throw StateError('cancelled');
      }
      if (!ok) {
        throw StateError(transcodeResume ? '缓冲失败，请稍后重试' : '转码流准备失败，分片转码较慢请稍后重试');
      }
    } else {
      final warm = await rbWarmTranscodePlaylist(playUri, timeout: const Duration(seconds: 90), fragmentCount: 6, retries: 2);
      if (epoch != _loadEpoch || !mounted) {
        throw StateError('cancelled');
      }
      if (!warm) {
        throw StateError('转码流准备失败，请稍后重试');
      }
    }
    return VideoPlayerController.networkUrl(playUri, formatHint: VideoFormat.hls);
  }

  Future<void> _resumeVideoPlayback(VideoPlayerController player, Duration at, {required bool play}) async {
    if (!player.value.isInitialized) {
      return;
    }
    var target = at;
    if (target > Duration.zero) {
      for (var i = 0; i < 24 && player.value.duration <= Duration.zero; i++) {
        await Future<void>.delayed(const Duration(milliseconds: 50));
        if (!player.value.isInitialized) {
          return;
        }
      }
      final dur = player.value.duration;
      if (dur > Duration.zero && target > dur) {
        target = dur;
      }
      try {
        await player.seekTo(target);
      } catch (e) {
        rbLog.warning('video resume seek: $e');
      }
    }
    if (play) {
      try {
        await player.play();
      } catch (_) {}
    }
  }

  Future<void> _releaseVideo(VideoPlayerController? c, RbVideoStreamServer? stream) async {
    _nav.thumbs?.setPaused(false);
    if (stream != null) {
      try {
        await stream.close().timeout(const Duration(seconds: 2));
      } catch (_) {}
    }
    if (c != null) {
      _clearVideoErrListener(c);
      _clearVideoTimelineListener(c);
      try {
        await c.dispose().timeout(const Duration(seconds: 2));
      } catch (_) {}
    }
  }

  void _safeSetState(VoidCallback fn) {
    if (!mounted) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        setState(fn);
      }
    });
  }

  void _syncWireFromParent() {
    final w = widget.wireResolver?.call();
    if (w != null && w.isOpen) {
      _nav.wire = w;
    }
    final t = widget.thumbsResolver?.call();
    if (t != null && t != _nav.thumbs) {
      _nav.thumbs?.removeListener(_onThumbLoaderChanged);
      _nav.thumbs = t;
      t.addListener(_onThumbLoaderChanged);
    }
  }

  void _onThumbLoaderChanged() {
    if (!mounted || _disposing || !_imagePreviewKind) {
      return;
    }
    if (_bytes != null && _bytes!.isNotEmpty) {
      return;
    }
    final slot = _nav.current;
    final b = _nav.thumbs?.bytesFor(slot.relPath);
    if (b == null || b.isEmpty || identical(b, _imageThumbBytes)) {
      return;
    }
    setState(() {
      _imageThumbBytes = b;
      if (_loading && (_bytes == null || _bytes!.isEmpty)) {
        _loading = false;
      }
    });
  }

  void _resetPreviewLoadProgress() {
    _previewLoadReceived = 0;
    _previewLoadTotal = 0;
    _previewProgressThrottle = null;
  }

  void _beginPreviewFetch({int? expectedTotal}) {
    _previewFetchActive = true;
    _resetPreviewLoadProgress();
    if (expectedTotal != null && expectedTotal > 0) {
      _previewLoadTotal = expectedTotal;
    }
  }

  void _endPreviewFetch() {
    _previewFetchActive = false;
    _resetPreviewLoadProgress();
  }

  void _onPreviewReadProgress(int received, int total) {
    if (!mounted) {
      return;
    }
    final now = DateTime.now();
    final throttle = _previewProgressThrottle;
    if (throttle != null && now.difference(throttle) < const Duration(milliseconds: 80)) {
      _previewLoadReceived = received;
      if (total > 0) {
        _previewLoadTotal = total;
      }
      return;
    }
    _previewProgressThrottle = now;
    setState(() {
      _previewLoadReceived = received;
      if (total > 0) {
        _previewLoadTotal = total;
      }
    });
  }

  void _primeImageThumb(RbPreviewFileSlot slot) {
    final loader = _nav.thumbs;
    if (loader == null) {
      _imageThumbBytes = null;
      return;
    }
    var b = loader.bytesFor(slot.relPath);
    if (b == null || b.isEmpty) {
      loader.requestThumb(slot.relPath, slot.entry, bypassPause: true);
      b = loader.bytesFor(slot.relPath);
    }
    _imageThumbBytes = b;
  }

  void _onWireEpoch() {
    if (!mounted || _disposing) {
      return;
    }
    final v = widget.wireEpoch?.value ?? 0;
    if (v == _lastWireEpoch) {
      return;
    }
    _lastWireEpoch = v;
    _syncWireFromParent();
    if (!_nav.wire.isOpen) {
      return;
    }
    if (_error.contains('重新连接') || _error.contains('连接已断开') || _error.contains('连接已关闭')) {
      unawaited(_loadCurrent());
    }
  }

  bool _requireOpenWire() {
    _syncWireFromParent();
    if (_nav.wire.isOpen) {
      return true;
    }
    if (!_wireLostHandling) {
      unawaited(_handleWireLost(StateError('连接已关闭')));
    }
    return false;
  }

  /// 先停掉视频/PDF 再通知列表重连，避免 gRPC 已关时 ExoPlayer 仍拉流导致 MediaCodec 异常与闪退。
  Future<void> _handleWireLost(Object e) async {
    if (_wireLostHandling) {
      return;
    }
    _wireLostHandling = true;
    try {
      await _disposeVideo(blocking: false);
      await _disposePdf();
      widget.onWireLost?.call();
      if (!mounted) {
        return;
      }
      _safeSetState(() {
        _error = '${rbUserMessage(e)}，正在重新连接…';
        _loading = false;
      });
    } finally {
      _wireLostHandling = false;
    }
  }

  void _onVideoCtrlError(VideoPlayerController ctrl, int epoch) {
    if (epoch != _loadEpoch ||
        !mounted ||
        ctrl != _videoCtrl ||
        _disposing ||
        _fileSwitching ||
        _videoPresetSwitching ||
        _videoTranscodeSeeking) {
      return;
    }
    final err = ctrl.value.errorDescription;
    if (err == null || err.isEmpty) {
      return;
    }
    _syncWireFromParent();
    if (rbPlaybackShouldReconnect(err, wireOpen: _nav.wire.isOpen)) {
      unawaited(_handleWireLost(StateError(err)));
      return;
    }
    final isTranscode = _videoTranscodePreset != RbTranscodePreset.source.id;
    if (isTranscode && _videoPlaybackRetry < 1 && (err.contains('Source error') || err.contains('ExoPlaybackException'))) {
      _videoPlaybackRetry++;
      var resumeAt = _absoluteVideoPosition();
      if (resumeAt <= Duration.zero && _videoPendingSeekMs > 0) {
        resumeAt = Duration(milliseconds: _videoPendingSeekMs);
      }
      final resumePlaying = _videoWasPlaying();
      rbLog.info('video playback retry transcode preset=$_videoTranscodePreset at=${resumeAt.inMilliseconds}ms');
      unawaited(() async {
        if (!mounted || epoch != _loadEpoch) {
          return;
        }
        _setVideoLoadLabel(_videoTranscodeUseHls ? '缓冲转码分片，继续播放…' : '正在恢复转码流…');
        try {
          await _detachVideoFromTree();
          if (!mounted || epoch != _loadEpoch) {
            return;
          }
          await _loadVideo(
            _nav.current,
            epoch,
            transcodeResume: true,
            presetResumeAt: resumeAt,
            presetResumePlaying: resumePlaying,
          );
        } catch (e) {
          rbLog.warning('video playback retry failed: $e');
        }
      }());
      return;
    }
    rbLog.fine('video playback error (no reconnect): $err');
    unawaited(() async {
      await _disposeVideo(blocking: false);
      if (!mounted || epoch != _loadEpoch) {
        return;
      }
      _safeSetState(() => _error = _videoPlaybackErrorMessage(err));
    }());
  }

  String _videoPlaybackErrorMessage(String err) {
    if (err.contains('Source error') || err.contains('ExoPlaybackException')) {
      return '视频播放失败，请重试';
    }
    return rbUserMessage(StateError(err));
  }

  Future<void> _loadCurrent() async {
    if (!mounted || _disposing) {
      return;
    }
    try {
      await _loadCurrentBody();
    } catch (e, st) {
      rbLog.warning('preview load failed', e, st);
      if (rbIsReadSuperseded(e)) {
        return;
      }
      if (rbIsWireLost(e)) {
        _syncWireFromParent();
        if (!_nav.wire.isOpen) {
          unawaited(_handleWireLost(e));
        } else {
          _safeSetState(() {
            _loading = false;
            _error = rbUserMessage(e);
          });
        }
        return;
      }
      _safeSetState(() {
        _loading = false;
        _error = rbUserMessage(e);
      });
    }
  }

  Future<void> _loadCurrentBody() async {
    if (!_requireOpenWire()) {
      return;
    }
    var handledWireLost = false;
    final (focused, focusMsg) = await _nav.focusPreviewableFile();
    if (!focused) {
      if (mounted) {
        setState(() {
          _loading = false;
          _error = focusMsg ?? '目录中没有可预览文件';
        });
      }
      return;
    }
    final slot = _nav.current;
    final prevKind = _kind;
    final nextKind = rbPreviewKindForEntry(slot.entry);
    final epoch = ++_loadEpoch;
    final hadVideo = prevKind == RbPreviewKind.video || _videoCtrl != null || _videoStream != null;
    if (mounted && !_disposing && (hadVideo || nextKind == RbPreviewKind.video)) {
      setState(() {
        _loading = true;
        if (prevKind == RbPreviewKind.video && nextKind == RbPreviewKind.video) {
          _videoLoadLabel = '正在切换…';
        } else if (prevKind != RbPreviewKind.video && nextKind == RbPreviewKind.video) {
          _videoLoadLabel = '正在准备播放…';
        } else {
          _videoLoadLabel = '';
        }
      });
      await WidgetsBinding.instance.endOfFrame;
    }
    _syncWireFromParent();
    if (!_nav.wire.isOpen) {
      unawaited(_handleWireLost(StateError('连接已关闭')));
      return;
    }
    await _nav.wire.cancelPendingReadsAndWait();
    _syncWireFromParent();
    if (!_nav.wire.isOpen) {
      unawaited(_handleWireLost(StateError('连接已关闭')));
      return;
    }
    if (hadVideo) {
      await _disposeVideo(blocking: false);
      if (prevKind == RbPreviewKind.video && nextKind != RbPreviewKind.video) {
        _nav.thumbs?.setPaused(false);
      }
    }
    await _disposePdf();
    if (epoch != _loadEpoch || !mounted || _disposing) {
      return;
    }
    _motionPlaying = false;
    _kind = rbPreviewKindForEntry(slot.entry);
    _resetScrollEdges();
    if (!mounted || _disposing) {
      return;
    }
    setState(() {
      _showImmersiveChrome = false;
      // 视频/PDF 自管 loading，勿占全页 _loading
      _loading = _kind != RbPreviewKind.video && _kind != RbPreviewKind.pdf;
      _error = '';
      _bytes = null;
      _text = '';
      _videoLoadLabel = '';
      _pdfMediaUri = null;
      _imageShowingOriginal = true;
      _imageCanShowOriginalButton = false;
      _imageLoadingOriginal = false;
      _imagePreviewMime = '';
      _imageThumbBytes = null;
      _endPreviewFetch();
    });
    try {
      if (_kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto) {
        _primeImageThumb(slot);
        await _loadImage(slot, epoch);
      } else if (_kind == RbPreviewKind.video) {
        await _loadVideo(slot, epoch);
      } else if (_kind == RbPreviewKind.pdf) {
        await _loadPdf(slot, epoch);
      } else if (_kind == RbPreviewKind.unsupported) {
        _error = '目录中没有可预览文件';
      } else {
        final max = rbPreviewMaxBytes(_kind);
        final bytes = await _readPreviewBytes(slot, maxBytes: max, epoch: epoch);
        if (epoch != _loadEpoch || !mounted) {
          return;
        }
        _bytes = bytes;
        if (_kind == RbPreviewKind.markdown || _kind == RbPreviewKind.text || _kind == RbPreviewKind.html) {
          _text = utf8.decode(bytes, allowMalformed: true);
        }
      }
    } catch (e) {
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (rbIsReadSuperseded(e)) {
        return;
      }
      if (rbIsWireLost(e)) {
        _syncWireFromParent();
        if (!_nav.wire.isOpen) {
          handledWireLost = true;
          unawaited(_handleWireLost(e));
        } else {
          _error = rbUserMessage(e);
        }
        return;
      }
      _error = rbUserMessage(e);
    } finally {
      if (!handledWireLost && mounted && epoch == _loadEpoch && !_disposing) {
        _endPreviewFetch();
        _safeSetState(() => _loading = false);
        if (_kind == RbPreviewKind.text) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!mounted || epoch != _loadEpoch || _disposing) {
              return;
            }
            _syncTextScrollEdges();
          });
        }
        final loader = _nav.thumbs;
        final prefetchThumbs = loader != null &&
            _nav.wire.isOpen &&
            _error.isEmpty &&
            _kind != RbPreviewKind.video &&
            !_docScrollKind &&
            !_docWebKind;
        if (prefetchThumbs) {
          WidgetsBinding.instance.addPostFrameCallback((_) {
            if (!mounted || epoch != _loadEpoch || _disposing || !_nav.wire.isOpen) {
              return;
            }
            loader.prefetchVisibleRange(_nav.filesAround(_nav.fileIndex, before: 2, after: 4), maxEdge: RbThumbLoader.sharedThumbEdge);
          });
        }
      }
    }
  }

  Future<Uint8List> _readPreviewBytes(
    RbPreviewFileSlot slot, {
    required int maxBytes,
    required int epoch,
    bool imagePreview = false,
  }) async {
    final store = _previewStore;
    if (store != null) {
      final hit = await store.load(slot.entry, slot.relPath, maxBytes: maxBytes);
      if (hit != null && hit.isNotEmpty) {
        if (!imagePreview || !rbIsLikelyStaleImagePreviewCache(slot.entry, hit.length)) {
          return hit;
        }
      }
    }
    _beginPreviewFetch(expectedTotal: slot.entry.size > 0 ? slot.entry.size : maxBytes);
    if (mounted) {
      setState(() {});
    }
    try {
      final r = await _nav.wire.readFileFull(
        slot.relPath,
        maxBytes: maxBytes,
        onProgress: _onPreviewReadProgress,
      );
      if (epoch != _loadEpoch) {
        throw StateError(rbReadSuperseded);
      }
      if (store != null && r.bytes.isNotEmpty) {
        unawaited(store.put(slot.entry, slot.relPath, r.bytes));
      }
      return r.bytes;
    } finally {
      if (epoch == _loadEpoch) {
        _endPreviewFetch();
      }
    }
  }

  Future<void> _loadPdf(RbPreviewFileSlot slot, int epoch) async {
    final store = _previewStore;
    final cached = await store?.mediaFileIfComplete(slot.entry, slot.relPath, logicalTotal: slot.entry.size);
    Uri? uri = cached != null ? Uri.file(cached.path) : _nav.wire.agentMediaPlayUri(slot.relPath);
    RbVideoStreamServer? stream;
    try {
      if (uri == null) {
        if (slot.entry.size > rbPreviewPdfMaxBytes) {
          _error = 'PDF 超过 ${rbPreviewPdfMaxBytes >> 20}MB 且无法经 HTTP 直连 Agent';
          return;
        }
        stream = await RbVideoStreamServer.start(
          wire: _nav.wire,
          relPath: slot.relPath,
          totalSize: slot.entry.size,
          fileName: slot.entry.name,
          mime: 'application/pdf',
          entry: slot.entry,
          previewStore: store,
        );
        if (epoch != _loadEpoch || !mounted) {
          return;
        }
        uri = stream.playUri;
        _pdfStream = stream;
      }
      _pdfMediaUri = uri;
    } catch (e) {
      await stream?.close();
      if (epoch == _loadEpoch && mounted) {
        rethrow;
      }
    }
  }

  Future<void> _onVideoTranscodePreset(String presetId) async {
    if (presetId == _videoTranscodePreset || _kind != RbPreviewKind.video) {
      return;
    }
    if (presetId != RbTranscodePreset.source.id && !_videoAgentHttpAvailable) {
      _notice('清晰度', '当前连接无法转码，请使用原画或直连 Agent', tone: RbBannerTone.warning);
      return;
    }
    if (presetId == RbTranscodePreset.source.id && _videoAgentHttpAvailable) {
      final meta = await _nav.wire.fetchAgentMediaMeta(_nav.current.relPath);
      if (meta != null && !meta.directPlay) {
        _notice('清晰度', '该格式无法原画直播，已建议使用转码档', tone: RbBannerTone.warning);
        return;
      }
    }
    if (_videoPresetSwitching) {
      return;
    }
    final prevPreset = _videoTranscodePreset;
    final prevVcodec = _videoTranscodeVcodec;
    final resumeAt = _absoluteVideoPosition();
    final resumePlaying = _videoWasPlaying();
    final wasTranscode = prevPreset != RbTranscodePreset.source.id;
    final toTranscode = presetId != RbTranscodePreset.source.id;
    rbLog.info('video preset switch resume=${resumeAt.inMilliseconds}ms hlsOff=$_videoHlsOffsetMs');
    final slot = _nav.current;
    _videoPresetSwitching = true;
    _videoTranscodePreset = presetId;
    try {
      _syncWireFromParent();
      _nav.wire.cancelPendingReads();
      if (wasTranscode && toTranscode && _videoCtrl != null) {
        await _swapTranscodePlayerAt(
          resumeAt: resumeAt,
          playing: resumePlaying,
          hint: '正在切换清晰度…',
          timeoutMsg: '切换清晰度超时，请稍后重试',
          presetSwitch: true,
        );
      } else {
        final epoch = ++_loadEpoch;
        await _detachVideoFromTree();
        if (!mounted || epoch != _loadEpoch) {
          return;
        }
        await _loadVideo(slot, epoch, presetSwitch: true, presetResumeAt: resumeAt, presetResumePlaying: resumePlaying);
      }
    } catch (e, st) {
      rbLog.warning('video preset switch failed', e, st);
      if (!mounted) {
        return;
      }
      _videoTranscodePreset = prevPreset;
      _videoTranscodeVcodec = prevVcodec;
      _notice('清晰度', '切换失败：${rbUserMessage(e)}', tone: RbBannerTone.warning);
      if (wasTranscode && prevPreset != RbTranscodePreset.source.id && _videoCtrl != null) {
        try {
          await _swapTranscodePlayerAt(
            resumeAt: resumeAt,
            playing: resumePlaying,
            hint: '正在恢复播放…',
            timeoutMsg: '恢复播放超时',
            presetSwitch: true,
          );
        } catch (e2, st2) {
          rbLog.warning('video preset recover failed', e2, st2);
        }
      }
    } finally {
      _videoPresetSwitching = false;
    }
  }

  Future<void> _loadVideo(
    RbPreviewFileSlot slot,
    int epoch, {
    bool presetSwitch = false,
    bool transcodeResume = false,
    bool transcodeSeek = false,
    Duration presetResumeAt = Duration.zero,
    bool presetResumePlaying = true,
  }) async {
    if (mounted && !presetSwitch) {
      final hint = transcodeSeek ? '正在跳转…' : (transcodeResume ? '正在缓冲…' : '正在准备播放…');
      _setVideoLoadLabel(hint);
    }
    _videoAgentHttpAvailable = _nav.wire.agentMediaPlayUri(slot.relPath) != null;
    if (_videoAgentHttpAvailable) {
      final meta = await _nav.wire.fetchAgentMediaMeta(slot.relPath);
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (meta != null) {
        if (meta.directPlay && _videoTranscodePreset == RbTranscodePreset.source.id) {
          // HEVC/AV1/VP8/VP9 等可原画直传
        } else if (!meta.directPlay) {
          if (_videoTranscodePreset == RbTranscodePreset.source.id) {
            _videoTranscodePreset = RbTranscodePreset.p720.id;
          }
          _videoTranscodeVcodec = RbTranscodeVcodec.pickForMeta(meta.videoCodec);
        }
      }
    }
    await _startRemoteVideo(
      epoch: epoch,
      entry: slot.entry,
      relPath: slot.relPath,
      totalSize: slot.entry.size,
      fileName: slot.entry.name,
      timeoutMsg: '视频加载超时，请返回后重试',
      presetSwitch: presetSwitch,
      transcodeResume: transcodeResume,
      transcodeSeek: transcodeSeek,
      presetResumeAt: presetResumeAt,
      presetResumePlaying: presetResumePlaying,
    );
  }

  Future<void> _startRemoteVideo({
    required int epoch,
    RbFileEntry? entry,
    required String relPath,
    required int totalSize,
    String? fileName,
    String? mime,
    int fileBaseOffset = 0,
    required String timeoutMsg,
    bool presetSwitch = false,
    bool transcodeResume = false,
    bool transcodeSeek = false,
    Duration presetResumeAt = Duration.zero,
    bool presetResumePlaying = true,
  }) async {
    RbVideoStreamServer? stream;
    VideoPlayerController? ctrl;
    var adopted = false;
    try {
      final store = _previewStore;
      final cached = entry != null && store != null
          ? await store.mediaFileIfComplete(entry, relPath, logicalTotal: totalSize, fileBaseOffset: fileBaseOffset)
          : null;
      final agentUri = cached == null ? _nav.wire.agentMediaPlayUri(relPath, fileBaseOffset: fileBaseOffset) : null;
      _videoAgentHttpAvailable = agentUri != null;
      _videoTimelineTotalMs = entry?.durationMs ?? 0;
      final transcodeReload = presetSwitch || transcodeResume;
      final beginMs = transcodeReload && presetResumeAt > Duration.zero
          ? (presetResumeAt.inMilliseconds ~/ rbTranscodeFragmentMs) * rbTranscodeFragmentMs
          : 0;
      final transcodeHlsUri = cached == null && _videoTranscodePreset != RbTranscodePreset.source.id
          ? _nav.wire.agentTranscodeHlsUri(
              relPath,
              presetId: _videoTranscodePreset,
              vcodec: _videoTranscodeVcodec,
              fileBaseOffset: fileBaseOffset,
              beginMs: beginMs,
            )
          : null;
      final transcodeStreamUri = cached == null && _videoTranscodePreset != RbTranscodePreset.source.id
          ? _nav.wire.agentTranscodeStreamUri(
              relPath,
              presetId: _videoTranscodePreset,
              vcodec: _videoTranscodeVcodec,
              fileBaseOffset: fileBaseOffset,
              startMs: beginMs,
            )
          : null;
      late VideoPlayerController player;
      if (cached != null) {
        rbLog.info('video play cache ${cached.path}');
        player = VideoPlayerController.file(cached);
        ctrl = player;
      } else if (transcodeHlsUri != null) {
        if (epoch == _loadEpoch && mounted) {
          final label = presetSwitch
              ? '正在切换清晰度…'
              : (transcodeSeek ? '正在跳转…' : (transcodeResume ? '正在缓冲…' : '正在准备转码流…'));
          _setVideoLoadLabel(label);
        }
        player = await _prepareTranscodePlayer(
          epoch: epoch,
          streamUri: transcodeStreamUri,
          hlsUri: transcodeHlsUri,
          transcodeReload: transcodeReload,
          presetSwitch: presetSwitch,
          transcodeSeek: transcodeSeek,
          transcodeResume: transcodeResume,
          beginMs: beginMs,
        );
        ctrl = player;
      } else if (agentUri != null) {
        rbLog.info('video play agent http (原画直读) $agentUri');
        player = VideoPlayerController.networkUrl(agentUri);
        ctrl = player;
      } else {
        rbLog.warning(
          'video play loopback+ICE 代理（无 Agent 媒体端口 ${RbDirectDialer.defaultPort}，大文件易卡）；'
          '请确认 Agent 已广播内网 IP 且防火墙放行直连端口',
        );
        stream = await RbVideoStreamServer.start(
          wire: _nav.wire,
          relPath: relPath,
          totalSize: totalSize,
          fileName: fileName,
          mime: mime,
          fileBaseOffset: fileBaseOffset,
          entry: entry,
          previewStore: store,
        );
        if (epoch != _loadEpoch || !mounted) {
          return;
        }
        player = VideoPlayerController.networkUrl(stream.playUri);
        ctrl = player;
      }
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      try {
        await player.initialize().timeout(
          Duration(seconds: transcodeReload ? 120 : 45),
          onTimeout: () => throw StateError(presetSwitch ? '切换清晰度超时，分片转码较慢请重试' : timeoutMsg),
        );
      } on PlatformException catch (e) {
        throw StateError(e.message ?? e.code);
      } catch (e) {
        if (transcodeHlsUri != null && !_videoTranscodeUseHls) {
          rbLog.warning('transcode stream init failed, fallback hls: $e');
          await player.dispose();
          player = await _prepareTranscodePlayer(
            epoch: epoch,
            streamUri: transcodeStreamUri,
            hlsUri: transcodeHlsUri,
            transcodeReload: transcodeReload,
            presetSwitch: presetSwitch,
            transcodeSeek: transcodeSeek,
            transcodeResume: transcodeResume,
            beginMs: beginMs,
            forceHls: true,
          );
          ctrl = player;
          await player.initialize().timeout(
            Duration(seconds: transcodeReload ? 120 : 45),
            onTimeout: () => throw StateError(presetSwitch ? '切换清晰度超时，分片转码较慢请重试' : timeoutMsg),
          );
        } else {
          rethrow;
        }
      }
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (transcodeHlsUri != null) {
        if (presetSwitch && presetResumeAt > Duration.zero) {
          _videoStreamPausedAbsMs = presetResumeAt.inMilliseconds;
          _videoStreamPlayAt = null;
        }
        _syncTranscodeTimeline(player, beginMs: beginMs);
        _videoStreamAnchorMs = _videoTranscodeUseHls ? 0 : beginMs;
      } else {
        _videoHlsOffsetMs = 0;
        _videoStreamAnchorMs = 0;
      }
      _videoStream = stream;
      _videoCtrl = player;
      adopted = true;
      _videoPlaybackRetry = 0;
      _videoPendingSeekMs = 0;
      _nav.thumbs?.setPaused(true);
      _attachVideoListeners(player, epoch);
      _setVideoLoadLabel('');
      final resumeSeek = transcodeReload && !presetSwitch
          ? Duration(milliseconds: presetResumeAt.inMilliseconds - beginMs)
          : Duration.zero;
      final resumePlaying = transcodeReload ? presetResumePlaying : true;
      if (transcodeReload && resumeSeek > Duration.zero) {
        await _resumeVideoPlayback(player, resumeSeek, play: false);
      }
      if (_videoTranscodeUseHls && transcodeHlsUri != null) {
        final warmUri = transcodeReload ? rbTranscodeUriCacheBust(transcodeHlsUri) : transcodeHlsUri;
        _lastHlsWarmChunk = beginMs ~/ rbTranscodeFragmentMs;
        rbWarmTranscodeAhead(warmUri, beginMs + 2 * rbTranscodeFragmentMs, fragmentCount: 12);
      }
      if (!mounted || epoch != _loadEpoch || _videoCtrl != player) {
        return;
      }
      if (resumePlaying) {
        try {
          await player.play();
          if (_usesAgentTranscode && !_videoTranscodeUseHls) {
            _videoStreamPlayAt = DateTime.now();
            _videoStreamPausedAbsMs = null;
            _videoWasPlayingFlag = true;
          }
        } catch (_) {}
      }
    } finally {
      if (!adopted) {
        if (stream != null) {
          await stream.close();
        }
        if (ctrl != null) {
          await ctrl.dispose();
        }
      }
    }
  }

  Future<void> _loadImage(RbPreviewFileSlot slot, int epoch) async {
    var fetchStarted = false;
    void applyOriginalState({required String mime, required int dataBytes, Uint8List? bytesHead}) {
      _imagePreviewMime = mime;
      _imageShowingOriginal = rbImagePreviewIsOriginal(
        fileName: slot.entry.name,
        mime: mime,
        dataBytes: dataBytes,
        entrySize: slot.entry.size,
        bytesHead: bytesHead,
      );
      _imageCanShowOriginalButton = !_imageShowingOriginal;
    }

    final store = _previewStore;
    if (store != null) {
      final hit = await store.load(slot.entry, slot.relPath, maxBytes: rbPreviewImageMaxBytes);
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (hit != null && hit.isNotEmpty && !rbIsLikelyStaleImagePreviewCache(slot.entry, hit.length)) {
        _bytes = hit;
        applyOriginalState(mime: '', dataBytes: hit.length, bytesHead: hit);
        return;
      }
    }
    fetchStarted = true;
    _beginPreviewFetch(expectedTotal: slot.entry.size > 0 ? slot.entry.size : rbPreviewImageMaxBytes);
    if (mounted) {
      setState(() {});
    }
    try {
      final r = await _nav.wire.readFileFull(
        slot.relPath,
        maxBytes: rbPreviewImageMaxBytes,
        onProgress: _onPreviewReadProgress,
      );
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (r.bytes.isEmpty) {
        throw StateError('预览数据为空');
      }
      if (!rbPreviewBytesDisplayable(fileName: slot.entry.name, mime: r.mime, bytesHead: r.bytes)) {
        throw StateError('预览格式无法显示（${r.mime.isEmpty ? slot.entry.name : r.mime}），请查看 Agent 转码日志');
      }
      _bytes = r.bytes;
      applyOriginalState(mime: r.mime, dataBytes: r.totalSize > 0 ? r.totalSize : r.bytes.length, bytesHead: r.bytes);
      if (store != null && r.bytes.isNotEmpty) {
        unawaited(store.put(slot.entry, slot.relPath, r.bytes));
      }
    } finally {
      if (fetchStarted && epoch == _loadEpoch) {
        _endPreviewFetch();
      }
    }
  }

  Future<void> _loadOriginalImage() async {
    if (_imageLoadingOriginal || !_imageCanShowOriginalButton || _imageShowingOriginal) {
      return;
    }
    final epoch = _loadEpoch;
    setState(() {
      _imageLoadingOriginal = true;
      _beginPreviewFetch(expectedTotal: _nav.current.entry.size > 0 ? _nav.current.entry.size : rbPreviewImageMaxBytes);
    });
    try {
      final slot = _nav.current;
      final r = await _nav.wire.readFileFullOriginal(
        slot.relPath,
        maxBytes: rbPreviewImageMaxBytes,
        onProgress: _onPreviewReadProgress,
      );
      if (epoch != _loadEpoch || !mounted) {
        return;
      }
      if (r.bytes.isEmpty || !rbPreviewBytesDisplayable(fileName: slot.entry.name, mime: r.mime, bytesHead: r.bytes)) {
        _notice('原图', '高清 WebP 加载失败，请重试或下载原文件', tone: RbBannerTone.warning, duration: const Duration(seconds: 4));
        return;
      }
      setState(() {
        _bytes = r.bytes;
        _imagePreviewMime = r.mime;
        _imageShowingOriginal = true;
        _imageCanShowOriginalButton = false;
      });
    } catch (e, st) {
      rbLog.warning('load original image', e, st);
      if (rbIsReadSuperseded(e)) {
        return;
      }
      _toast(rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 5));
    } finally {
      if (mounted && _imageLoadingOriginal) {
        setState(() {
          _imageLoadingOriginal = false;
          _endPreviewFetch();
        });
      }
    }
  }

  Future<void> _jumpToPreviewIndex(int index) async {
    if (_fileSwitching || _disposing || !mounted || index == _nav.fileIndex) {
      return;
    }
    if (index < 0 || index >= _nav.files.length || !rbIsPreviewableEntry(_nav.files[index].entry)) {
      return;
    }
    _fileSwitching = true;
    _syncWireFromParent();
    try {
      if (_motionPlaying) {
        _stopMotionClip();
      }
      _nav.fileIndex = index;
      await _loadCurrent();
    } catch (e, st) {
      rbLog.warning('preview strip pick', e, st);
      if (rbIsReadSuperseded(e)) {
        return;
      }
      _syncWireFromParent();
      if (rbIsWireLost(e) && !_nav.wire.isOpen) {
        unawaited(_handleWireLost(e));
      } else {
        _toast(rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 5));
      }
    } finally {
      _fileSwitching = false;
    }
  }

  Future<void> _onPreviewFileSwipe(bool next) async {
    if (_fileSwitching || _disposing || !mounted) {
      return;
    }
    _fileSwitching = true;
    _syncWireFromParent();
    try {
      final (ok, msg) = next ? await _nav.nextFile() : await _nav.prevFile();
      if (!ok) {
        _toast(msg);
        return;
      }
      await _loadCurrent();
    } catch (e, st) {
      rbLog.warning('preview file swipe', e, st);
      if (rbIsReadSuperseded(e)) {
        return;
      }
      _syncWireFromParent();
      if (rbIsWireLost(e) && !_nav.wire.isOpen) {
        unawaited(_handleWireLost(e));
      } else {
        _toast(rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 5));
      }
    } finally {
      _fileSwitching = false;
    }
  }

  void _toast(String? msg, {RbBannerTone tone = RbBannerTone.info, Duration duration = const Duration(seconds: 2)}) {
    if (msg == null || !mounted) {
      return;
    }
    rbFlashTopNotice(context, msg, tone: tone, duration: duration);
  }

  void _notice(String title, String message, {RbBannerTone tone = RbBannerTone.info, Duration duration = const Duration(seconds: 3)}) {
    if (!mounted) {
      return;
    }
    rbFlashTopNotice(context, '$title：$message', tone: tone, duration: duration);
  }

  bool get _immersiveMedia => _kind != RbPreviewKind.unsupported;

  bool get _videoKind => _kind == RbPreviewKind.video;

  void _toggleImmersiveChrome() {
    if (!mounted) {
      return;
    }
    setState(() => _showImmersiveChrome = !_showImmersiveChrome);
  }

  bool _canPreviewFileStep(bool next) => next ? _nav.canNextFile : _nav.canPrevFile;

  /// 加载/占位阶段也允许横竖滑切文件（与 [RbZoomableImage] / 缩略图占位一致）。
  Widget _wrapPreviewFileSwipe(Widget child, {VoidCallback? onTap, bool? horizontalSwipe}) {
    final horiz = horizontalSwipe ?? !_disableHorizontalFileSwipe;
    return RbPreviewFileSwipeListener(
      onTap: onTap,
      onVerticalFileSwipe: (next) => unawaited(_onPreviewFileSwipe(next)),
      canVerticalFileSwipe: _canPreviewFileStep,
      onHorizontalFileSwipe: horiz ? (next) => unawaited(_onPreviewFileSwipe(next)) : null,
      canHorizontalFileSwipe: _canPreviewFileStep,
      child: child,
    );
  }

  bool get _imageFullScreenVerticalNav =>
      _kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto;

  /// 渲染态 Markdown/HTML 含横向滚动，禁用横滑切文件。
  bool get _disableHorizontalFileSwipe =>
      (_kind == RbPreviewKind.markdown && _renderMarkdown) || (_kind == RbPreviewKind.html && _renderHtml);

  Widget _wrapEdgeFileNav(Widget child) {
    if ((_loading && !_canShowImageWhileLoading) ||
        _error.isNotEmpty ||
        _videoKind ||
        _imageFullScreenVerticalNav) {
      return child;
    }
    return RbPreviewEdgeSwipeLayer(
      leftEdgeVerticalEnabled: !_docScrollKind && !_docWebKind,
      edgeVerticalSwipeEnabled: !_imageFullScreenVerticalNav && !_docWebKind,
      canTriggerVerticalSwipe: _canEdgeVerticalSwipe,
      onEdgeVerticalSwipe: (next) => unawaited(_onPreviewFileSwipe(next)),
      onHorizontalFileSwipe: _disableHorizontalFileSwipe ? null : (next) => unawaited(_onPreviewFileSwipe(next)),
      canTriggerHorizontalFileSwipe: _canPreviewFileStep,
      child: child,
    );
  }

  Widget _chromeGradient({required bool top}) {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: top ? Alignment.topCenter : Alignment.bottomCenter,
          end: top ? Alignment.bottomCenter : Alignment.topCenter,
          colors: [Colors.black.withValues(alpha: 0.72), Colors.transparent],
        ),
      ),
    );
  }

  Widget _previewFolderStrip() {
    return RbPreviewFolderStrip(
      nav: _nav,
      thumbs: _nav.thumbs,
      compact: true,
      onPick: (i) => unawaited(_jumpToPreviewIndex(i)),
    );
  }

  Widget _immersiveBottomChrome() {
    return DecoratedBox(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.bottomCenter,
          end: Alignment.topCenter,
          colors: [Colors.black.withValues(alpha: 0.75), Colors.transparent],
        ),
      ),
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [_previewFolderStrip(), _videoDeleteChrome()],
        ),
      ),
    );
  }

  Widget _previewLoadProgressBar() {
    final value = _previewLoadProgressValue;
    final pct = value != null ? (value * 100).round() : null;
    return IgnorePointer(
      child: DecoratedBox(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.bottomCenter,
            end: Alignment.topCenter,
            colors: [Colors.black.withValues(alpha: 0.55), Colors.transparent],
          ),
        ),
        child: SafeArea(
          top: false,
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 12, 20, 10),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                if (pct != null)
                  Text('加载预览 $pct%', textAlign: TextAlign.center, style: const TextStyle(color: Colors.white70, fontSize: 12)),
                const SizedBox(height: 6),
                LinearProgressIndicator(
                  value: value,
                  minHeight: 3,
                  backgroundColor: Colors.white24,
                  color: Colors.white,
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _photoLoadBackButton() {
    return SafeArea(
      bottom: false,
      child: Align(
        alignment: Alignment.topLeft,
        child: IconButton(onPressed: _popPreview, icon: const Icon(Icons.arrow_back, color: Colors.white)),
      ),
    );
  }

  Widget _immersiveChromeOverlay(String pos) {
    if (_docImmersiveKind) {
      return Stack(
        children: [
          Positioned(left: 0, right: 0, bottom: 0, child: _immersiveBottomChrome()),
        ],
      );
    }
    return Stack(
      children: [
        Positioned(left: 0, right: 0, top: 0, child: Stack(children: [
          Positioned(left: 0, right: 0, top: 0, height: 120, child: _chromeGradient(top: true)),
          SafeArea(bottom: false, child: _videoTopChrome(pos)),
        ])),
        Positioned(left: 0, right: 0, bottom: 0, child: _immersiveBottomChrome()),
      ],
    );
  }

  static const _bottomBarH = 56.0;

  void _notifyListChanged() {
    _listChanged = true;
    widget.onListChanged?.call();
  }

  void _popPreview() {
    try {
      AppNavigator.pop<Object?>(_listChanged ? true : null);
    } catch (e, st) {
      rbLog.warning('pop preview failed', e, st);
    }
  }

  void _openFileInfoPanel() {
    final rel = _nav.current.relPath;
    unawaited(RbFileInfoTopPanel.show(
      context,
      nav: _nav,
      previewKind: _kind,
      onNavigateToDir: (dir) => AppNavigator.pop<Object?>(RbPreviewNavigateToDir(dir)),
      onDownload: () => rbRunFileInfoDownload(context, wire: _nav.wire, relPath: rel, entry: _nav.current.entry),
    ));
  }

  Widget _fileInfoTapTarget({required Widget child}) {
    return InkWell(onTap: _openFileInfoPanel, child: child);
  }

  bool get _bottomBusy => _deleting || _downloading || _imageLoadingOriginal;

  bool get _showImageOriginalButton =>
      !_imageShowingOriginal &&
      _imageCanShowOriginalButton &&
      (_kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto);

  Future<void> _downloadCurrent() async {
    final slot = _nav.current;
    setState(() => _downloading = true);
    try {
      await rbRunFileInfoDownload(context, wire: _nav.wire, relPath: slot.relPath, entry: slot.entry);
    } finally {
      if (mounted) {
        setState(() => _downloading = false);
      }
    }
  }

  Future<void> _deleteCurrent() async {
    final slot = _nav.current;
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除文件'),
        content: Text('确定删除？\n${_nav.displayFilePath}'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(
            style: FilledButton.styleFrom(backgroundColor: Theme.of(ctx).colorScheme.error),
            onPressed: () => Navigator.pop(ctx, true),
            child: const Text('删除'),
          ),
        ],
      ),
    );
    if (ok != true || !mounted) {
      return;
    }
    setState(() => _deleting = true);
    try {
      await _nav.wire.deleteFile(slot.relPath);
      _nav.thumbs?.evict(slot.relPath, entry: slot.entry);
      unawaited(_previewStore?.evict(slot.relPath));
      if (!mounted) {
        return;
      }
      _notifyListChanged();
      setState(() => _deleting = false);
      if (!_nav.removeCurrentFile()) {
        _popPreview();
        return;
      }
      await _nav.focusPreviewableFile();
      await _loadCurrent();
      if (mounted) {
        _toast('已删除');
      }
    } catch (e) {
      if (rbIsWireLost(e)) {
        unawaited(_handleWireLost(e));
        return;
      }
      _toast(rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 5));
    } finally {
      if (mounted && _deleting) {
        setState(() => _deleting = false);
      }
    }
  }

  Widget _bottomToolbar() {
    return SafeArea(
      top: false,
      child: SizedBox(height: _bottomBarH, child: Center(child: _previewBottomActionButtons())),
    );
  }

  Widget _previewBottomActionButtons({bool onVideoChrome = false}) {
    final busy = _bottomBusy;
    final err = onVideoChrome ? Colors.redAccent : Theme.of(context).colorScheme.error;
    final fg = onVideoChrome ? Colors.white : null;
    Widget dlIcon() => _downloading
        ? SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: fg))
        : Icon(Icons.download_outlined, color: fg);
    Widget delIcon() => _deleting
        ? SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: err))
        : Icon(Icons.delete_outline, color: err);
    Widget origIcon() => _imageLoadingOriginal
        ? SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: fg))
        : Icon(Icons.photo_outlined, color: fg);
    final buttons = <Widget>[];
    if (_showImageOriginalButton) {
      buttons.add(FilledButton.tonalIcon(
        onPressed: busy ? null : () => unawaited(_loadOriginalImage()),
        icon: origIcon(),
        label: Text(_imageLoadingOriginal ? '加载中…' : '原图', style: fg != null ? TextStyle(color: fg) : null),
      ));
      buttons.add(const SizedBox(width: 12));
    }
    buttons.addAll([
      FilledButton.tonalIcon(
        onPressed: busy ? null : () => unawaited(_downloadCurrent()),
        icon: dlIcon(),
        label: Text(_downloading ? '下载中…' : '下载', style: fg != null ? TextStyle(color: fg) : null),
      ),
      const SizedBox(width: 12),
      FilledButton.tonalIcon(
        onPressed: busy ? null : () => unawaited(_deleteCurrent()),
        icon: delIcon(),
        label: Text(_deleting ? '删除中…' : '删除', style: TextStyle(color: err)),
      ),
    ]);
    return Row(mainAxisSize: MainAxisSize.min, children: buttons);
  }

  Widget _videoDeleteChrome() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      child: _previewBottomActionButtons(onVideoChrome: true),
    );
  }

  Widget _videoBottomBarTrailing() {
    final busy = _bottomBusy;
    Widget icon(bool loading, IconData data, Color color) => loading
        ? SizedBox(width: 22, height: 22, child: CircularProgressIndicator(strokeWidth: 2, color: color))
        : Icon(data, color: color);
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        IconButton(
          visualDensity: VisualDensity.compact,
          padding: EdgeInsets.zero,
          constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
          tooltip: _downloading ? '下载中…' : '下载',
          onPressed: busy ? null : () => unawaited(_downloadCurrent()),
          icon: icon(_downloading, Icons.download_outlined, Colors.white),
        ),
        IconButton(
          visualDensity: VisualDensity.compact,
          padding: EdgeInsets.zero,
          constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
          tooltip: _deleting ? '删除中…' : '删除',
          onPressed: busy ? null : () => unawaited(_deleteCurrent()),
          icon: icon(_deleting, Icons.delete_outline, Colors.redAccent),
        ),
      ],
    );
  }

  Widget _videoTopChrome(String pos) {
    return Row(
      children: [
        IconButton(onPressed: _popPreview, icon: const Icon(Icons.arrow_back, color: Colors.white)),
        Expanded(
          child: _fileInfoTapTarget(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  _nav.current.entry.name,
                  maxLines: 1,
                  overflow: TextOverflow.ellipsis,
                  style: const TextStyle(color: Colors.white, fontSize: 15),
                ),
                if (pos.isNotEmpty) Text(pos, style: const TextStyle(color: Colors.white70, fontSize: 12)),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _videoBottomFileNav() {
    VoidCallback? tap(VoidCallback f) => _fileSwitching ? null : f;
    Widget navBtn({required String tip, required IconData icon, required bool enabled, required VoidCallback onTap}) {
      return IconButton(
        visualDensity: VisualDensity.compact,
        padding: EdgeInsets.zero,
        constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
        tooltip: tip,
        icon: Icon(icon, color: Colors.white),
        onPressed: enabled ? tap(onTap) : null,
      );
    }

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        navBtn(
          tip: '上一个文件',
          icon: Icons.chevron_left,
          enabled: _nav.canPrevFile,
          onTap: () => unawaited(_onPreviewFileSwipe(false)),
        ),
        navBtn(
          tip: '下一个文件',
          icon: Icons.chevron_right,
          enabled: _nav.canNextFile,
          onTap: () => unawaited(_onPreviewFileSwipe(true)),
        ),
      ],
    );
  }

  Widget? _buildRenderToggle() {
    if (_kind == RbPreviewKind.markdown) {
      final rendered = _renderMarkdown;
      return IconButton(
        tooltip: rendered ? '切到源码' : '切到渲染',
        onPressed: () {
          setState(() => _renderMarkdown = !_renderMarkdown);
          _resetScrollEdges();
        },
        icon: Icon(rendered ? Icons.code : Icons.article_outlined),
      );
    }
    if (_kind == RbPreviewKind.html) {
      final rendered = _renderHtml;
      return IconButton(
        tooltip: rendered ? '切到源码' : '切到渲染',
        onPressed: () {
          setState(() => _renderHtml = !_renderHtml);
          _resetScrollEdges();
        },
        icon: Icon(rendered ? Icons.code : Icons.language),
      );
    }
    return null;
  }

  List<Widget> _navActions() {
    final toggle = _buildRenderToggle();
    return toggle == null ? const [] : [toggle];
  }

  PreferredSizeWidget? _buildPreviewAppBar(String pos) {
    if (!_immersiveMedia) {
      return AppBar(
        leading: BackButton(onPressed: _popPreview),
        title: _previewTitleColumn(pos),
        actions: _navActions(),
      );
    }
    if (_videoKind || _kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto) {
      return null;
    }
    if (_docImmersiveKind) {
      return AppBar(
        backgroundColor: Colors.black,
        foregroundColor: Colors.white,
        leading: BackButton(onPressed: _popPreview, color: Colors.white),
        title: _previewTitleColumn(pos, onDark: true),
        actions: _navActions(),
      );
    }
    return null;
  }

  Widget _previewTitleColumn(String pos, {bool onDark = false}) {
    final subStyle = onDark ? const TextStyle(color: Colors.white70, fontSize: 12) : null;
    return _fileInfoTapTarget(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            _nav.current.entry.name,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: TextStyle(fontSize: 15, height: 1.25, color: onDark ? Colors.white : null),
          ),
          Text(
            _nav.displayFilePath,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
            style: subStyle ?? Theme.of(context).textTheme.bodySmall,
          ),
          if (pos.isNotEmpty) Text(pos, style: subStyle ?? Theme.of(context).textTheme.bodySmall),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final pos = _nav.files.isEmpty ? '' : '${_nav.fileIndex + 1}/${_nav.files.length}';
    return PopScope(
      canPop: false,
      onPopInvokedWithResult: (didPop, _) {
        if (didPop) {
          return;
        }
        try {
          _popPreview();
        } catch (e, st) {
          rbLog.warning('pop invoked', e, st);
        }
      },
      child: Scaffold(
        backgroundColor: _immersiveMedia ? Colors.black : null,
        appBar: _buildPreviewAppBar(pos),
        body: _immersiveMedia ? _buildImmersiveBody(pos) : _wrapEdgeFileNav(_buildBody()),
        bottomNavigationBar: _immersiveMedia ? null : _bottomToolbar(),
      ),
    );
  }

  Widget _buildImmersiveBody(String pos) {
    if (_videoKind) {
      return _buildVideoScaffoldBody(pos);
    }
    if (_kind == RbPreviewKind.image || _kind == RbPreviewKind.motionPhoto) {
      return _buildPhotoImmersiveBody(pos);
    }
    return _buildDocImmersiveBody(pos);
  }

  /// 文本 / Markdown / HTML / PDF：单击显隐底栏缩略图条与下载/删除。
  Widget _buildDocImmersiveBody(String pos) {
    if (_loading) {
      return Stack(
        fit: StackFit.expand,
        children: [
          Positioned.fill(
            child: _wrapPreviewFileSwipe(
              const Center(child: CircularProgressIndicator(color: Colors.white70)),
              onTap: _toggleImmersiveChrome,
            ),
          ),
          _photoLoadBackButton(),
          if (_showPreviewLoadBar) Positioned(left: 0, right: 0, bottom: 0, child: _previewLoadProgressBar()),
        ],
      );
    }
    if (_error.isNotEmpty) {
      return Stack(
        children: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(_error, textAlign: TextAlign.center, style: const TextStyle(color: Colors.white70)),
                  const SizedBox(height: 16),
                  FilledButton(onPressed: () => unawaited(_loadCurrent()), child: const Text('重试')),
                ],
              ),
            ),
          ),
          SafeArea(
            child: Align(
              alignment: Alignment.topLeft,
              child: IconButton(onPressed: _popPreview, icon: const Icon(Icons.arrow_back, color: Colors.white)),
            ),
          ),
        ],
      );
    }
    return Stack(
      fit: StackFit.expand,
      children: [
        Positioned.fill(
          child: _RbPreviewLightTapDetector(
            onLightTap: _toggleImmersiveChrome,
            child: _wrapEdgeFileNav(_buildBody()),
          ),
        ),
        if (_showImmersiveChrome) _immersiveChromeOverlay(pos),
      ],
    );
  }

  Widget _buildPhotoImmersiveBody(String pos) {
    if (_error.isNotEmpty) {
      return Stack(
        children: [
          Center(child: Text(_error, textAlign: TextAlign.center, style: const TextStyle(color: Colors.white70))),
          _photoLoadBackButton(),
        ],
      );
    }
    final waitingContent = _loading && !_canShowImageWhileLoading;
    return Stack(
      fit: StackFit.expand,
      children: [
        Positioned.fill(
          child: waitingContent
              ? _wrapPreviewFileSwipe(
                  const Center(child: CircularProgressIndicator(color: Colors.white70)),
                  onTap: _toggleImmersiveChrome,
                )
              : _wrapEdgeFileNav(_buildBody()),
        ),
        if (!_showImmersiveChrome) _photoLoadBackButton(),
        if (_showImmersiveChrome) _immersiveChromeOverlay(pos),
        if (_showPreviewLoadBar) Positioned(left: 0, right: 0, bottom: 0, child: _previewLoadProgressBar()),
      ],
    );
  }

  Widget _buildVideoScaffoldBody(String pos) {
    if (_loading) {
      return _wrapPreviewFileSwipe(
        const Center(child: CircularProgressIndicator(color: Colors.white70)),
      );
    }
    if (_error.isNotEmpty) {
      return Stack(
        children: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Text(_error, textAlign: TextAlign.center, style: const TextStyle(color: Colors.white70)),
                  const SizedBox(height: 16),
                  FilledButton(onPressed: () => unawaited(_loadCurrent()), child: const Text('重试')),
                ],
              ),
            ),
          ),
          SafeArea(
            child: Align(
              alignment: Alignment.topLeft,
              child: IconButton(
                onPressed: _popPreview,
                icon: const Icon(Icons.arrow_back, color: Colors.white),
              ),
            ),
          ),
        ],
      );
    }
    return _buildVideo(pos);
  }

  Widget _buildBody() {
    if (_loading && !_canShowImageWhileLoading) {
      return _wrapPreviewFileSwipe(const Center(child: CircularProgressIndicator()));
    }
    if (_error.isNotEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(_error, textAlign: TextAlign.center),
              const SizedBox(height: 16),
              FilledButton(onPressed: () => unawaited(_loadCurrent()), child: const Text('重试')),
            ],
          ),
        ),
      );
    }
    return switch (_kind) {
      RbPreviewKind.image => _buildImage(),
      RbPreviewKind.motionPhoto => _buildMotionPhoto(),
      RbPreviewKind.markdown => _renderMarkdown
          ? NotificationListener<ScrollNotification>(
              onNotification: _onScrollNotification,
              child: Markdown(data: _text),
            )
          : SingleChildScrollView(
              controller: _textScroll,
              padding: const EdgeInsets.fromLTRB(12, 12, 12, 48),
              child: SelectableText(_text, style: const TextStyle(fontFamily: 'monospace', fontSize: 13)),
            ),
      RbPreviewKind.text => SingleChildScrollView(
          controller: _textScroll,
          padding: const EdgeInsets.fromLTRB(12, 12, 12, 48),
          child: SelectableText(_text, style: const TextStyle(fontFamily: 'monospace', fontSize: 13)),
        ),
      RbPreviewKind.html => _renderHtml
          ? WebViewWidget(
              controller: WebViewController()
                ..setJavaScriptMode(JavaScriptMode.disabled)
                ..loadHtmlString(_text),
            )
          : SingleChildScrollView(
              controller: _textScroll,
              padding: const EdgeInsets.fromLTRB(12, 12, 12, 48),
              child: SelectableText(_text, style: const TextStyle(fontFamily: 'monospace', fontSize: 13)),
            ),
      RbPreviewKind.pdf => _buildPdf(),
      RbPreviewKind.video => const SizedBox.shrink(),
      RbPreviewKind.unsupported => Center(child: rbFileTypeIconWidget(_nav.current.entry.name, size: 72)),
    };
  }

  Widget _buildPdf() {
    final uri = _pdfMediaUri;
    if (uri == null) {
      return Center(child: Text(_error.isNotEmpty ? _error : 'PDF 未就绪', textAlign: TextAlign.center));
    }
    return RbPdfPreview(mediaUri: uri);
  }

  Widget _buildVideo(String pos) {
    final ctrl = _videoCtrl;
    if (ctrl == null) {
      final hint = _videoLoadLabel.isEmpty ? '正在准备播放…' : _videoLoadLabel;
      return Stack(
        fit: StackFit.expand,
        children: [
          const Center(
            child: SizedBox(
              width: 36,
              height: 36,
              child: CircularProgressIndicator(strokeWidth: 2.5, color: Colors.white70),
            ),
          ),
          Positioned(top: 0, left: 0, right: 0, child: RbStatusOverlay(message: hint, onDark: true, showProgress: false)),
        ],
      );
    }
    return RbRemoteVideoPlayer(
      key: ValueKey('v$_loadEpoch-${_nav.current.relPath}'),
      controller: ctrl,
      topChrome: _videoTopChrome(pos),
      folderStrip: _previewFolderStrip(),
      bottomBarLeading: _videoBottomFileNav(),
      bottomBarTrailing: _videoBottomBarTrailing(),
      transcodePresetId: _videoTranscodePreset,
      transcodeAgentAvailable: _videoAgentHttpAvailable,
      onTranscodePreset: (id) => unawaited(_onVideoTranscodePreset(id)),
      onHlsSeek: _usesAgentTranscode ? _seekTranscodeHls : null,
      timelineOffsetMs: _videoHlsOffsetMs,
      timelineTotalMs: _videoTimelineTotalMs,
      playbackAnchorMs: _videoStreamAnchorMs,
      timelineUsePlayClock: _usesAgentTranscode && !_videoTranscodeUseHls,
      timelinePausedAbsPosMs: _videoStreamPausedAbsMs,
      timelinePlayStartedAt: _videoStreamPlayAt,
      hintLabel: _videoLoadLabel,
      hideChromeInitially: true,
      onEdgeVerticalSwipe: (next) => unawaited(_onPreviewFileSwipe(next)),
    );
  }

  Widget _buildImage() => _buildPreviewImageLayer();

  Widget _buildPreviewImageLayer() {
    final slot = _nav.current;
    return _RbPreviewImageLayer(
      key: ValueKey(slot.relPath),
      thumbBytes: _imageThumbBytes,
      previewBytes: _bytes,
      fileName: slot.entry.name,
      previewMime: _imagePreviewMime,
      onTap: _toggleImmersiveChrome,
      onVerticalFileSwipe: (next) => unawaited(_onPreviewFileSwipe(next)),
      canVerticalFileSwipe: _canPreviewFileStep,
      onHorizontalFileSwipe: (next) => unawaited(_onPreviewFileSwipe(next)),
      canHorizontalFileSwipe: _canPreviewFileStep,
    );
  }

  Future<void> _playMotionClip() async {
    final slot = _nav.current;
    final rel = rbMotionVideoRelPath(slot.relPath, slot.entry);
    if (rel == null) {
      return;
    }
    final epoch = ++_loadEpoch;
    if (mounted) {
      setState(() => _motionPlaying = true);
    }
    await _disposeVideo();
    if (epoch != _loadEpoch || !mounted) {
      return;
    }
    try {
      await _startRemoteVideo(
        epoch: epoch,
        entry: slot.entry,
        relPath: rel,
        totalSize: rbMotionVideoTotalSize(slot.entry),
        fileName: slot.entry.motionCompanion.isNotEmpty ? slot.entry.motionCompanion : slot.entry.name,
        mime: rbMotionVideoMime(slot.entry),
        fileBaseOffset: rbMotionVideoBaseOffset(slot.entry),
        timeoutMsg: '动态照片视频加载超时',
      );
      if (mounted) {
        setState(() {});
      }
    } catch (e) {
      if (mounted) {
        setState(() => _motionPlaying = false);
        _toast(rbUserMessage(e), tone: RbBannerTone.warning, duration: const Duration(seconds: 5));
      }
    }
  }

  void _stopMotionClip() {
    if (!_motionPlaying || !mounted) {
      return;
    }
    setState(() => _motionPlaying = false);
    unawaited(_disposeVideo());
  }

  bool get _showMotionVideo =>
      _motionPlaying && _videoCtrl != null && _videoCtrl!.value.isInitialized;

  Widget _buildMotionPhoto() {
    final path = _nav.current.relPath;
    final ctrl = _videoCtrl;
    return GestureDetector(
      onLongPressStart: (_) => unawaited(_playMotionClip()),
      onLongPressEnd: (_) => _stopMotionClip(),
      onLongPressCancel: _stopMotionClip,
      child: Stack(
        fit: StackFit.expand,
        children: [
          _buildPreviewImageLayer(),
          if (_showMotionVideo)
            Center(
              child: FittedBox(
                fit: BoxFit.contain,
                child: SizedBox(
                  width: ctrl!.value.size.width > 0 ? ctrl.value.size.width : 1,
                  height: ctrl.value.size.height > 0 ? ctrl.value.size.height : 1,
                  child: VideoPlayer(ctrl, key: ValueKey('motion-$path')),
                ),
              ),
            ),
          if (!_showImmersiveChrome)
            Positioned(
              left: 12,
              right: 12,
              bottom: 24,
              child: Center(
                child: IgnorePointer(
                  child: DecoratedBox(
                    decoration: BoxDecoration(color: Colors.black54, borderRadius: BorderRadius.circular(8)),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                      child: Text(
                        _motionPlaying ? '松手停止' : '长按播放 Live / 动态照片',
                        style: const TextStyle(color: Colors.white, fontSize: 12),
                      ),
                    ),
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}

/// 预览图加载前用列表缩略图占位，预览就绪后淡入（底图仅展示，手势在预览层）。
class _RbPreviewImageLayer extends StatefulWidget {
  const _RbPreviewImageLayer({
    super.key,
    required this.thumbBytes,
    required this.previewBytes,
    required this.fileName,
    this.previewMime = '',
    this.onTap,
    this.onVerticalFileSwipe,
    this.canVerticalFileSwipe,
    this.onHorizontalFileSwipe,
    this.canHorizontalFileSwipe,
  });

  final Uint8List? thumbBytes;
  final Uint8List? previewBytes;
  final String fileName;
  final String previewMime;
  final VoidCallback? onTap;
  final void Function(bool next)? onVerticalFileSwipe;
  final bool Function(bool next)? canVerticalFileSwipe;
  final void Function(bool next)? onHorizontalFileSwipe;
  final bool Function(bool next)? canHorizontalFileSwipe;

  @override
  State<_RbPreviewImageLayer> createState() => _RbPreviewImageLayerState();
}

class _RbPreviewImageLayerState extends State<_RbPreviewImageLayer> with SingleTickerProviderStateMixin {
  static const _fadeMs = 220;

  late final AnimationController _fadeCtrl;
  late final Animation<double> _fade;

  @override
  void initState() {
    super.initState();
    _fadeCtrl = AnimationController(vsync: this, duration: const Duration(milliseconds: _fadeMs));
    _fade = CurvedAnimation(parent: _fadeCtrl, curve: Curves.easeOut);
    _applyPreviewFade();
  }

  @override
  void dispose() {
    _fadeCtrl.dispose();
    super.dispose();
  }

  @override
  void didUpdateWidget(covariant _RbPreviewImageLayer oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (!identical(widget.previewBytes, oldWidget.previewBytes)) {
      _applyPreviewFade();
    }
  }

  void _applyPreviewFade() {
    final preview = widget.previewBytes;
    if (preview == null || preview.isEmpty) {
      _fadeCtrl.value = 0;
      return;
    }
    final thumb = widget.thumbBytes;
    if (thumb == null || thumb.isEmpty) {
      _fadeCtrl.value = 1;
      return;
    }
    _fadeCtrl.forward(from: 0);
  }

  @override
  Widget build(BuildContext context) {
    final thumb = widget.thumbBytes;
    final preview = widget.previewBytes;
    if (thumb == null && preview == null) {
      return Center(child: rbFileTypeIconWidget(widget.fileName, size: 72));
    }
    final hasThumb = thumb != null && thumb.isNotEmpty;
    final hasPreview = preview != null && preview.isNotEmpty;
    if (!hasThumb && !hasPreview) {
      return Center(child: rbFileTypeIconWidget(widget.fileName, size: 72));
    }
    return AnimatedBuilder(
      animation: _fade,
      builder: (context, _) {
        final previewOpaque = hasPreview && _fade.value >= 0.05;
        final thumbInteractive = hasThumb && (!hasPreview || _fade.value < 0.95);
        return Stack(
          fit: StackFit.expand,
          alignment: Alignment.center,
          children: [
            if (hasPreview)
              IgnorePointer(
                ignoring: !previewOpaque,
                child: FadeTransition(
                  opacity: _fade,
                  child: RbZoomableImage(
                    bytes: preview,
                    fileName: widget.fileName,
                    previewMime: widget.previewMime,
                    onTap: widget.onTap,
                    onVerticalFileSwipe: widget.onVerticalFileSwipe,
                    canVerticalFileSwipe: widget.canVerticalFileSwipe,
                    onHorizontalFileSwipe: widget.onHorizontalFileSwipe,
                    canHorizontalFileSwipe: widget.canHorizontalFileSwipe,
                  ),
                ),
              ),
            if (hasThumb)
              IgnorePointer(
                ignoring: !thumbInteractive,
                child: FadeTransition(
                  opacity: Tween<double>(begin: 1, end: 0).animate(_fade),
                  child: RbPreviewFileSwipeListener(
                    onTap: widget.onTap,
                    onVerticalFileSwipe: widget.onVerticalFileSwipe,
                    canVerticalFileSwipe: widget.canVerticalFileSwipe,
                    onHorizontalFileSwipe: widget.onHorizontalFileSwipe,
                    canHorizontalFileSwipe: widget.canHorizontalFileSwipe,
                    child: RbMemoryPicture(
                      bytes: thumb,
                      fileName: widget.fileName,
                      fit: BoxFit.contain,
                      filterQuality: FilterQuality.low,
                      gaplessPlayback: true,
                    ),
                  ),
                ),
              ),
          ],
        );
      },
    );
  }
}

/// 轻点（几乎无位移）时回调；滚动/拖拽不触发，供文本类预览显隐底栏。
class _RbPreviewLightTapDetector extends StatefulWidget {
  const _RbPreviewLightTapDetector({required this.onLightTap, required this.child});

  final VoidCallback onLightTap;
  final Widget child;

  @override
  State<_RbPreviewLightTapDetector> createState() => _RbPreviewLightTapDetectorState();
}

class _RbPreviewLightTapDetectorState extends State<_RbPreviewLightTapDetector> {
  static const _slop = 36.0;
  int? _ptr;
  Offset? _down;

  void _onDown(PointerDownEvent e) {
    _ptr = e.pointer;
    _down = e.position;
  }

  void _onUp(PointerUpEvent e) {
    if (_ptr != e.pointer) {
      return;
    }
    _ptr = null;
    final down = _down;
    _down = null;
    if (down != null && (e.position - down).distance <= _slop) {
      widget.onLightTap();
    }
  }

  void _onCancel(PointerCancelEvent e) {
    if (_ptr == e.pointer) {
      _ptr = null;
      _down = null;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Listener(
      onPointerDown: _onDown,
      onPointerUp: _onUp,
      onPointerCancel: _onCancel,
      behavior: HitTestBehavior.translucent,
      child: widget.child,
    );
  }
}
