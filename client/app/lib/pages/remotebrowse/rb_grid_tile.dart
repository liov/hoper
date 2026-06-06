import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_memory_picture.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:flutter/material.dart';

/// 目录宫格 / 媒体墙共用的文件格（缩略图走同一 [RbThumbLoader] 缓存）。
class RbGridMediaTile extends StatelessWidget {
  const RbGridMediaTile({
    super.key,
    required this.thumbs,
    required this.entry,
    required this.filePath,
    required this.thumbMaxEdge,
    this.selecting = false,
    this.selected = false,
    this.onTap,
  });

  final RbThumbLoader? thumbs;
  final RbFileEntry entry;
  final String filePath;
  final int thumbMaxEdge;
  final bool selecting;
  final bool selected;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final showName = !rbFileIsPhotoOrVideo(entry.name);
    return Material(
      color: selected ? cs.primaryContainer : cs.surfaceContainerHighest,
      child: InkWell(
        onTap: onTap,
        child: Stack(
          fit: StackFit.expand,
          children: [
            RbEntryThumb(loader: thumbs, entry: entry, filePath: filePath, maxEdge: thumbMaxEdge, fit: BoxFit.cover),
            if (entry.isMotionPhoto) const Positioned(right: 3, top: 3, child: RbGridLiveBadge()),
            if ((rbFileIsVideo(entry.name) || entry.isMotionPhoto) && entry.durationMs > 0)
              Positioned(left: 3, bottom: 3, child: RbGridVideoDurationLabel(ms: entry.durationMs)),
            if (showName) RbGridNameOverlay(name: entry.name),
            if (selecting)
              Positioned(
                top: 2,
                right: 2,
                child: Icon(
                  selected ? Icons.check_circle : Icons.circle_outlined,
                  size: 18,
                  color: selected ? cs.primary : (showName ? Colors.white : cs.onSurfaceVariant),
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class RbGridUpTile extends StatelessWidget {
  const RbGridUpTile({super.key, required this.onTap});

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Theme.of(context).colorScheme.surfaceContainerHighest,
      child: InkWell(
        onTap: onTap,
        child: const Stack(
          fit: StackFit.expand,
          children: [
            Center(child: Icon(Icons.arrow_upward, size: 24)),
            RbGridNameOverlay(name: '..', maxLines: 1),
          ],
        ),
      ),
    );
  }
}

class RbGridDirTile extends StatelessWidget {
  const RbGridDirTile({
    super.key,
    required this.entry,
    required this.onTap,
    this.selecting = false,
    this.selected = false,
  });

  final RbFileEntry entry;
  final VoidCallback onTap;
  final bool selecting;
  final bool selected;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Material(
      color: selected ? cs.primaryContainer : cs.surfaceContainerHighest,
      child: InkWell(
        onTap: onTap,
        child: Stack(
          fit: StackFit.expand,
          children: [
            Center(child: Icon(Icons.folder, size: 30, color: Colors.amber[800])),
            RbGridNameOverlay(name: entry.name),
            if (selecting)
              Positioned(
                top: 2,
                right: 2,
                child: Icon(
                  selected ? Icons.check_circle : Icons.circle_outlined,
                  size: 18,
                  color: selected ? cs.primary : cs.onSurfaceVariant,
                ),
              ),
          ],
        ),
      ),
    );
  }
}

class RbEntryThumb extends StatelessWidget {
  const RbEntryThumb({
    super.key,
    required this.loader,
    required this.entry,
    required this.filePath,
    required this.maxEdge,
    this.fit = BoxFit.cover,
  });

  final RbThumbLoader? loader;
  final RbFileEntry entry;
  final String filePath;
  final int maxEdge;
  final BoxFit fit;

  @override
  Widget build(BuildContext context) {
    if (loader == null) {
      return Center(child: Icon(rbFileTypeIcon(entry.name), size: maxEdge.clamp(22, 36).toDouble()));
    }
    return _RbThumbPreview(loader: loader!, entry: entry, filePath: filePath, maxEdge: maxEdge, fit: fit);
  }
}

class _RbThumbPreview extends StatefulWidget {
  const _RbThumbPreview({
    required this.loader,
    required this.entry,
    required this.filePath,
    required this.maxEdge,
    required this.fit,
  });

  final RbThumbLoader loader;
  final RbFileEntry entry;
  final String filePath;
  final int maxEdge;
  final BoxFit fit;

  @override
  State<_RbThumbPreview> createState() => _RbThumbPreviewState();
}

class _RbThumbPreviewState extends State<_RbThumbPreview> {
  @override
  void initState() {
    super.initState();
    _requestThumb();
  }

  @override
  void didUpdateWidget(covariant _RbThumbPreview oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.loader != widget.loader ||
        oldWidget.filePath != widget.filePath ||
        oldWidget.entry.id != widget.entry.id ||
        oldWidget.maxEdge != widget.maxEdge) {
      _requestThumb();
    }
  }

  void _requestThumb() {
    widget.loader.requestThumb(widget.filePath, widget.entry, maxEdge: widget.maxEdge, bypassPause: true);
  }

  @override
  Widget build(BuildContext context) {
    return ListenableBuilder(
      listenable: widget.loader,
      builder: (context, _) => _buildThumbBody(context),
    );
  }

  Widget _buildThumbBody(BuildContext context) {
    final entry = widget.entry;
    final iconSize = (widget.maxEdge * 0.42).clamp(28.0, 56.0);
    final icon = Center(
      child: rbFileTypeIconWidget(entry.name, size: iconSize, color: Theme.of(context).colorScheme.onSurfaceVariant),
    );
    if (!rbFileSupportsAgentThumb(entry.name)) {
      return icon;
    }
    final bytes = widget.loader.bytesFor(widget.filePath);
    if (bytes != null && bytes.isNotEmpty) {
      return RbMemoryPicture(bytes: bytes, fileName: entry.name, fit: widget.fit, gaplessPlayback: true);
    }
    if (widget.loader.isLoading(widget.filePath)) {
      return const Center(child: SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2)));
    }
    return icon;
  }
}

class RbGridLiveBadge extends StatelessWidget {
  const RbGridLiveBadge({super.key});

