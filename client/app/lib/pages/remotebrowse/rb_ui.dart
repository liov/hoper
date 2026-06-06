import 'dart:async';

import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:flutter/material.dart';

/// 宫格列数滞回，避免宽度抖动时 [GridView] 换列把滚动位置打回顶部。
int rbGridCrossForWidth({
  required double width,
  required double cellBaseWidth,
  required int currentCross,
  required double currentWidth,
}) {
  final step = cellBaseWidth;
  final target = (width / step).floor().clamp(3, 10);
  if (currentCross < 3 || currentWidth <= 0) {
    return target;
  }
  final lo = (currentCross - 0.55) * step;
  final hi = (currentCross + 0.55) * step;
  if (width >= lo && width <= hi) {
    return currentCross;
  }
  return target;
}

/// 媒体墙单次拉取条数（约 12 行；与 [RbMediaDfsScanner.fetchItems] 上限一致）。
int rbGridMediaFetchBatchSize(int crossAxisCount) {
  final cross = crossAxisCount.clamp(3, 10);
  return (cross * 12).clamp(48, 128);
}

/// 远程相册模块共用 UI。
class RbEmptyState extends StatelessWidget {
  const RbEmptyState({
    super.key,
    required this.icon,
    required this.title,
    this.subtitle,
    this.action,
  });

  final IconData icon;
  final String title;
  final String? subtitle;
  final Widget? action;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final subtitleStyle = Theme.of(context).textTheme.bodyMedium?.copyWith(color: cs.onSurfaceVariant, height: 1.45);
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 28),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 56, color: cs.primary.withValues(alpha: 0.55)),
          const SizedBox(height: 16),
          Text(title, textAlign: TextAlign.center, style: Theme.of(context).textTheme.titleMedium),
          if (subtitle != null) ...[
            const SizedBox(height: 8),
            Text(subtitle!, textAlign: TextAlign.center, style: subtitleStyle),
          ],
          if (action != null) ...[const SizedBox(height: 20), action!],
        ],
      ),
    );
  }
}

enum RbBannerTone { info, warning, loading }

/// 顶部 transient 提示（info / error / loading），统一样式；info/error 可自动消失。
mixin RbTopNoticeMixin<T extends StatefulWidget> on State<T> {
  String _rbInfoNotice = '';
  String _rbErrorNotice = '';
  String _rbLoadingNotice = '';
  Timer? _rbInfoClearTimer;
  Timer? _rbErrorClearTimer;

  @override
  void dispose() {
    _rbInfoClearTimer?.cancel();
    _rbErrorClearTimer?.cancel();
    super.dispose();
  }

  void rbShowInfoNotice(String msg, {Duration duration = const Duration(seconds: 3)}) {
    final text = msg.trim();
    if (text.isEmpty) {
      return;
    }
    _rbInfoClearTimer?.cancel();
    setState(() => _rbInfoNotice = text);
    _rbInfoClearTimer = Timer(duration, () {
      if (mounted && _rbInfoNotice == text) {
        setState(() => _rbInfoNotice = '');
      }
    });
  }

  void rbShowErrorNotice(String msg, {Duration duration = const Duration(seconds: 5)}) {
    final text = msg.trim();
    if (text.isEmpty) {
      return;
    }
    _rbErrorClearTimer?.cancel();
    setState(() => _rbErrorNotice = text);
    _rbErrorClearTimer = Timer(duration, () {
      if (mounted && _rbErrorNotice == text) {
        setState(() => _rbErrorNotice = '');
      }
    });
  }

  void rbClearErrorNotice() {
    _rbErrorClearTimer?.cancel();
    if (_rbErrorNotice.isEmpty) {
      return;
    }
    setState(() => _rbErrorNotice = '');
  }

  void rbSetLoadingNotice(String msg) {
    final text = msg.trim();
    if (_rbLoadingNotice == text) {
      return;
    }
    setState(() => _rbLoadingNotice = text);
  }

  void rbClearLoadingNotice() {
    if (_rbLoadingNotice.isEmpty) {
      return;
    }
    setState(() => _rbLoadingNotice = '');
  }

  String? get rbTopNoticeMessage {
    final info = _rbInfoNotice.trim();
    if (info.isNotEmpty) {
      return info;
    }
    final err = _rbErrorNotice.trim();
    if (err.isNotEmpty) {
      return err;
    }
    final load = _rbLoadingNotice.trim();
    if (load.isNotEmpty) {
      return load;
    }
    return null;
  }

  RbBannerTone get rbTopNoticeTone {
    if (_rbInfoNotice.trim().isNotEmpty) {
      return RbBannerTone.info;
    }
    if (_rbErrorNotice.trim().isNotEmpty) {
      return RbBannerTone.warning;
    }
    return RbBannerTone.loading;
  }

  bool get rbTopNoticeShowProgress =>
      _rbLoadingNotice.trim().isNotEmpty &&
      _rbInfoNotice.trim().isEmpty &&
      _rbErrorNotice.trim().isEmpty;
}

