import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/agent_session.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_agent_paths.dart';
import 'package:app/remotebrowse/rb_user_message.dart';
import 'package:flutter/material.dart';

class RemoteBrowseAgentView extends StatefulWidget {
  const RemoteBrowseAgentView({super.key});

  @override
  State<RemoteBrowseAgentView> createState() => _RemoteBrowseAgentViewState();
}

class _RemoteBrowseAgentViewState extends State<RemoteBrowseAgentView> with AutomaticKeepAliveClientMixin {
  @override
  bool get wantKeepAlive => true;
  final _roomCtrl = TextEditingController();
  final _sandboxCtrl = TextEditingController();
  final _signalCtrl = TextEditingController();
  var _running = false;
  var _status = '';
  var _advancedOpen = false;
  String? _winDrive;

  @override
  void initState() {
    super.initState();
    _signalCtrl.text = '';
    _sandboxCtrl.text = '';
    final drives = RbAgentPaths.windowsDrives();
    if (drives.isNotEmpty) {
      _winDrive = drives.first;
    }
  }

  @override
  void dispose() {
    _roomCtrl.dispose();
    _sandboxCtrl.dispose();
    _signalCtrl.dispose();
    super.dispose();
  }

  void _applyWindowsDrive(String? drive) {
    if (drive == null || drive.isEmpty) {
      return;
    }
    setState(() {
      _winDrive = drive;
      _sandboxCtrl.text = drive;
    });
  }

  Future<void> _pickDirectory() async {
    final path = await RbAgentPaths.pickDirectory(initial: _sandboxCtrl.text.trim());
    if (path != null && path.isNotEmpty && mounted) {
      setState(() => _sandboxCtrl.text = path);
    }
  }

  Future<void> _start() async {
    if (_roomCtrl.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('请填写房间码')));
      return;
    }
    setState(() {
      _running = true;
      _status = '等待对方连接…';
    });
    try {
      final root = _sandboxCtrl.text.trim();
      await RbAgentSession.run(
        parseSignalWsUri(_signalCtrl.text.trim()),
        _roomCtrl.text.trim(),
        sandbox: root.isEmpty ? null : root,
      );
      if (mounted) {
        setState(() => _status = '共享已结束');
      }
    } catch (e) {
      if (mounted) {
        setState(() => _status = rbUserMessage(e));
      }
    } finally {
      if (mounted) {
        setState(() => _running = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    final drives = RbAgentPaths.windowsDrives();
    final hint = Theme.of(context).textTheme.bodySmall?.copyWith(color: Theme.of(context).colorScheme.onSurfaceVariant);
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text('在本机开启相册共享，对方在「浏览」页用相同房间码即可访问。', style: hint),
          const SizedBox(height: 16),
          if (_running)
            RbStatusBanner(
              message: _status.isEmpty ? '共享进行中…' : _status,
              tone: RbBannerTone.loading,
            ),
          if (_running) const SizedBox(height: 12),
          RbClearTextField(
            controller: _roomCtrl,
            readOnly: _running,
            decoration: const InputDecoration(labelText: '房间码', hintText: '告诉对方此房间码', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 12),
          if (drives.isNotEmpty) ...[
            InputDecorator(
              decoration: const InputDecoration(labelText: '磁盘（Windows）', border: OutlineInputBorder()),
              child: DropdownButtonHideUnderline(
                child: DropdownButton<String>(
                  isExpanded: true,
                  value: drives.contains(_winDrive) ? _winDrive : drives.first,
                  items: drives.map((d) => DropdownMenuItem(value: d, child: Text(d))).toList(),
                  onChanged: _running ? null : _applyWindowsDrive,
                ),
              ),
            ),
            const SizedBox(height: 12),
          ],
          RbClearTextField(
            controller: _sandboxCtrl,
            readOnly: _running,
            decoration: InputDecoration(
              labelText: '路径限制（可选）',
              hintText: '留空=可访问任意目录；仅默认打开用户目录',
              border: const OutlineInputBorder(),
            ),
            extraSuffix: IconButton(
              tooltip: '选择文件夹',
              onPressed: _running ? null : _pickDirectory,
              icon: const Icon(Icons.folder_open),
            ),
          ),
          const SizedBox(height: 12),
          ExpansionTile(
            initiallyExpanded: _advancedOpen,
            onExpansionChanged: _running ? null : (v) => setState(() => _advancedOpen = v),
            title: const Text('高级设置'),
            subtitle: Text('自定义服务器地址', style: hint),
            children: [
              RbClearTextField(
                controller: _signalCtrl,
                readOnly: _running,
                decoration: const InputDecoration(
                  labelText: '信令服务器',
                  hintText: 'host:8079 或 192.168.1.10',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 8),
            ],
          ),
          const SizedBox(height: 16),
          FilledButton.icon(
            onPressed: _running ? null : _start,
            icon: Icon(_running ? Icons.hourglass_top : Icons.play_arrow),
            label: Text(_running ? '共享进行中' : '开启共享'),
          ),
          if (!_running && _status.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text(_status, style: hint),
          ],
        ],
      ),
    );
  }
}
