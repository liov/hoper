import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_media.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_preview_nav.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:path/path.dart' as p;

/// 文件信息面板数据（列表 / 预览共用）。
class RbFileInfoModel {
  RbFileInfoModel({
    required this.entry,
    required this.displayFilePath,
    required this.remoteCurrentDir,
    required this.remoteBrowseRoot,
    required this.previewKind,
    this.relPath = '',
    this.browsePath = '',
    this.fileIndex,
    this.fileCount,
  });

  final RbFileEntry entry;
  final String displayFilePath;
  final String remoteCurrentDir;
  final String remoteBrowseRoot;
  final RbPreviewKind previewKind;
  final String relPath;
  final String browsePath;
  final int? fileIndex;
  final int? fileCount;

  factory RbFileInfoModel.fromNav(RbPreviewNav nav, RbPreviewKind previewKind) {
    return RbFileInfoModel(
      entry: nav.current.entry,
      displayFilePath: nav.displayFilePath,
      remoteCurrentDir: nav.remoteCurrentDir,
      remoteBrowseRoot: nav.remoteBrowseRoot,
      previewKind: previewKind,
      relPath: nav.current.relPath,
      browsePath: nav.browsePath,
      fileIndex: nav.files.isEmpty ? null : nav.fileIndex,
      fileCount: nav.files.isEmpty ? null : nav.files.length,
    );
  }

  factory RbFileInfoModel.fromListEntry({
    required RbFileEntry entry,
    required String remoteCurrentDir,
    required String remoteBrowseRoot,
    required String relPath,
  }) {
    final cur = remoteCurrentDir.trim();
    final display = cur.isNotEmpty ? rbJoinRemotePath(cur, entry.name) : entry.name;
    return RbFileInfoModel(
      entry: entry,
      displayFilePath: display,
      remoteCurrentDir: cur,
      remoteBrowseRoot: remoteBrowseRoot.trim().isNotEmpty ? remoteBrowseRoot.trim() : cur,
      previewKind: rbPreviewKindForEntry(entry),
      relPath: relPath,
      browsePath: cur,
    );
  }

  List<({String name, String absPath, bool isFile})> get pathTreeSegments {
    final nodes = rbRemotePathSegments(displayFilePath);
    if (nodes.isEmpty) {
      return [(name: entry.name, absPath: displayFilePath, isFile: true)];
    }
    return [for (var i = 0; i < nodes.length; i++) (name: nodes[i].label, absPath: nodes[i].path, isFile: i == nodes.length - 1)];
  }

  String absoluteDirForTreeIndex(int index) {
    final segs = pathTreeSegments;
    if (index < 0 || index >= segs.length) {
      return remoteCurrentDir;
    }
    final seg = segs[index];
    if (seg.isFile) {
      return rbRemotePathParent(seg.absPath) ?? remoteCurrentDir;
    }
    return seg.absPath;
  }
}

/// 预览页顶部下滑面板。
class RbFileInfoTopPanel extends StatelessWidget {
  const RbFileInfoTopPanel({super.key, required this.model, required this.onClose, this.onNavigateToDir, this.onDownload});

  final RbFileInfoModel model;
  final VoidCallback onClose;
  final void Function(String absoluteDirPath)? onNavigateToDir;
  final Future<void> Function()? onDownload;