/// 一次性顶部提示（预览、下载等无 [RbTopNoticeMixin] 的页面）。
void rbFlashTopNotice(
  BuildContext context,
  String message, {
  RbBannerTone tone = RbBannerTone.info,
  Duration duration = const Duration(seconds: 2),
  bool showProgress = false,
}) {
  final text = message.trim();
  if (text.isEmpty) {
    return;
  }
  final overlay = Overlay.maybeOf(context, rootOverlay: true);
  if (overlay == null) {
    return;
  }
  late OverlayEntry entry;
  entry = OverlayEntry(
    builder: (ctx) => Positioned(
      top: 0,
      left: 0,
      right: 0,
      child: RbStatusOverlay(message: text, tone: tone, showProgress: showProgress),
    ),
  );
  overlay.insert(entry);
  Future<void>.delayed(duration, () {
    entry.remove();
  });
}

/// 顶部浮层提示：叠在 [child] 上，不占布局高度。
class RbStatusOverlayHost extends StatelessWidget {
  const RbStatusOverlayHost({
    super.key,
    required this.child,
    this.message,
    this.tone = RbBannerTone.loading,
    this.showProgress = true,
  });

  final Widget child;
  final String? message;
  final RbBannerTone tone;
  /// 为 false 时横幅仅文字（页面其它位置已有转圈时用）。
  final bool showProgress;

  @override
  Widget build(BuildContext context) {
    final msg = message?.trim() ?? '';
    if (msg.isEmpty) {
      return child;
    }
    return Stack(
      fit: StackFit.expand,
      children: [
        child,
        Positioned(top: 0, left: 0, right: 0, child: RbStatusOverlay(message: msg, tone: tone, showProgress: showProgress)),
      ],
    );
  }
}

/// 顶部横条提示（遮罩，不顶开下方内容）。
class RbStatusOverlay extends StatelessWidget {
  const RbStatusOverlay({
    super.key,
    required this.message,
    this.tone = RbBannerTone.loading,
    this.onDark = false,
    this.showProgress = true,
  });

  final String message;
  final RbBannerTone tone;
  final bool onDark;
  final bool showProgress;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final (bg, fg, icon) = onDark
        ? (Colors.black.withValues(alpha: 0.72), Colors.white, null)
        : switch (tone) {
            RbBannerTone.info => (cs.primaryContainer.withValues(alpha: 0.94), cs.onPrimaryContainer, Icons.info_outline),
            RbBannerTone.warning => (cs.errorContainer.withValues(alpha: 0.94), cs.onErrorContainer, Icons.cloud_off_outlined),
            RbBannerTone.loading => (cs.secondaryContainer.withValues(alpha: 0.94), cs.onSecondaryContainer, null),
          };
    return SafeArea(
      bottom: false,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(12, 6, 12, 0),
        child: Material(
          elevation: onDark ? 0 : 2,
          borderRadius: BorderRadius.circular(8),
          color: bg,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            child: Row(
              children: [
                if (showProgress && (tone == RbBannerTone.loading || onDark)) ...[
                  SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: fg)),
                  const SizedBox(width: 10),
                ] else if (icon != null) ...[
                  Icon(icon, size: 18, color: fg),
                  const SizedBox(width: 10),
                ],
                Expanded(child: Text(message, style: TextStyle(fontSize: 13, color: fg, height: 1.3))),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class RbStatusBanner extends StatelessWidget {
  const RbStatusBanner({super.key, required this.message, this.tone = RbBannerTone.info, this.trailing});

  final String message;
  final RbBannerTone tone;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final (bg, fg, icon) = switch (tone) {
      RbBannerTone.info => (cs.primaryContainer, cs.onPrimaryContainer, Icons.info_outline),
      RbBannerTone.warning => (cs.errorContainer, cs.onErrorContainer, Icons.cloud_off_outlined),
      RbBannerTone.loading => (cs.secondaryContainer, cs.onSecondaryContainer, null),
    };
    return Material(
      color: bg,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        child: Row(
          children: [
            if (tone == RbBannerTone.loading)
              SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: fg))
            else
              Icon(icon, size: 18, color: fg),
            const SizedBox(width: 10),
            Expanded(child: Text(message, style: TextStyle(fontSize: 13, color: fg, height: 1.3))),
            ?trailing,
          ],
        ),
      ),
    );
  }
}

