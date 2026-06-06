import 'dart:async';

import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/pages/remotebrowse/remote_browse_connection_edit_view.dart';
import 'package:app/pages/remotebrowse/remote_browse_files_entry.dart';
import 'package:app/remotebrowse/api.dart';
import 'package:app/remotebrowse/rb_connect_endpoint.dart';
import 'package:app/remotebrowse/rb_connection_profile.dart';
import 'package:app/remotebrowse/rb_connection_store.dart';
import 'package:app/remotebrowse/rb_remote_os.dart';
import 'package:flutter/foundation.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/material.dart';

class RemoteBrowseConnectionListView extends StatefulWidget {
  const RemoteBrowseConnectionListView({super.key});

  @override
  State<RemoteBrowseConnectionListView> createState() => _RemoteBrowseConnectionListViewState();
}

class _RemoteBrowseConnectionListViewState extends State<RemoteBrowseConnectionListView> {
  List<RbConnectionProfile> _profiles = [];
  final _platformById = <String, String>{};
  var _loading = true;

  @override
  void initState() {
    super.initState();
    _reload();
  }

  Future<void> _reload() async {
    setState(() => _loading = true);
    final list = await RbConnectionStore.loadAll();
    if (!mounted) {
      return;
    }
    final plat = <String, String>{};
    for (final p in list) {
      if (p.agentPlatform.isNotEmpty) {
        plat[p.id] = p.agentPlatform;
      }
    }
    setState(() {
      _profiles = list;
      _platformById
        ..clear()
        ..addAll(plat);
      _loading = false;
    });
    unawaited(_refreshAgentPlatforms(list));
  }

  RbRemoteOs _osFor(RbConnectionProfile p) => rbParseRemoteOs(_platformById[p.id] ?? p.agentPlatform);

  Future<void> _refreshAgentPlatforms(List<RbConnectionProfile> list) async {
    final next = Map<String, String>.from(_platformById);
    var changed = false;
    await Future.wait(list.map((p) async {
      if (p.room.trim().isEmpty) {
        return;
      }
      final target = RbConnectTarget.fromProfile(p);
      if (!target.isSignal) {
        return;
      }
      try {
        final plat = await probeAgentPlatform(target.signalUri!, p.room);
        if (plat == null || plat.isEmpty) {
          return;
        }
        if (next[p.id] != plat) {
          next[p.id] = plat;
          changed = true;
        }
        if (p.agentPlatform != plat) {
          await _persistAgentPlatform(p.id, plat);
        }
      } catch (_) {}
    }));
    if (changed && mounted) {
      setState(() => _platformById.addAll(next));
    }
  }

  Future<void> _persistAgentPlatform(String id, String platform) async {
    final list = await RbConnectionStore.loadAll();
    final i = list.indexWhere((p) => p.id == id);
    if (i < 0) {
      return;
    }
    final old = list[i];
    if (old.agentPlatform == platform) {
      return;
    }
    list[i] = RbConnectionProfile(
      id: old.id,
      name: old.name,
      signalUrl: old.signalUrl,
      direct: old.direct,
      room: old.room,
      path: old.path,
      agentPlatform: platform,
    );
    await RbConnectionStore.persist(list);
  }

  Future<void> _openAdd() async {
    final saved = await AppNavigator.push<bool>(const RemoteBrowseConnectionEditView());
    if (saved == true) {
      await _reload();
    }
  }

  Future<void> _openBrowse(RbConnectionProfile profile) async {
    if (kIsWeb) {
      if (!mounted) {
        return;
      }
      await showDialog<void>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('请使用客户端 App'),
          content: Text('Web 端暂不支持浏览「${profile.name}」中的文件。\n请使用 iOS、Android 或 macOS 客户端。'),
          actions: [TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('知道了'))],
        ),
      );
      return;
    }
    await RbConnectionStore.persist(_profiles, lastId: profile.id);
    await AppNavigator.push(RemoteBrowseFilesView(profile: profile));
    await _reload();
  }

  Future<void> _openEdit(RbConnectionProfile profile) async {
    final saved = await AppNavigator.push<bool>(RemoteBrowseConnectionEditView(profile: profile));
    if (saved == true) {
      await _reload();
    }
  }

  Future<void> _delete(RbConnectionProfile profile) async {
    final ok = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('删除连接'),
        content: Text('确定删除「${profile.name}」？'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('取消')),
          FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('删除')),
        ],
      ),
    );
    if (ok != true || !mounted) {
      return;
    }
    await RbConnectionStore.remove(profile.id);
    await _reload();
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('已删除「${profile.name}」')));
    }
  }

  String _subtitle(RbConnectionProfile p) {
    final ep = RbConnectTarget.endpointDisplay(p);
    final room = p.room.isEmpty ? '' : p.room;
    if (p.direct.trim().isNotEmpty) {
      return room.isEmpty ? '直连 · $ep' : '$room · 直连 $ep';
    }
    return room.isEmpty ? '信令 · $ep' : '$room · 信令 $ep';
  }

  @override
  Widget build(BuildContext context) {
    Widget body;
    if (_loading) {
      body = const Center(child: CircularProgressIndicator());
    } else if (_profiles.isEmpty) {
      body = Center(
        child: RbEmptyState(
          icon: Icons.devices_outlined,
          title: '还没有远程连接',
          subtitle: '添加一台已开启「共享」的设备，即可浏览对方相册与文件。',
          action: FilledButton.icon(onPressed: _openAdd, icon: const Icon(Icons.add), label: const Text('添加连接')),
        ),
      );
    } else {
      body = RefreshIndicator(
        onRefresh: _reload,
        child: ListView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.only(top: 8, bottom: 80),
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 4, 16, 8),
              child: Text(
                '点选设备进入相册；对方需先在「共享」页开启服务。',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Theme.of(context).colorScheme.onSurfaceVariant),
              ),
            ),
            ..._profiles.map(
              (p) => RbConnectionTile(
                name: p.name,
                subtitle: _subtitle(p),
                remoteOs: _osFor(p),
                onTap: () => _openBrowse(p),
                onEdit: () => _openEdit(p),
                onDelete: () => _delete(p),
              ),
            ),
          ],
        ),
      );
    }
    return Scaffold(
      body: body,
      floatingActionButton: _profiles.isEmpty && !_loading
          ? null
          : FloatingActionButton(
              onPressed: _openAdd,
              child: const Icon(Icons.add),
            ),
    );
  }
}