  static Future<void> show(
    BuildContext context, {
    required RbPreviewNav nav,
    required RbPreviewKind previewKind,
    void Function(String absoluteDirPath)? onNavigateToDir,
    Future<void> Function()? onDownload,
  }) {
    final model = RbFileInfoModel.fromNav(nav, previewKind);
    return showGeneralDialog<void>(
      context: context,
      barrierDismissible: true,
      barrierLabel: '关闭',
      barrierColor: Colors.black54,
      transitionDuration: const Duration(milliseconds: 240),
      pageBuilder: (context, animation, secondary) => const SizedBox.shrink(),
      transitionBuilder: (ctx, anim, secondary, child) {
        final curved = CurvedAnimation(parent: anim, curve: Curves.easeOutCubic, reverseCurve: Curves.easeInCubic);
        return SlideTransition(
          position: Tween(begin: const Offset(0, -1), end: Offset.zero).animate(curved),
          child: Align(
            alignment: Alignment.topCenter,
            child: Material(
              color: Colors.transparent,
              child: RbFileInfoTopPanel(
                model: model,
                onClose: () => Navigator.pop(ctx),
                onNavigateToDir: onNavigateToDir,
                onDownload: onDownload,
              ),
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final maxH = MediaQuery.sizeOf(context).height * 0.78;
    return ConstrainedBox(
      constraints: BoxConstraints(maxHeight: maxH),
      child: DecoratedBox(
        decoration: BoxDecoration(
          color: cs.surface,
          borderRadius: const BorderRadius.vertical(bottom: Radius.circular(16)),
          boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.18), blurRadius: 12, offset: const Offset(0, 4))],
        ),
        child: SafeArea(
          bottom: false,
          child: RbFileInfoBody(model: model, onClose: onClose, onNavigateToDir: onNavigateToDir, onDownload: onDownload),
        ),
      ),
    );
  }
}

/// 列表页不可预览文件：底部上滑面板。
class RbFileInfoBottomPanel extends StatelessWidget {
  const RbFileInfoBottomPanel({super.key, required this.model, this.onNavigateToDir, this.onDownload});

  final RbFileInfoModel model;
  final void Function(String absoluteDirPath)? onNavigateToDir;
  final Future<void> Function()? onDownload;

  static Future<void> showFromList(
    BuildContext context, {
    required RbFileEntry entry,
    required String remoteCurrentDir,
    required String remoteBrowseRoot,
    required String relPath,
    void Function(String absoluteDirPath)? onNavigateToDir,
    Future<void> Function()? onDownload,
  }) {
    final model = RbFileInfoModel.fromListEntry(
      entry: entry,
      remoteCurrentDir: remoteCurrentDir,
      remoteBrowseRoot: remoteBrowseRoot,
      relPath: relPath,
    );
    return showModalBottomSheet<void>(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => RbFileInfoBottomPanel(model: model, onNavigateToDir: onNavigateToDir, onDownload: onDownload),
    );
  }

  @override
  Widget build(BuildContext context) {
    final maxH = MediaQuery.sizeOf(context).height * 0.78;
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.viewInsetsOf(context).bottom),
      child: ConstrainedBox(
        constraints: BoxConstraints(maxHeight: maxH),
        child: Material(
          color: Theme.of(context).colorScheme.surface,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(16)),
          child: SafeArea(
            top: false,
            child: RbFileInfoBody(
              model: model,
              onClose: () => Navigator.pop(context),
              onNavigateToDir: onNavigateToDir,
              onDownload: onDownload,
            ),
          ),
        ),
      ),
    );
  }
}

class RbFileInfoBody extends StatelessWidget {
  const RbFileInfoBody({super.key, required this.model, required this.onClose, this.onNavigateToDir, this.onDownload});

  final RbFileInfoModel model;
  final VoidCallback onClose;
  final void Function(String absoluteDirPath)? onNavigateToDir;
  final Future<void> Function()? onDownload;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        const SizedBox(height: 8),
        Center(child: Container(width: 36, height: 4, decoration: BoxDecoration(color: cs.outlineVariant, borderRadius: BorderRadius.circular(2)))),
        Padding(
          padding: const EdgeInsets.fromLTRB(8, 4, 4, 0),
          child: Row(
            children: [
              Expanded(child: Text('文件信息', style: Theme.of(context).textTheme.titleMedium)),
              IconButton(tooltip: '关闭', onPressed: onClose, icon: const Icon(Icons.close)),
            ],
          ),
        ),
        Flexible(
          child: SingleChildScrollView(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                if (!rbIsPreviewableEntry(model.entry))
                  Padding(
                    padding: const EdgeInsets.only(bottom: 12),
                    child: Text('此文件类型不支持预览', style: Theme.of(context).textTheme.bodySmall?.copyWith(color: cs.onSurfaceVariant)),
                  ),
                if (onDownload != null && !model.entry.isDirectory)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 16),
                    child: _RbFileInfoDownloadButton(onDownload: onDownload!),
                  ),
                Text('路径', style: Theme.of(context).textTheme.labelLarge?.copyWith(color: cs.primary)),
                const SizedBox(height: 4),
                Text('点击文件夹进入该目录', style: Theme.of(context).textTheme.bodySmall?.copyWith(color: cs.onSurfaceVariant)),
                const SizedBox(height: 8),
                _PathTree(model: model, onNavigateToDir: onNavigateToDir, onClose: onClose),
                const SizedBox(height: 16),
                Text('详细信息', style: Theme.of(context).textTheme.labelLarge?.copyWith(color: cs.primary)),
                const SizedBox(height: 8),
                ..._detailTiles(context, model),
              ],
            ),
          ),
        ),
      ],
    );
  }

  List<Widget> _detailTiles(BuildContext context, RbFileInfoModel model) {
    final entry = model.entry;
    final ext = p.extension(entry.name);
    final root = model.remoteBrowseRoot.trim();
    final curDir = model.remoteCurrentDir.trim();
    final rows = <(String, String)>[
      ('文件名', entry.name),
      ('完整路径', model.displayFilePath),
      if (root.isNotEmpty) ('浏览根目录', root),
      if (curDir.isNotEmpty && curDir != root) ('当前目录', curDir),
      if (model.relPath.trim().isNotEmpty && model.relPath != model.displayFilePath) ('相对路径', model.relPath),
      if (model.browsePath.isNotEmpty && model.browsePath != '.' && model.browsePath != curDir) ('浏览路径', model.browsePath),
      ('大小', rbFormatFileSize(entry.size)),
      if (entry.mtimeUnixMs > 0) ('修改时间', DateFormat('yyyy-MM-dd HH:mm:ss').format(DateTime.fromMillisecondsSinceEpoch(entry.mtimeUnixMs))),
      ('类型', rbPreviewKindLabel(model.previewKind)),
      if (ext.isNotEmpty) ('扩展名', ext),
      if (entry.durationMs > 0) ('时长', rbFormatDurationMs(entry.durationMs)),
      if (entry.isMotionPhoto) ('动态照片', '是'),
      if (entry.motionCompanion.isNotEmpty) ('配对视频', entry.motionCompanion),
      if (entry.motionOffset > 0) ('内嵌视频偏移', '${entry.motionOffset} B'),
      if (entry.motionLength > 0) ('内嵌视频长度', rbFormatFileSize(entry.motionLength)),
      if (entry.thumbHash.isNotEmpty) ('缩略图哈希', entry.thumbHash),
      if (entry.id.isNotEmpty) ('文件 ID', entry.id),
      if (model.fileCount != null && model.fileIndex != null) ('目录内序号', '${model.fileIndex! + 1} / ${model.fileCount}'),
    ];
    return rows.map((r) => _DetailRow(label: r.$1, value: r.$2)).toList();
  }
}