/// 带尾部清除图标的输入框；有内容时显示 [Icons.clear]。
class RbClearTextField extends StatefulWidget {
  const RbClearTextField({
    super.key,
    required this.controller,
    this.decoration,
    this.readOnly = false,
    this.autocorrect,
    this.autofocus = false,
    this.textInputAction,
    this.onChanged,
    this.onSubmitted,
    this.extraSuffix,
  });

  final TextEditingController controller;
  final InputDecoration? decoration;
  final bool readOnly;
  final bool? autocorrect;
  final bool autofocus;
  final TextInputAction? textInputAction;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onSubmitted;
  final Widget? extraSuffix;

  @override
  State<RbClearTextField> createState() => _RbClearTextFieldState();
}

class _RbClearTextFieldState extends State<RbClearTextField> {
  @override
  void initState() {
    super.initState();
    widget.controller.addListener(_onText);
  }

  @override
  void dispose() {
    widget.controller.removeListener(_onText);
    super.dispose();
  }

  void _onText() {
    if (mounted) {
      setState(() {});
    }
  }

  void _clear() {
    widget.controller.clear();
    widget.onChanged?.call('');
  }

  @override
  Widget build(BuildContext context) {
    final showClear = !widget.readOnly && widget.controller.text.isNotEmpty;
    final suffix = (showClear || widget.extraSuffix != null)
        ? Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (showClear)
                IconButton(
                  icon: const Icon(Icons.clear, size: 20),
                  tooltip: '清除',
                  onPressed: _clear,
                ),
              if (widget.extraSuffix != null) widget.extraSuffix!,
            ],
          )
        : null;
    // TextField 内部 Scrollable 固定 restorationId=editable；多输入框 + 路由切换会污染恢复桶导致
    // restoreScrollOffset 把 bool 当成 double?。用 UnmanagedRestorationScope 隔离各输入框。
    return UnmanagedRestorationScope(
      child: TextField(
        controller: widget.controller,
        readOnly: widget.readOnly,
        autocorrect: widget.autocorrect ?? true,
        autofocus: widget.autofocus,
        textInputAction: widget.textInputAction,
        onChanged: widget.onChanged,
        onSubmitted: widget.onSubmitted,
        decoration: (widget.decoration ?? const InputDecoration()).copyWith(suffixIcon: suffix),
      ),
    );
  }
}

class RbConnectionTile extends StatelessWidget {
  const RbConnectionTile({
    super.key,
    required this.name,
    required this.subtitle,
    required this.onTap,
    this.remoteOs = RbRemoteOs.unknown,
    this.onEdit,
    this.onDelete,
  });

  final String name;
  final String subtitle;
  final RbRemoteOs remoteOs;
  final VoidCallback onTap;
  final VoidCallback? onEdit;
  final VoidCallback? onDelete;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
      elevation: 0,
      color: cs.surfaceContainerLow,
      child: InkWell(
        borderRadius: BorderRadius.circular(12),
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(14, 10, 6, 10),
          child: Row(
            children: [
              RbRemoteOsAvatar(os: remoteOs, radius: 22),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w600)),
                    const SizedBox(height: 2),
                    Text(subtitle, maxLines: 1, overflow: TextOverflow.ellipsis, style: Theme.of(context).textTheme.bodySmall),
                  ],
                ),
              ),
              if (onEdit != null)
                IconButton(icon: const Icon(Icons.edit_outlined, size: 20), tooltip: '编辑', onPressed: onEdit),
              if (onDelete != null)
                IconButton(icon: const Icon(Icons.delete_outline, size: 20), tooltip: '删除', onPressed: onDelete),
              Icon(Icons.chevron_right, color: cs.onSurfaceVariant),
            ],
          ),
        ),
      ),
    );
  }
}
