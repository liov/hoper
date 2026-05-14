import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/viewer_session.dart';
import 'package:app/remotebrowse/wire_session.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

class RemoteBrowseView extends StatefulWidget {
  const RemoteBrowseView({super.key});

  @override
  State<RemoteBrowseView> createState() => _RemoteBrowseViewState();
}

class _RemoteBrowseViewState extends State<RemoteBrowseView> {
  final _api = RemoteBrowseApi();
  final _roomCtrl = TextEditingController(text: 'demo');
  final _pathCtrl = TextEditingController();
  final _signalCtrl = TextEditingController(text: signalWsUrlFrom(RemoteBrowseApi.rbDebugBaseUrl));
  final _directCtrl = TextEditingController();
  final _entries = <RbFileEntry>[];
  var _loading = false;
  var _status = '';
  int _previewIndex = 0;
  RbWireSession? _wire;

  @override
  void dispose() {
    _roomCtrl.dispose();
    _pathCtrl.dispose();
    _signalCtrl.dispose();
    _directCtrl.dispose();
    _wire?.close();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _status = '连接文件端…';
    });
    try {
      final list = kIsWeb ? await _loadHttp() : await _loadP2P();
      setState(() {
        _entries
          ..clear()
          ..addAll(list);
        _previewIndex = 0;
        _status = kIsWeb ? 'H5 调试 HTTP · 房间 ${_roomCtrl.text.trim()}' : 'P2P · 房间 ${_roomCtrl.text.trim()}';
      });
    } catch (e) {
      setState(() => _status = '失败: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<List<RbFileEntry>> _loadHttp() => _api.listFiles(_pathCtrl.text.trim());

  Future<List<RbFileEntry>> _loadP2P() async {
    await _wire?.close();
    final direct = _parseDirect(_directCtrl.text.trim());
    if (direct != null && _roomCtrl.text.trim().isEmpty) {
      _wire = await RbViewerSession.connectDirect(direct.$1, direct.$2);
      return _wire!.listFiles(_pathCtrl.text.trim());
    }
    final ws = Uri.parse(_signalCtrl.text.trim());
    _wire = await RbViewerSession.connect(ws, _roomCtrl.text.trim(), directHost: direct?.$1, directPort: direct?.$2);
    return _wire!.listFiles(_pathCtrl.text.trim());
  }

  (String, int)? _parseDirect(String raw) {
    if (raw.isEmpty) {
      return null;
    }
  final hostPort = raw.split(':');
    if (hostPort.length == 1) {
      return (hostPort[0], RbViewerSession.directPort);
    }
    final port = int.tryParse(hostPort.last);
    if (port == null) {
      return null;
    }
    return (hostPort.sublist(0, hostPort.length - 1).join(':'), port);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('远程相册')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              children: [
                if (!kIsWeb)
                  TextField(
                    controller: _signalCtrl,
                    decoration: const InputDecoration(labelText: '文件端信令 wss://…/rb/signal', border: OutlineInputBorder()),
                  ),
                if (!kIsWeb)
                  TextField(
                    controller: _directCtrl,
                    decoration: const InputDecoration(labelText: '直连 IP（可选，如 192.168.1.5:19091）', border: OutlineInputBorder()),
                  ),
                if (!kIsWeb) const SizedBox(height: 8),
                TextField(controller: _roomCtrl, decoration: const InputDecoration(labelText: '房间码', border: OutlineInputBorder())),
                const SizedBox(height: 8),
                TextField(controller: _pathCtrl, decoration: const InputDecoration(labelText: '目录 path', border: OutlineInputBorder())),
                const SizedBox(height: 8),
                FilledButton(onPressed: _loading ? null : _load, child: Text(_loading ? '加载中' : '拉取列表')),
                if (_status.isNotEmpty) Padding(padding: const EdgeInsets.only(top: 8), child: Text(_status, style: Theme.of(context).textTheme.bodySmall)),
              ],
            ),
          ),
          if (_entries.isNotEmpty)
            Expanded(
              flex: 2,
              child: PageView.builder(
                itemCount: _entries.length,
                onPageChanged: (i) => setState(() => _previewIndex = i),
                itemBuilder: (ctx, i) => _ThumbPreview(wire: _wire, api: _api, entry: _entries[i]),
              ),
            ),
          Expanded(
            flex: 3,
            child: GridView.builder(
              padding: const EdgeInsets.all(8),
              gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 3, mainAxisSpacing: 6, crossAxisSpacing: 6),
              itemCount: _entries.length,
              itemBuilder: (ctx, i) {
                final e = _entries[i];
                return InkWell(
                  onTap: () => setState(() => _previewIndex = i),
                  child: DecoratedBox(
                    decoration: BoxDecoration(
                      border: Border.all(color: i == _previewIndex ? Theme.of(context).colorScheme.primary : Colors.black12),
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        Expanded(child: _ThumbPreview(wire: _wire, api: _api, entry: e, maxEdge: 128, fit: BoxFit.cover)),
                        Padding(padding: const EdgeInsets.all(4), child: Text(e.name, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontSize: 11))),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class _ThumbPreview extends StatelessWidget {
  const _ThumbPreview({required this.api, required this.entry, this.wire, this.maxEdge = 256, this.fit = BoxFit.contain});

  final RemoteBrowseApi api;
  final RbWireSession? wire;
  final RbFileEntry entry;
  final int maxEdge;
  final BoxFit fit;

  Future<Uint8List> _load() {
    final path = entry.id.isNotEmpty ? entry.id : entry.name;
    if (wire != null) {
      return wire!.fetchThumb(path, maxEdge: maxEdge);
    }
    return api.fetchThumb(path, hash: entry.thumbHash, maxEdge: maxEdge);
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: _load(),
      builder: (ctx, snap) {
        if (snap.connectionState != ConnectionState.done) {
          return const Center(child: CircularProgressIndicator(strokeWidth: 2));
        }
        final bytes = snap.data;
        if (bytes == null || bytes.isEmpty) {
          return const Center(child: Icon(Icons.broken_image_outlined));
        }
        return Image.memory(bytes, fit: fit, gaplessPlayback: true);
      },
    );
  }
}

class RemoteBrowseBinding extends Bindings {
  @override
  void dependencies() {}
}