class _PathTree extends StatelessWidget {
  const _PathTree({required this.model, required this.onClose, this.onNavigateToDir});

  final RbFileInfoModel model;
  final VoidCallback onClose;
  final void Function(String absoluteDirPath)? onNavigateToDir;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final segments = model.pathTreeSegments;
    return DecoratedBox(
      decoration: BoxDecoration(
        color: cs.surfaceContainerHighest.withValues(alpha: 0.45),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: cs.outlineVariant.withValues(alpha: 0.5)),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 4),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            for (var i = 0; i < segments.length; i++)
              _PathTreeRow(
                depth: i,
                name: segments[i].name,
                isFile: segments[i].isFile,
                isLast: i == segments.length - 1,
                tappable: onNavigateToDir != null,
                onTap: onNavigateToDir == null
                    ? null
                    : () {
                        onClose();
                        onNavigateToDir!(rbNormRemotePath(model.absoluteDirForTreeIndex(i)));
                      },
              ),
          ],
        ),
      ),
    );
  }
}

class _PathTreeRow extends StatelessWidget {
  const _PathTreeRow({
    required this.depth,
    required this.name,
    required this.isFile,
    required this.isLast,
    required this.tappable,
    this.onTap,
  });

  final int depth;
  final String name;
  final bool isFile;
  final bool isLast;
  final bool tappable;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final indent = 12.0 + depth * 18.0;
    final icon = isFile ? Icons.insert_drive_file_outlined : Icons.folder_outlined;
    final style = Theme.of(context).textTheme.bodyMedium?.copyWith(
          fontWeight: isLast ? FontWeight.w600 : FontWeight.normal,
          color: isLast ? cs.onSurface : (tappable ? cs.primary : cs.onSurfaceVariant),
        );
    final row = Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (depth > 0)
          Padding(
            padding: const EdgeInsets.only(top: 2, right: 6),
            child: Icon(Icons.subdirectory_arrow_right, size: 16, color: cs.outline),
          ),
        Icon(icon, size: 20, color: isLast ? cs.primary : cs.onSurfaceVariant),
        const SizedBox(width: 8),
        Expanded(child: Text(name, style: style)),
        if (tappable) Icon(Icons.chevron_right, size: 18, color: cs.onSurfaceVariant),
      ],
    );
    return Padding(
      padding: EdgeInsets.only(left: indent, right: 8, top: depth == 0 ? 0 : 4, bottom: isLast ? 0 : 4),
      child: tappable
          ? InkWell(
              borderRadius: BorderRadius.circular(8),
              onTap: onTap,
              child: Padding(padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 4), child: row),
            )
          : row,
    );
  }
}

class _RbFileInfoDownloadButton extends StatefulWidget {
  const _RbFileInfoDownloadButton({required this.onDownload});

  final Future<void> Function() onDownload;

  @override
  State<_RbFileInfoDownloadButton> createState() => _RbFileInfoDownloadButtonState();
}

class _RbFileInfoDownloadButtonState extends State<_RbFileInfoDownloadButton> {
  var _busy = false;

  Future<void> _tap() async {
    if (_busy) {
      return;
    }
    setState(() => _busy = true);
    try {
      await widget.onDownload();
    } finally {
      if (mounted) {
        setState(() => _busy = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return FilledButton.tonalIcon(
      onPressed: _busy ? null : _tap,
      icon: _busy
          ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
          : const Icon(Icons.download_outlined),
      label: Text(_busy ? '下载中…' : '下载到本机'),
    );
  }
}

class _DetailRow extends StatelessWidget {
  const _DetailRow({required this.label, required this.value});

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(width: 88, child: Text(label, style: Theme.of(context).textTheme.bodySmall?.copyWith(color: cs.onSurfaceVariant))),
          Expanded(child: SelectableText(value, style: Theme.of(context).textTheme.bodyMedium)),
        ],
      ),
    );
  }
}
