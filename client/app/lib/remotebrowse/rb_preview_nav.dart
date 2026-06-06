import 'dart:math';

import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:app/remotebrowse/rb_thumb_loader.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:path/path.dart' as p;

/// 预览页路径树跳转：回到列表并打开该目录。
class RbPreviewNavigateToDir {
  const RbPreviewNavigateToDir(this.absoluteDirPath);

  final String absoluteDirPath;
}

class RbPreviewFileSlot {
  RbPreviewFileSlot({required this.entry, required this.relPath});

  final RbFileEntry entry;
  final String relPath;
}

/// 预览页导航：上下换文件、左右换同级目录。
class RbPreviewNav {
  RbPreviewNav({
    required this.wire,
    required this.thumbs,
    required this.browsePath,
    required this.parentBrowsePath,
    required this.files,
    required this.fileIndex,
    required this.entryRelPath,
    this.remoteCurrentDir = '',
    this.remoteBrowseRoot = '',
  });

  RbGrpcSession wire;
  RbThumbLoader? thumbs;
  String browsePath;
  /// 当前所在目录的绝对路径（Agent `resolved_root_path`）。
  String remoteCurrentDir;
  /// 本次连接浏览的根目录绝对路径。
  String remoteBrowseRoot;
  final String parentBrowsePath;
  List<RbPreviewFileSlot> files;
  int fileIndex;
  final String Function(RbFileEntry) entryRelPath;

  List<String>? _siblingDirNames;
  var _siblingsLoading = false;

  factory RbPreviewNav.initial({
    required RbGrpcSession wire,
    required RbThumbLoader? thumbs,
    required String browsePath,
    required String parentBrowsePath,
    required List<RbFileEntry> fileEntries,
    required int fileIndex,
    required String Function(RbFileEntry) entryRelPath,
    String remoteCurrentDir = '',
    String remoteBrowseRoot = '',
  }) {
    final files = fileEntries.map((e) => RbPreviewFileSlot(entry: e, relPath: entryRelPath(e))).toList();
    final idx = fileIndex.clamp(0, files.isEmpty ? 0 : files.length - 1);
    final curDir = remoteCurrentDir.trim();
    final root = remoteBrowseRoot.trim().isNotEmpty ? remoteBrowseRoot.trim() : curDir;
    return RbPreviewNav(
      wire: wire,
      thumbs: thumbs,
      browsePath: browsePath,
      parentBrowsePath: parentBrowsePath,
      files: files,
      fileIndex: idx,
      entryRelPath: entryRelPath,
      remoteCurrentDir: curDir,
      remoteBrowseRoot: root,
    );
  }

  factory RbPreviewNav.fromSlots({
    required RbGrpcSession wire,
    required RbThumbLoader? thumbs,
    required List<RbPreviewFileSlot> files,
    required int fileIndex,
    String remoteBrowseRoot = '',
  }) {
    final idx = fileIndex.clamp(0, files.isEmpty ? 0 : files.length - 1);
    final root = remoteBrowseRoot.trim();
    return RbPreviewNav(
      wire: wire,
      thumbs: thumbs,
      browsePath: '',
      parentBrowsePath: '',
      files: files,
      fileIndex: idx,
      entryRelPath: (_) => '',
      remoteCurrentDir: '',
      remoteBrowseRoot: root.isNotEmpty ? root : '',
    );
  }

  /// 完整路径按层级拆分，最后一级为当前文件。
  List<({String name, String absPath, bool isFile})> get pathTreeSegments {
    final nodes = rbRemotePathSegments(displayFilePath);
    if (nodes.isEmpty) {
      return [(name: current.entry.name, absPath: displayFilePath, isFile: true)];
    }
    return [for (var i = 0; i < nodes.length; i++) (name: nodes[i].label, absPath: nodes[i].path, isFile: i == nodes.length - 1)];
  }

  /// 路径树点击后应打开的目录绝对路径。
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

  /// 远端主机上的文件绝对路径（Agent resolved 目录 + 文件名）。
  String get displayFilePath {
    final slot = current;
    final dir = remoteCurrentDir.trim();
    if (dir.isNotEmpty) {
      return rbJoinRemotePath(dir, slot.entry.name);
    }
    final rel = slot.relPath.trim();
    return rel.isNotEmpty ? rel : slot.entry.name;
  }

  void _clampFileIndex() {
    if (files.isEmpty) {
      fileIndex = 0;
      return;
    }
    if (fileIndex < 0 || fileIndex >= files.length) {
      fileIndex = fileIndex.clamp(0, files.length - 1);
    }
  }

  RbPreviewFileSlot get current {
    _clampFileIndex();
    return files[fileIndex];
  }

  String get currentDirName {
    if (browsePath.isEmpty || browsePath == '.') {
      return '.';
    }
    return p.basename(browsePath);
  }

  bool get canPrevFile => _peekPreviewableStep(-1) != null;

  bool get canNextFile => _peekPreviewableStep(1) != null;

  int? _peekPreviewableStep(int delta) {
    if (files.isEmpty) {
      return null;
    }
    _clampFileIndex();
    final start = fileIndex;
    if (delta > 0) {
      for (var i = start + 1; i < files.length; i++) {
        if (rbIsPreviewableEntry(files[i].entry)) {
          return i;
        }
      }
      return null;
    }
    for (var i = start - 1; i >= 0; i--) {
      if (rbIsPreviewableEntry(files[i].entry)) {
        return i;
      }
    }
    return null;
  }

