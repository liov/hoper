import 'dart:async';
import 'dart:math';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_icon.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_preview_nav.dart';
import 'package:app/remotebrowse/rb_memory_picture.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:flutter/material.dart';
import 'package:flutter/rendering.dart' show ScrollCacheExtent;

/// 预览底栏：当前目录可预览项横向缩略图，高亮当前文件。
class RbPreviewFolderStrip extends StatefulWidget {
  const RbPreviewFolderStrip({
    super.key,
    required this.nav,
    required this.thumbs,
    required this.onPick,
    this.thumbEdge = RbThumbLoader.sharedThumbEdge,
    this.compact = true,
  });

  final RbPreviewNav nav;
  final RbThumbLoader? thumbs;
  final void Function(int fileIndex) onPick;
  final int thumbEdge;
  final bool compact;

  @override
  State<RbPreviewFolderStrip> createState() => _RbPreviewFolderStripState();
}

class _RbPreviewFolderStripState extends State<RbPreviewFolderStrip> {
  static const _itemW = 48.0;
  static const _itemH = 48.0;
  static const _itemWNormal = 54.0;
  static const _itemHNormal = 54.0;
  static const _gap = 2.0;
  static const _hPad = 8.0;
  final _scroll = ScrollController();

  @override
  void initState() {
    super.initState();
    _scroll.addListener(_onScroll);
    widget.thumbs?.addListener(_onThumbs);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollToCurrent(animate: false);
      WidgetsBinding.instance.addPostFrameCallback((_) => _syncVisibleThumbs());
    });
  }

  @override
  void didUpdateWidget(covariant RbPreviewFolderStrip oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.nav.fileIndex != widget.nav.fileIndex) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _scrollToCurrent();
        _syncVisibleThumbs();
      });
    }
    if (oldWidget.nav.files != widget.nav.files) {
      WidgetsBinding.instance.addPostFrameCallback((_) => _syncVisibleThumbs());
    }
  }

  @override
  void dispose() {
    _scroll.removeListener(_onScroll);
    widget.thumbs?.removeListener(_onThumbs);
    _scroll.dispose();
    super.dispose();
  }

  void _onScroll() {
    _syncVisibleThumbs();
    if (_scroll.position.isScrollingNotifier.value) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) => _syncVisibleThumbs());
  }

  bool _onStripScrollNotification(ScrollNotification n) {
    if (n.metrics.axis != Axis.horizontal) {
      return false;
    }
    if (n is ScrollUpdateNotification || n is ScrollEndNotification) {
      _syncVisibleThumbs();
    }
    return false;
  }

  void _onThumbs() {
    if (!mounted) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        setState(() {});
      }
    });
  }

  List<int> get _indices {
    final out = <int>[];
    for (var i = 0; i < widget.nav.files.length; i++) {
      if (rbIsPreviewableEntry(widget.nav.files[i].entry)) {
        out.add(i);
      }
    }
    return out;
  }

  double get _itemStride => (widget.compact ? _itemW : _itemWNormal) + _gap;

  void _syncVisibleThumbs() {
    final loader = widget.thumbs;
    final indices = _indices;
    if (loader == null || indices.isEmpty) {
      return;
    }
    final files = widget.nav.files;
    final slice = <({String path, RbFileEntry entry})>[];
    if (!_scroll.hasClients) {
      final cur = indices.indexOf(widget.nav.fileIndex);
      final center = cur >= 0 ? cur : 0;
      final lo = max(0, center - RbThumbLoader.prefetchBehind);
      final hi = min(indices.length - 1, center + RbThumbLoader.prefetchAhead);
      for (var li = lo; li <= hi; li++) {
        final slot = files[indices[li]];
        if (rbFileSupportsAgentThumb(slot.entry.name)) {
          slice.add((path: slot.relPath, entry: slot.entry));
        }
      }
    } else {
      final offset = max(0.0, _scroll.offset - _hPad);
      final viewEnd = _scroll.offset + _scroll.position.viewportDimension + _hPad;
      final stride = _itemStride;
      final first = max(0, (offset / stride).floor() - RbThumbLoader.prefetchBehind);
      final last = min(indices.length - 1, (viewEnd / stride).ceil() + RbThumbLoader.prefetchAhead);
      for (var li = first; li <= last; li++) {
        final slot = files[indices[li]];
        if (rbFileSupportsAgentThumb(slot.entry.name)) {
          slice.add((path: slot.relPath, entry: slot.entry));
        }
      }
    }
    if (slice.isNotEmpty) {
      loader.prefetchStripVisibleRange(slice, maxEdge: widget.thumbEdge);
    }
  }

  void _scrollToCurrent({bool animate = true}) {
    final indices = _indices;
    if (!_scroll.hasClients || indices.isEmpty) {
      return;
    }
    final listIdx = indices.indexOf(widget.nav.fileIndex);
    final idx = listIdx >= 0 ? listIdx : 0;
    final viewport = _scroll.position.viewportDimension;
    final itemW = widget.compact ? _itemW : _itemWNormal;
    final target = idx * (itemW + _gap) - (viewport - itemW) / 2;
    final max = _scroll.position.maxScrollExtent;
    final off = target.clamp(0.0, max);
    if (animate) {
      unawaited(_scroll.animateTo(off, duration: const Duration(milliseconds: 220), curve: Curves.easeOut));
    } else {
      _scroll.jumpTo(off);
    }
  }

  @override
  Widget build(BuildContext context) {
    final indices = _indices;
    if (indices.length <= 1) {
      return const SizedBox.shrink();
    }
    final files = widget.nav.files;
    final stride = _itemStride;
    final list = NotificationListener<ScrollNotification>(
      onNotification: _onStripScrollNotification,
      child: ListView.separated(
        controller: _scroll,
        scrollCacheExtent: ScrollCacheExtent.pixels(stride * (RbThumbLoader.prefetchAhead + 2)),
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: _hPad),
        itemCount: indices.length,
        separatorBuilder: (_, _) => const SizedBox(width: _gap),
        itemBuilder: (ctx, li) {
          final fi = indices[li];
          final slot = files[fi];
          return _StripCell(
            loader: widget.thumbs,
            relPath: slot.relPath,
            entry: slot.entry,
            thumbEdge: widget.thumbEdge,
            selected: fi == widget.nav.fileIndex,
            compact: widget.compact,
            onTap: () => widget.onPick(fi),
          );
        },
      ),
    );
    final itemH = widget.compact ? _itemH : _itemHNormal;
    return SizedBox(height: itemH + 12, child: list);
  }
}