  @override
  Widget build(BuildContext context) {
    return DecoratedBox(
      decoration: BoxDecoration(color: Colors.black.withValues(alpha: 0.62), borderRadius: BorderRadius.circular(3)),
      child: const Padding(
        padding: EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        child: Text(
          'Live',
          style: TextStyle(fontSize: 9, height: 1.1, color: Colors.white, fontWeight: FontWeight.w700, letterSpacing: 0.3),
        ),
      ),
    );
  }
}

class RbGridVideoDurationLabel extends StatelessWidget {
  const RbGridVideoDurationLabel({super.key, required this.ms});

  final int ms;

  @override
  Widget build(BuildContext context) {
    final text = rbFormatDurationMs(ms);
    if (text.isEmpty) {
      return const SizedBox.shrink();
    }
    return DecoratedBox(
      decoration: BoxDecoration(color: Colors.black.withValues(alpha: 0.65), borderRadius: BorderRadius.circular(3)),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
        child: Text(text, style: const TextStyle(fontSize: 10, height: 1.1, color: Colors.white, fontWeight: FontWeight.w600)),
      ),
    );
  }
}

class RbGridNameOverlay extends StatelessWidget {
  const RbGridNameOverlay({super.key, required this.name, this.maxLines = 2});

  final String name;
  final int maxLines;

  @override
  Widget build(BuildContext context) {
    return Positioned(
      left: 0,
      right: 0,
      bottom: 0,
      child: DecoratedBox(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [Colors.transparent, Colors.black.withValues(alpha: 0.72)],
          ),
        ),
        child: Padding(
          padding: const EdgeInsets.fromLTRB(3, 12, 3, 3),
          child: Text(
            name,
            maxLines: maxLines,
            overflow: TextOverflow.ellipsis,
            style: const TextStyle(fontSize: 10, height: 1.15, color: Colors.white, fontWeight: FontWeight.w500),
          ),
        ),
      ),
    );
  }
}
