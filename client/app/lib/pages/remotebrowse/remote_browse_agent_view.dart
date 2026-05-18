import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/agent_session.dart';
import 'package:flutter/material.dart';

class RemoteBrowseAgentView extends StatefulWidget {
  const RemoteBrowseAgentView({super.key});

  @override
  State<RemoteBrowseAgentView> createState() => _RemoteBrowseAgentViewState();
}

class _RemoteBrowseAgentViewState extends State<RemoteBrowseAgentView> {
  final _roomCtrl = TextEditingController(text: 'demo');
  final _rootCtrl = TextEditingController();
  final _signalCtrl = TextEditingController();
  var _running = false;
  var _status = '';

  @override
  void initState() {
    super.initState();
    _signalCtrl.text = defaultSignalWs();
  }

  @override
  void dispose() {
    _roomCtrl.dispose();
    _rootCtrl.dispose();
    _signalCtrl.dispose();
    super.dispose();
  }

  Future<void> _start() async {
    final root = _rootCtrl.text.trim();
    if (root.isEmpty) {
      setState(() => _status = '请填写本机共享目录');
      return;
    }
    setState(() {
      _running = true;
      _status = '等待浏览端连接…';
    });
    try {
      await RbAgentSession.run(Uri.parse(_signalCtrl.text.trim()), _roomCtrl.text.trim(), root);
      if (mounted) {
        setState(() => _status = '会话已结束');
      }
    } catch (e) {
      if (mounted) {
        setState(() => _status = '失败: $e');
      }
    } finally {
      if (mounted) {
        setState(() => _running = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          TextField(
            controller: _signalCtrl,
            decoration: const InputDecoration(labelText: '信令 wss://…/rb/signal', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 8),
          TextField(controller: _roomCtrl, decoration: const InputDecoration(labelText: '房间码', border: OutlineInputBorder())),
          const SizedBox(height: 8),
          TextField(
            controller: _rootCtrl,
            decoration: const InputDecoration(labelText: '本机目录（绝对路径）', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 12),
          FilledButton(onPressed: _running ? null : _start, child: Text(_running ? '服务中…' : '启动 Agent')),
          if (_status.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(top: 12),
              child: Text(_status, style: Theme.of(context).textTheme.bodySmall),
            ),
          const SizedBox(height: 8),
          Text(
            '共享端由 Rust（rfv）负责打洞与缩略图；连接顺序：直连 → ICE → 中继。浏览端使用相同房间码。',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
          ),
        ],
      ),
    );
  }
}
