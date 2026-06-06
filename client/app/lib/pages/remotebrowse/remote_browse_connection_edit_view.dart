import 'dart:async';

import 'package:app/remotebrowse/rb_connect_endpoint.dart';
import 'package:app/remotebrowse/rb_connection_profile.dart';
import 'package:app/remotebrowse/rb_connection_store.dart';
import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/material.dart';

class RemoteBrowseConnectionEditView extends StatefulWidget {
  const RemoteBrowseConnectionEditView({super.key, this.profile});

  final RbConnectionProfile? profile;

  bool get isEdit => profile != null;

  @override
  State<RemoteBrowseConnectionEditView> createState() => _RemoteBrowseConnectionEditViewState();
}

class _RemoteBrowseConnectionEditViewState extends State<RemoteBrowseConnectionEditView> {
  late final TextEditingController _nameCtrl;
  late final TextEditingController _endpointCtrl;
  late final TextEditingController _roomCtrl;
  late final TextEditingController _pathCtrl;
  var _saving = false;
  RbConnectTargetKind? _endpointKind;
  var _probing = false;
  Timer? _probeDebounce;
  var _probeGen = 0;

  @override
  void initState() {
    super.initState();
    final p = widget.profile;
    _nameCtrl = TextEditingController(text: p?.name ?? '');
    final ep = p == null ? '' : RbConnectTarget.endpointDisplay(p);
    _endpointCtrl = TextEditingController(text: ep);
    _roomCtrl = TextEditingController(text: p?.room ?? (p == null ? '' : 'demo'));
    _pathCtrl = TextEditingController(text: p?.path ?? '');
    _endpointCtrl.addListener(_onEndpointChanged);
    _scheduleEndpointProbe();
  }

  @override
  void dispose() {
    _probeDebounce?.cancel();
    _endpointCtrl.removeListener(_onEndpointChanged);
    _nameCtrl.dispose();
    _endpointCtrl.dispose();
    _roomCtrl.dispose();
    _pathCtrl.dispose();
    super.dispose();
  }

  void _onEndpointChanged() => _scheduleEndpointProbe();

  void _scheduleEndpointProbe() {
    _probeDebounce?.cancel();
    _probeDebounce = Timer(const Duration(milliseconds: 450), () {
      unawaited(_probeEndpointKind(_endpointCtrl.text.trim()));
    });
  }

  Future<void> _probeEndpointKind(String raw) async {
    final gen = ++_probeGen;
    if (raw.isEmpty) {
      if (mounted) {
        setState(() {
          _endpointKind = null;
          _probing = false;
        });
      }
      return;
    }
    if (mounted) {
      setState(() {
        _probing = true;
        _endpointKind = null;
      });
    }
    try {
      final target = await RbConnectTarget.resolve(raw, probeTimeout: const Duration(seconds: 2));
      if (gen != _probeGen || !mounted) {
        return;
      }
      setState(() {
        _endpointKind = target.kind;
        _probing = false;
      });
    } catch (_) {
      if (gen != _probeGen || !mounted) {
        return;
      }
      setState(() {
        _endpointKind = null;
        _probing = false;
      });
    }
  }

  RbConnectionProfile _buildProfile(RbConnectTarget target) {
    final id = widget.profile?.id ?? DateTime.now().millisecondsSinceEpoch.toString();
    final stored = RbConnectTarget.endpointStorage(target);
    final (signalUrl, direct) = target.isDirect ? ('', stored) : (stored, '');
    return RbConnectionProfile(
      id: id,
      name: _nameCtrl.text.trim(),
      signalUrl: signalUrl,
      direct: direct,
      room: _roomCtrl.text.trim(),
      path: _pathCtrl.text.trim(),
      agentPlatform: widget.profile?.agentPlatform ?? '',
    );
  }

