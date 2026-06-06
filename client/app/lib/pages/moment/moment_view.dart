import 'package:app/global/state.dart';
import 'package:app/util/nav.dart';
import 'package:app/pages/route.dart';
import 'package:flutter/material.dart';

import '../user/login_view.dart';
import 'list/moment_list_view.dart';

class MomentView extends StatefulWidget {
  const MomentView({super.key});

  @override
  State<MomentView> createState() => _MomentState();
}

class _MomentState extends State<MomentView> with AutomaticKeepAliveClientMixin, TickerProviderStateMixin {
  static const _tabValues = ['关注', '推荐', '刚刚'];
  late final TabController _tabController = TabController(length: _tabValues.length, vsync: this, initialIndex: 2);

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: TabBar(
          isScrollable: true,
          tabs: [for (final choice in _tabValues) Tab(text: choice)],
          controller: _tabController,
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              if (globalState.authState.userAuth != null) {
                AppNavigator.pushNamed(Routes.MOMENT_ADD);
              } else {
                AppNavigator.push(LoginView());
              }
            },
          ),
        ],
      ),
      body: TabBarView(
        controller: _tabController,
        children: const [
          Center(child: Text('TODO')),
          Center(child: Text('TODO')),
          MomentListView(tag: 'newest'),
        ],
      ),
    );
  }

  @override
  bool get wantKeepAlive => true;
}
