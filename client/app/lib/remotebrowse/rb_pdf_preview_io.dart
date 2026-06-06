import 'dart:async';
import 'dart:io';

import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:flutter/material.dart';
import 'package:flutter_pdfview/flutter_pdfview.dart';
import 'package:open_file/open_file.dart';
import 'package:path_provider/path_provider.dart';

/// PDF：HTTP 下载后用系统 PdfRenderer 渲染（避免 WebView/pdfrx PDFium native assets 问题）。
class RbPdfPreview extends StatefulWidget {
  const RbPdfPreview({super.key, required this.mediaUri, this.maxBytes = rbPreviewPdfMaxBytes});

  final Uri mediaUri;
  final int maxBytes;

  @override
  State<RbPdfPreview> createState() => _RbPdfPreviewState();
}

class _RbPdfPreviewState extends State<RbPdfPreview> {
  var _failed = false;
  String? _failMsg;
  String? _tempPath;
  var _phase = '正在连接…';
  var _pageCount = 0;
  var _pageIndex = 0;

  @override
  void initState() {
    super.initState();
    unawaited(_prepare());
  }

  @override
  void dispose() {
    _removeTemp();
    super.dispose();
  }

  Future<void> _prepare() async {
    try {
      if (mounted) {
        setState(() => _phase = '正在下载…');
      }
      final path = await _downloadToTemp(widget.mediaUri, widget.maxBytes);
      if (!mounted) {
        return;
      }
      setState(() {
        _tempPath = path;
        _phase = '';
      });
    } catch (e) {
      if (mounted) {
        setState(() {
          _failed = true;
          _failMsg = e.toString();
        });
      }
    }
  }

  Future<String> _downloadToTemp(Uri uri, int maxBytes) async {
    final client = HttpClient();
    try {
      final headReq = await client.headUrl(uri).timeout(const Duration(seconds: 30));
      final headResp = await headReq.close().timeout(const Duration(seconds: 30));
      final total = headResp.contentLength;
      if (total > 0 && total > maxBytes) {
        throw StateError('PDF 超过预览上限 ${maxBytes >> 20}MB');
      }
      final getReq = await client.getUrl(uri).timeout(const Duration(seconds: 30));
      final resp = await getReq.close().timeout(const Duration(minutes: 5));
      if (resp.statusCode != HttpStatus.ok && resp.statusCode != HttpStatus.partialContent) {
        throw HttpException('HTTP ${resp.statusCode}', uri: uri);
      }
      final dir = await getTemporaryDirectory();
      final path = '${dir.path}/rb_preview_${DateTime.now().millisecondsSinceEpoch}.pdf';
      final f = File(path);
      final sink = f.openWrite();
      var done = 0;
      await for (final chunk in resp) {
        done += chunk.length;
        if (done > maxBytes) {
          await sink.close();
          await f.delete();
          throw StateError('PDF 超过预览上限 ${maxBytes >> 20}MB');
        }
        sink.add(chunk);
        if (mounted && total > 0) {
          final pct = (100 * done / total).clamp(0, 100).toStringAsFixed(0);
          setState(() => _phase = '正在下载 $pct%');
        }
      }
      await sink.flush();
      await sink.close();
      return path;
    } finally {
      client.close(force: true);
    }
  }

  void _removeTemp() {
    final path = _tempPath;
    _tempPath = null;
    if (path == null) {
      return;
    }
    try {
      File(path).deleteSync();
    } catch (_) {}
  }

  Future<void> _openExternally() async {
    final path = _tempPath;
    if (path == null) {
      return;
    }
    final r = await OpenFile.open(path, type: 'application/pdf');
    if (!mounted || r.type == ResultType.done) {
      return;
    }
    setState(() => _failMsg = '${_failMsg ?? ''}\n${r.message}');
  }

  @override
  Widget build(BuildContext context) {
    if (_failed) {
      return _buildErrorPanel();
    }
    final path = _tempPath;
    if (path == null) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const CircularProgressIndicator(),
            if (_phase.isNotEmpty) ...[
              const SizedBox(height: 12),
              Text(_phase, style: Theme.of(context).textTheme.bodyMedium),
            ],
          ],
        ),
      );
    }
    return Stack(
      fit: StackFit.expand,
      children: [
        PDFView(
          filePath: path,
          enableSwipe: true,
          swipeHorizontal: false,
          autoSpacing: true,
          pageFling: true,
          onRender: (pages) {
            if (!mounted) {
              return;
            }
            setState(() => _pageCount = pages ?? 0);
          },
          onError: (err) {
            if (!mounted) {
              return;
            }
            setState(() {
              _failed = true;
              _failMsg = err;
            });
          },
          onPageError: (page, err) {
            if (!mounted) {
              return;
            }
            setState(() {
              _failed = true;
              _failMsg = '第 $page 页: $err';
            });
          },
          onPageChanged: (page, total) {
            if (!mounted) {
              return;
            }
            setState(() {
              _pageIndex = page ?? 0;
              if (total != null && total > 0) {
                _pageCount = total;
              }
            });
          },
        ),
        if (_pageCount > 0)
          Positioned(
            left: 0,
            right: 0,
            bottom: 8,
            child: Center(
              child: DecoratedBox(
                decoration: BoxDecoration(color: Colors.black54, borderRadius: BorderRadius.circular(16)),
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                  child: Text(
                    '${_pageIndex + 1} / $_pageCount',
                    style: const TextStyle(color: Colors.white, fontSize: 12),
                  ),
                ),
              ),
            ),
          ),
      ],
    );
  }

  Widget _buildErrorPanel() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(_failMsg ?? '无法预览 PDF', textAlign: TextAlign.center),
            if (_tempPath != null) ...[
              const SizedBox(height: 16),
              FilledButton.icon(
                onPressed: () => unawaited(_openExternally()),
                icon: const Icon(Icons.open_in_new),
                label: const Text('用系统应用打开'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
