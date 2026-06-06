/// 转码连续流在部分平台上 [VideoPlayerController.value.position] 长期为 0，用锚点+播放时钟合成全片时间。
int rbVideoAbsPosMs({
  required int playerPosMs,
  required int timelineOffsetMs,
  required int playbackAnchorMs,
  required bool isPlaying,
  DateTime? playStartedAt,
  int? pausedAbsPosMs,
  double playbackSpeed = 1.0,
}) {
  if (timelineOffsetMs > 0) {
    if (pausedAbsPosMs != null && !isPlaying && playerPosMs < 1500) {
      return pausedAbsPosMs;
    }
    return playerPosMs + timelineOffsetMs;
  }
  final useClock = playStartedAt != null || pausedAbsPosMs != null;
  if (useClock) {
    if (playerPosMs >= playbackAnchorMs + 1500 && playerPosMs > 800) {
      return playerPosMs;
    }
    if (!isPlaying && pausedAbsPosMs != null) {
      return pausedAbsPosMs;
    }
    if (isPlaying && playStartedAt != null) {
      final el = DateTime.now().difference(playStartedAt).inMilliseconds;
      return playbackAnchorMs + (el * playbackSpeed).round();
    }
  }
  return playerPosMs;
}
