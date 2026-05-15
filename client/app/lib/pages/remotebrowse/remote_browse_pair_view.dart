import 'package:app/pages/remotebrowse/remote_browse_agent_view.dart';
import 'package:app/pages/remotebrowse/remote_browse_view.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

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
        body: const TabBarView(
          children: [RemoteBrowseView(), RemoteBrowseAgentView()],
        ),
      ),
    );
  }
}

class RemoteBrowsePairBinding extends Bindings {
  @override
  void dependencies() {}
}
