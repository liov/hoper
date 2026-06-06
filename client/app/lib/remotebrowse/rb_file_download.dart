import 'dart:io';
import 'dart:math' as math;

import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_agent_paths.dart';
import 'package:app/remotebrowse/rb_grpc_session.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:flutter/material.dart';
import 'package:path/path.dart' as p;

/// 单文件下载体积上限（2GB）。
const rbDownloadMaxBytes = 2 << 30;

/// 从远端分块读取并写入用户选择的路径；取消保存对话框返回 false。
Future<bool> rbDownloadRemoteFile({
  required RbGrpcSession wire,
  required String relPath,
  required String fileName,
  required int fileSize,
}) async {
  if (fileSize > rbDownloadMaxBytes) {
    throw StateError('文件超过 ${rbDownloadMaxBytes >> 30}GB，暂不支持下载');
  }
  final dir = await RbAgentPaths.pickDirectory();
  if (dir == null || dir.isEmpty) {
    return false;
  }
  final savePath = p.join(dir, fileName);
  final file = File(savePath);
  final sink = file.openWrite();
  var offset = 0;
  var total = fileSize > 0 ? fileSize : 0;
  try {
    while (true) {
      final remain = total > 0 ? total - offset : rbReadFileRangeMax;
      final want = math.min(rbReadFileRangeMax, remain > 0 ? remain : rbReadFileRangeMax);
      final chunk = await wire.readFileRange(relPath, offset: offset, length: want);
      if (total <= 0 && chunk.totalSize > 0) {
        total = chunk.totalSize;
        if (total > rbDownloadMaxBytes) {
          throw StateError('文件超过 ${rbDownloadMaxBytes >> 30}GB，暂不支持下载');
        }
      }
      if (chunk.bytes.isEmpty) {
        break;
      }
      sink.add(chunk.bytes);
      offset += chunk.bytes.length;
      if (total > 0 && offset >= total) {
        break;
      }
      if (chunk.bytes.length < want) {
        break;
      }
    }
  } finally {
    await sink.flush();
    await sink.close();
  }
  return true;
}

/// 文件信息面板「下载」：选路径、拉取、顶部提示反馈。
Future<void> rbRunFileInfoDownload(BuildContext context, {
  required RbGrpcSession wire,
  required String relPath,
  required RbFileEntry entry,
}) async {
  if (entry.isDirectory) {
    return;
  }
  rbFlashTopNotice(context, '正在下载 ${entry.name}…', tone: RbBannerTone.loading, duration: const Duration(seconds: 2), showProgress: true);
  try {
    final ok = await rbDownloadRemoteFile(
      wire: wire,
      relPath: relPath,
      fileName: entry.name,
      fileSize: entry.size,
    );
    if (!context.mounted) {
      return;
    }
    if (ok) {
      rbFlashTopNotice(context, '下载完成：${entry.name}');
    }
  } catch (e) {
    if (context.mounted) {
      rbFlashTopNotice(context, '下载失败：${rbUserMessage(e)}', tone: RbBannerTone.warning, duration: const Duration(seconds: 4));
    }
  }
}