bool _stripNameAsThumb(RbFileEntry entry) {
  final k = rbPreviewKindForEntry(entry);
  return k != RbPreviewKind.image && k != RbPreviewKind.video && k != RbPreviewKind.motionPhoto;
}

Widget _stripFileNameThumb(String name, {required bool compact, required bool selected}) {
  return ColoredBox(
    color: selected ? const Color(0xFF3A3A3A) : const Color(0xFF2A2A2A),
    child: Padding(
      padding: const EdgeInsets.all(3),
      child: Center(
        child: Text(
          name,
          maxLines: 3,
          overflow: TextOverflow.ellipsis,
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: compact ? 8 : 9,
            height: 1.05,
            color: selected ? Colors.white : Colors.white70,
            fontWeight: FontWeight.w500,
          ),
        ),
      ),
    ),
  );
}

class _StripCell extends StatefulWidget {
  const _StripCell({
    required this.loader,
    required this.relPath,
    required this.entry,
    required this.thumbEdge,
    required this.selected,
    required this.compact,
    required this.onTap,
  });

  final RbThumbLoader? loader;
  final String relPath;
  final RbFileEntry entry;
  final int thumbEdge;
  final bool selected;
  final bool compact;
  final VoidCallback onTap;

  @override
  State<_StripCell> createState() => _StripCellState();
}

class _StripCellState extends State<_StripCell> {
  @override
  void initState() {
    super.initState();
    widget.loader?.addListener(_onLoader);
  }

  @override
  void didUpdateWidget(covariant _StripCell oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.loader != widget.loader) {
      oldWidget.loader?.removeListener(_onLoader);
      widget.loader?.addListener(_onLoader);
    }
  }

  @override
  void dispose() {
    widget.loader?.removeListener(_onLoader);
    super.dispose();
  }

  void _onLoader() {
    if (!mounted) {
      return;
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (mounted) {
        setState(() {});
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final entry = widget.entry;
    final loader = widget.loader;
    final isVideo = rbFileIsVideo(entry.name) || entry.isMotionPhoto;
    final nameThumb = _stripNameAsThumb(entry);
    final bytes = nameThumb ? null : loader?.bytesFor(widget.relPath);
    final loading = !nameThumb && loader != null && rbFileSupportsAgentThumb(entry.name) &&
        (bytes == null || bytes.isEmpty) && loader.isLoading(widget.relPath);
    final compact = widget.compact;
    final itemW = compact ? _RbPreviewFolderStripState._itemW : _RbPreviewFolderStripState._itemWNormal;
    final itemH = compact ? _RbPreviewFolderStripState._itemH : _RbPreviewFolderStripState._itemHNormal;
    final borderRadius = BorderRadius.circular(compact ? 5 : 6);
    final borderWidth = widget.selected ? (compact ? 2.0 : 2.5) : 1.0;
    final iconSize = compact ? 22.0 : 26.0;
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: widget.onTap,
        borderRadius: borderRadius,
        child: Container(
          width: itemW,
          height: itemH,
          decoration: BoxDecoration(
            borderRadius: borderRadius,
            border: Border.all(color: widget.selected ? Theme.of(context).colorScheme.primary : Colors.white24, width: borderWidth),
          ),
          clipBehavior: Clip.antiAlias,
          child: Stack(
            fit: StackFit.expand,
            children: [
              if (nameThumb)
                _stripFileNameThumb(entry.name, compact: compact, selected: widget.selected)
              else if (bytes != null && bytes.isNotEmpty)
                RbMemoryPicture(bytes: bytes, fileName: entry.name, fit: BoxFit.cover, gaplessPlayback: true)
              else if (loading)
                const ColoredBox(
                  color: Colors.white12,
                  child: Center(child: SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white54))),
                )
              else
                ColoredBox(
                  color: Colors.white12,
                  child: Center(child: Icon(rbFileTypeIcon(entry.name), color: Colors.white54, size: iconSize)),
                ),
              if (isVideo)
                Align(
                  alignment: Alignment.bottomRight,
                  child: Padding(
                    padding: const EdgeInsets.all(2),
                    child: Icon(Icons.play_circle_fill, size: compact ? 12 : 14, color: Colors.white),
                  ),
                ),
              if (entry.isMotionPhoto)
                const Align(
                  alignment: Alignment.topRight,
                  child: Padding(
                    padding: EdgeInsets.all(2),
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        color: Color(0x99000000),
                        borderRadius: BorderRadius.all(Radius.circular(2)),
                      ),
                      child: Padding(
                        padding: EdgeInsets.symmetric(horizontal: 3, vertical: 1),
                        child: Text('Live', style: TextStyle(fontSize: 7, color: Colors.white, fontWeight: FontWeight.w700)),
                      ),
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
