



import 'package:app/global/state.dart';
import 'package:app/global/theme.dart';
import 'package:app/pages/splash/splash_view.dart';
import 'package:app/pages/index/index.dart';
import 'package:app/pages/moment/moment_view.dart';
import 'package:app/pages/user/user_view.dart';
import 'package:app/pages/route.dart';
import 'package:app/providers/providers.dart';
import 'package:app/util/nav.dart';

import 'package:convex_bottom_bar/convex_bottom_bar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:app/components/media/pick.dart';

class App extends ConsumerWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return const _AppShell();
  }
}

class _AppShell extends ConsumerStatefulWidget {
  const _AppShell();

  @override
  ConsumerState<_AppShell> createState() => _AppState();
}

class _AppState extends ConsumerState<_AppShell> with WidgetsBindingObserver {
  static const _tabs = <Widget>[
    MomentView(),
    IndexPage(),
    MediaPick(),
    UserView(),
  ];

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    final splash = ref.read(splashProvider.notifier);
    switch (state) {
      case AppLifecycleState.inactive:
        break;
      case AppLifecycleState.resumed:
        splash.advertising(splashWidget);
        break;
      case AppLifecycleState.paused:
        splash.pausedTime = DateTime.now();
        break;
      case AppLifecycleState.detached:
      case AppLifecycleState.hidden:
        break;
    }
  }

  @override
  Widget build(BuildContext context) {
    final pageIndex = ref.watch(homeProvider.select((s) => s.pageIndex));
    final selectedIndex = ref.watch(homeProvider.select((s) => s.selectedIndex));
    final home = ref.read(homeProvider.notifier);
    final theme = globalState.isDarkMode.value ? AppTheme.dark : AppTheme.light;
    return Scaffold(
      body: IndexedStack(
        index: pageIndex,
        sizing: StackFit.expand,
        children: _tabs,
      ),
      bottomNavigationBar: ConvexAppBar(
        disableDefaultTabController: true,
        initialActiveIndex: selectedIndex,
        onTabNotify: (i) {
          if (i == 2) {
            AppNavigator.pushNamed(Routes.MOMENT_ADD);
            return false;
          }
          return true;
        },
        onTap: home.onItemTapped,
        activeColor: theme.canvasColor,
        backgroundColor: theme.primaryColor,
        style: TabStyle.fixedCircle,
        items: Home.bottomNavigationBarList.map((e) => e.tabItem()).toList(),
      ),
    );
  }
}
