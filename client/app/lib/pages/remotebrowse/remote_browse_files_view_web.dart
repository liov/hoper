import 'package:app/pages/remotebrowse/rb_ui.dart';
import 'package:app/remotebrowse/rb_connection_profile.dart';
import 'package:flutter/material.dart';

class RemoteBrowseFilesView extends StatelessWidget {
  const RemoteBrowseFilesView({super.key, required this.profile});

  final RbConnectionProfile profile;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text(profile.name)),
      body: Center(
        child: RbEmptyState(
          icon: Icons.phone_iphone_outlined,
          title: '请使用客户端 App',
          subtitle: 'Web 端暂不支持浏览远程文件。\n请使用 iOS、Android 或 macOS 打开「${profile.name}」。',
        ),
      ),
    );
  }
}
