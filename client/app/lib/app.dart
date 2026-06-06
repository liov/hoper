import 'package:app/global/state.dart';
import 'package:app/global/theme.dart';
import 'package:app/pages/route.dart';
import 'package:app/translations/zh_CN/local.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';

class AppRoot extends StatefulWidget {
  const AppRoot({super.key});

  static final GlobalKey<_AppRootState> restartKey = GlobalKey<_AppRootState>();

  static void restartApp() => restartKey.currentState?._restart();

  @override
  State<AppRoot> createState() => _AppRootState();
}

class _AppRootState extends State<AppRoot> {
  Key _appKey = UniqueKey();

  void _restart() => setState(() => _appKey = UniqueKey());

  @override
  Widget build(BuildContext context) {
    return KeyedSubtree(
      key: _appKey,
      child: ValueListenableBuilder<bool>(
        valueListenable: globalState.isDarkMode,
        builder: (_, dark, _) => MaterialApp(
          navigatorKey: rootNavigatorKey,
          title: 'hoper',
          themeMode: dark ? ThemeMode.dark : ThemeMode.light,
          theme: AppTheme.light,
          darkTheme: AppTheme.dark,
          initialRoute: Routes.START,
          onGenerateRoute: Routes.onGenerateRoute,
          onUnknownRoute: (_) => MaterialPageRoute(builder: (_) => const _NotFoundPage()),
          builder: (ctx, child) => Listener(
            onPointerDown: (_) {
              final scope = FocusScope.of(ctx);
              if (!scope.hasPrimaryFocus && scope.focusedChild != null) {
                FocusManager.instance.primaryFocus?.unfocus();
              }
            },
            child: child,
          ),
          locale: const Locale('zh', 'CN'),
          localizationsDelegates: const [
            GlobalMaterialLocalizations.delegate,
            GlobalWidgetsLocalizations.delegate,
            GlobalCupertinoLocalizations.delegate,
            ZhCupertinoLocalizations.delegate,
          ],
          supportedLocales: const [
            Locale('zh', 'CN'),
            Locale('en', 'US'),
            Locale('vi', 'VN'),
          ],
        ),
      ),
    );
  }
}

class _NotFoundPage extends StatelessWidget {
  const _NotFoundPage();

  @override
  Widget build(BuildContext context) => const Scaffold(body: Center(child: Text('找不到页面')));
}
