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
          bottom: const TabBar(tabs: [Tab(text: '浏览'), Tab(text: '共享')]),
        ),
        body: TabBarView(
          children: [
            const RemoteBrowseConnectionListView(key: PageStorageKey<String>('rb_connections')),
            Center(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Text(
                  'Web 端暂不支持启动文件端共享。\n请使用 iOS / Android / macOS 客户端。',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.bodyLarge,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
