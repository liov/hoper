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
  final _sandboxCtrl = TextEditingController();
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
    _sandboxCtrl.dispose();
    _signalCtrl.dispose();
    super.dispose();
  }

  Future<void> _start() async {
    setState(() {
      _running = true;
      _status = '等待浏览端连接…（目录由浏览端填写）';
    });
    try {
      final sandbox = _sandboxCtrl.text.trim();
      await RbAgentSession.run(
        Uri.parse(_signalCtrl.text.trim()),
        _roomCtrl.text.trim(),
        sandbox: sandbox.isEmpty ? null : sandbox,
      );
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
            controller: _sandboxCtrl,
            decoration: const InputDecoration(
              labelText: '可选：本机沙箱根目录（限制浏览端只能访问其下路径）',
              border: OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 12),
          FilledButton(onPressed: _running ? null : _start, child: Text(_running ? '服务中…' : '启动共享')),
          if (_status.isNotEmpty)
            Padding(
              padding: const EdgeInsets.only(top: 12),
              child: Text(_status, style: Theme.of(context).textTheme.bodySmall),
            ),
          const SizedBox(height: 8),
          Text(
            '浏览端填写要看的目录；本机执行 rfv <房间码> 或设置 RB_ROOM 后执行 rfv。',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
          ),
        ],
      ),
    );
  }
}