  Future<void> _save() async {
    if (_nameCtrl.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('请填写名称')));
      return;
    }
    setState(() => _saving = true);
    RbConnectTarget target;
    try {
      target = await RbConnectTarget.resolve(_endpointCtrl.text.trim());
    } catch (e) {
      if (mounted) {
        setState(() => _saving = false);
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('地址无效：$e')));
      }
      return;
    }
    if (target.isSignal && _roomCtrl.text.trim().isEmpty) {
      if (mounted) {
        setState(() => _saving = false);
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('信令连接需填写房间码，与对方「共享」页一致')));
      }
      return;
    }
    final profile = _buildProfile(target);
    final list = await RbConnectionStore.loadAll();
    final next = [...list.where((p) => p.id != profile.id), profile];
    await RbConnectionStore.persist(next, lastId: profile.id);
    if (!mounted) {
      return;
    }
    AppNavigator.pop(true);
  }

  Future<void> _delete() async {
    final p = widget.profile;
    if (p == null) {
      return;
    }
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除连接'),
        content: Text('确定删除「${p.name}」？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('删除')),
        ],
      ),
    );
    if (ok != true || !mounted) {
      return;
    }
    await RbConnectionStore.remove(p.id);
    AppNavigator.pop(true);
  }

  String _kindHint() {
    if (_probing) {
      return '正在探测信令 /rb/health…';
    }
    return switch (_endpointKind) {
      RbConnectTargetKind.direct => '未检测到信令，将直连 TCP（默认端口 19091），可不填房间码',
      RbConnectTargetKind.signal => '已检测到信令服务，需填写房间码',
      null => '填写域名或 IP，可选端口（信令常见 8079；直连常见 19091，将自动探测）',
    };
  }

  @override
  Widget build(BuildContext context) {
    final hint = Theme.of(context).textTheme.bodySmall?.copyWith(color: Theme.of(context).colorScheme.onSurfaceVariant);
    return Scaffold(
      appBar: AppBar(title: Text(widget.isEdit ? '编辑连接' : '添加连接')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Text('只需填写服务器地址；协议与 /rb/signal 路径固定，保存前会自动探测信令。', style: hint),
          const SizedBox(height: 16),
          RbClearTextField(
            controller: _nameCtrl,
            textInputAction: TextInputAction.next,
            decoration: const InputDecoration(labelText: '名称', hintText: '例如：家里电脑', border: OutlineInputBorder()),
          ),
          const SizedBox(height: 12),
          RbClearTextField(
            controller: _endpointCtrl,
            textInputAction: TextInputAction.next,
            autocorrect: false,
            onChanged: (_) => _scheduleEndpointProbe(),
            decoration: const InputDecoration(
              labelText: '服务器地址',
              hintText: 'example.com:8079 或 192.168.1.10、[2001:db8::1]:8079',
              border: OutlineInputBorder(),
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(top: 4, left: 4),
            child: Text(_kindHint(), style: hint),
          ),
          const SizedBox(height: 12),
          RbClearTextField(
            controller: _roomCtrl,
            textInputAction: TextInputAction.next,
            decoration: InputDecoration(
              labelText: '房间码',
              hintText: _endpointKind == RbConnectTargetKind.direct ? '直连可留空' : '与对方共享页一致',
              border: const OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 12),
          RbClearTextField(
            controller: _pathCtrl,
            textInputAction: TextInputAction.done,
            decoration: InputDecoration(
              labelText: '起始目录（可选）',
              hintText: (widget.profile?.agentPlatform ?? '').toLowerCase().contains('win')
                  ? 'Windows 路径如 D:/Users/…，留空从 Agent 根目录'
                  : '留空则从共享根目录开始',
              border: const OutlineInputBorder(),
            ),
          ),
          const SizedBox(height: 24),
          FilledButton(onPressed: _saving ? null : _save, child: Text(_saving ? '保存中…' : '保存')),
          if (widget.isEdit) ...[
            const SizedBox(height: 12),
            OutlinedButton(onPressed: _saving ? null : _delete, child: const Text('删除连接')),
          ],
        ],
      ),
    );
  }
}