  List<({String path, RbFileEntry entry})> filesAround(int index, {int before = 3, int after = 8}) {
    if (files.isEmpty) {
      return [];
    }
    final lo = max(0, index - before);
    final hi = min(files.length - 1, index + after);
    return [for (var i = lo; i <= hi; i++) (path: files[i].relPath, entry: files[i].entry)];
  }

  bool get canSiblingNav => parentBrowsePath != browsePath;

  int? get siblingIndex {
    final names = _siblingDirNames;
    if (names == null || names.isEmpty) {
      return null;
    }
    final name = currentDirName;
    final i = names.indexOf(name);
    return i >= 0 ? i : null;
  }

  Future<String?> _ensureSiblings() async {
    if (!canSiblingNav) {
      return '已在根目录，无同级目录';
    }
    if (_siblingDirNames != null) {
      return null;
    }
    if (_siblingsLoading) {
      return '正在加载同级目录…';
    }
    _siblingsLoading = true;
    try {
      final list = await wire.listFiles(parentBrowsePath);
      _siblingDirNames = list.entries.where((e) => e.isDirectory).map((e) => e.name).toList()..sort();
      if (_siblingDirNames!.isEmpty) {
        return '父目录下没有其它文件夹';
      }
      return null;
    } catch (e) {
      return e.toString();
    } finally {
      _siblingsLoading = false;
    }
  }

  /// 删除当前项后是否还有文件；无文件时 [fileIndex] 为 0。
  bool removeCurrentFile() {
    if (files.isEmpty) {
      return false;
    }
    _clampFileIndex();
    files.removeAt(fileIndex);
    if (files.isEmpty) {
      fileIndex = 0;
      return false;
    }
    if (fileIndex >= files.length) {
      fileIndex = files.length - 1;
    }
    return true;
  }

  Future<(bool ok, String? msg)> prevFile() async => _stepPreviewable(-1);

  Future<(bool ok, String? msg)> nextFile() async => _stepPreviewable(1);

  /// 当前不可预览时，先向后再向前找最近的可预览项。
  Future<(bool ok, String? msg)> focusPreviewableFile() async {
    if (files.isEmpty) {
      return (false, '没有文件');
    }
    _clampFileIndex();
    if (rbIsPreviewableEntry(files[fileIndex].entry)) {
      return (true, null);
    }
    final start = fileIndex;
    for (var i = start + 1; i < files.length; i++) {
      if (rbIsPreviewableEntry(files[i].entry)) {
        fileIndex = i;
        return (true, null);
      }
    }
    for (var i = start - 1; i >= 0; i--) {
      if (rbIsPreviewableEntry(files[i].entry)) {
        fileIndex = i;
        return (true, null);
      }
    }
    return (false, '目录中没有可预览文件');
  }

  Future<(bool ok, String? msg)> _stepPreviewable(int delta) async {
    if (files.isEmpty) {
      return (false, '没有文件');
    }
    _clampFileIndex();
    final start = fileIndex;
    if (delta > 0) {
      for (var i = start + 1; i < files.length; i++) {
        if (rbIsPreviewableEntry(files[i].entry)) {
          fileIndex = i;
          return (true, null);
        }
      }
      return (false, '已是最后一个可预览文件');
    }
    for (var i = start - 1; i >= 0; i--) {
      if (rbIsPreviewableEntry(files[i].entry)) {
        fileIndex = i;
        return (true, null);
      }
    }
    return (false, '已是第一个可预览文件');
  }

  int? _firstPreviewableIndex() {
    for (var i = 0; i < files.length; i++) {
      if (rbIsPreviewableEntry(files[i].entry)) {
        return i;
      }
    }
    return null;
  }

  Future<(bool ok, String? msg)> prevSiblingDir() => _jumpSibling(-1);

  Future<(bool ok, String? msg)> nextSiblingDir() => _jumpSibling(1);

  Future<(bool ok, String? msg)> _jumpSibling(int delta) async {
    final err = await _ensureSiblings();
    if (err != null) {
      return (false, err);
    }
    final names = _siblingDirNames!;
    var idx = siblingIndex ?? 0;
    final next = idx + delta;
    if (next < 0) {
      return (false, '已是第一个同级目录');
    }
    if (next >= names.length) {
      return (false, '已是最后一个同级目录');
    }
    final targetDir = names[next];
    browsePath = parentBrowsePath == '.' ? targetDir : p.join(parentBrowsePath, targetDir);
    _siblingDirNames = names;
    try {
      final list = await wire.listFiles(browsePath);
      remoteCurrentDir = rbNormResolvedRootPath(list.resolvedRootPath);
      final entries = list.entries.where((e) => !e.isDirectory).toList();
      rbSortEntriesByMtimeDesc(entries);
      if (entries.isEmpty) {
        return (false, '「$targetDir」中没有可预览文件');
      }
      files = entries.map((e) => RbPreviewFileSlot(entry: e, relPath: _relInBrowse(e.name))).toList();
      final pick = _firstPreviewableIndex();
      if (pick == null) {
        return (false, '「$targetDir」中没有可预览文件');
      }
      fileIndex = pick;
      _clampFileIndex();
      return (true, null);
    } catch (e) {
      return (false, e.toString());
    }
  }

  String _relInBrowse(String name) {
    if (browsePath.isEmpty || browsePath == '.') {
      return name;
    }
    return rbJoinRemotePath(browsePath, name);
  }

  static String parentOfBrowsePath(String browsePath) {
    if (browsePath.isEmpty || browsePath == '.') {
      return '.';
    }
    return rbRemotePathParent(rbNormRemotePath(browsePath)) ?? '.';
  }
}
