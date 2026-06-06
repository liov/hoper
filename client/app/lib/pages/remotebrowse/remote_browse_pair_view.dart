import 'package:app/pages/remotebrowse/remote_browse_agent_view.dart';
import 'package:app/pages/remotebrowse/remote_browse_connection_list_view.dart';
import 'package:flutter/material.dart';

class RemoteBrowsePairView extends StatelessWidget {
  const RemoteBrowsePairView({super.key});

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 2,
      child: Scaffold(
        appBar: AppBar(
          title: const Text('远程相册'),
          bottom: const TabBar(
            tabs: [
              Tab(icon: Icon(Icons.photo_album_outlined, size: 20), text: '浏览'),
              Tab(icon: Icon(Icons.cloud_upload_outlined, size: 20), text: '共享'),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            const RemoteBrowseConnectionListView(key: PageStorageKey<String>('rb_connections')),
            RemoteBrowseAgentView(key: const PageStorageKey<String>('rb_agent')),
          ],
        ),
      ),
    );
  }
}
