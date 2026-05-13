import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/relay.dart';
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
  final _entries = <RbFileEntry>[];
  var _loading = false;
  var _status = '';
  var _useRelay = !kIsWeb;
  int _previewIndex = 0;
  var _thumbEdge = 256;
  RbRelaySession? _relay;

  int get thumbEdge => _useRelay && !kIsWeb ? 128 : _thumbEdge;

  @override
  void dispose() {
    _roomCtrl.dispose();
    _pathCtrl.dispose();
    _relay?.close();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() {
      _loading = true;
      _status = '加载中…';
    });
    try {
      final list = _useRelay && !kIsWeb ? await _loadRelay() : await _loadHttp();
      final health = await _api.health();
      setState(() {
        _entries
          ..clear()
          ..addAll(list);
        _previewIndex = 0;
        final mode = _useRelay && !kIsWeb ? '中继' : 'HTTP';
        _status = '$mode · 房间 ${_roomCtrl.text.trim()} · relay ${health['relayTcp'] ?? health['relay_tcp'] ?? '-'}';
      });
    } catch (e) {
      setState(() => _status = '失败: $e');
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<List<RbFileEntry>> _loadHttp() => _api.listFiles(_pathCtrl.text.trim());

  Future<List<RbFileEntry>> _loadRelay() async {
    await _relay?.close();
    _relay = await RbRelaySession.openViewer(Uri.parse(_api.signalWsUrl()), _roomCtrl.text.trim());
    return _relay!.listFiles(_pathCtrl.text.trim());
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
                TextField(controller: _roomCtrl, decoration: const InputDecoration(labelText: '房间码', border: OutlineInputBorder())),
                const SizedBox(height: 8),
                TextField(controller: _pathCtrl, decoration: const InputDecoration(labelText: '目录 path', border: OutlineInputBorder())),
                if (!kIsWeb)
                  SwitchListTile(
                    contentPadding: EdgeInsets.zero,
                    title: const Text('经中继 P2P'),
                    value: _useRelay,
                    onChanged: _loading ? null : (v) => setState(() => _useRelay = v),
                  ),
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
                itemBuilder: (ctx, i) => _ThumbPreview(api: _api, relay: _useRelay ? _relay : null, entry: _entries[i], maxEdge: thumbEdge),
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
                        Expanded(child: _ThumbPreview(api: _api, relay: _useRelay ? _relay : null, entry: e, maxEdge: thumbEdge, fit: BoxFit.cover)),
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
  const _ThumbPreview({required this.api, required this.entry, this.relay, this.maxEdge = 256, this.fit = BoxFit.contain});

  final RemoteBrowseApi api;
  final RbRelaySession? relay;
  final RbFileEntry entry;
  final int maxEdge;
  final BoxFit fit;

  Future<Uint8List> _load() {
    final path = entry.id.isNotEmpty ? entry.id : entry.name;
    if (relay != null) {
      return relay!.fetchThumb(path, maxEdge: maxEdge);
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
