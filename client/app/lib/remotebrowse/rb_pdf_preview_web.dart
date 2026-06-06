import 'package:app/remotebrowse/rb_file_preview.dart';
import 'package:flutter/material.dart';
import 'package:webview_flutter/webview_flutter.dart';

class RbPdfPreview extends StatefulWidget {
  const RbPdfPreview({super.key, required this.mediaUri, this.maxBytes = rbPreviewPdfMaxBytes});

  final Uri mediaUri;
  final int maxBytes;

  @override
  State<RbPdfPreview> createState() => _RbPdfPreviewState();
}

class _RbPdfPreviewState extends State<RbPdfPreview> {
  WebViewController? _controller;
  var _failed = false;
  String? _failMsg;

  @override
  void initState() {
    super.initState();
    final c = WebViewController()
      ..setJavaScriptMode(JavaScriptMode.unrestricted)
      ..setBackgroundColor(Colors.white)
      ..setNavigationDelegate(
        NavigationDelegate(
          onWebResourceError: (e) {
            if (!mounted) {
              return;
            }
            setState(() {
              _failed = true;
              _failMsg = e.description;
            });
          },
        ),
      )
      ..loadRequest(widget.mediaUri);
    _controller = c;
  }

  @override
  Widget build(BuildContext context) {
    if (_failed) {
      return Center(child: Text(_failMsg ?? '无法预览 PDF', textAlign: TextAlign.center));
    }
    final c = _controller;
    if (c == null) {
      return const Center(child: CircularProgressIndicator());
    }
    return SizedBox.expand(child: WebViewWidget(controller: c));
  }
}
