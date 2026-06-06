import 'package:app/gen/pb/content/content.model.pbenum.dart';
import 'package:app/global/state.dart';
import 'package:app/pages/moment/add/moment_add_view.dart';
import 'package:app/gen/pb/content/moment.model.pb.dart' as $pb;
import 'package:app/pages/moment/detail/moment_detail_view.dart';
import 'package:app/pages/moment/list/moment_list_view.dart';
import 'package:app/pages/remotebrowse/remote_browse_entry.dart';
import 'package:app/pages/splash/splash_view.dart';
import 'package:app/pages/splash/start_view.dart';
import 'package:app/pages/user/login_view.dart';
import 'package:app/pages/webview/webview.dart';
import 'package:fixnum/fixnum.dart';
import 'package:flutter/material.dart';

abstract class Routes {
  Routes._();

  static const START = '/';
  static const HOME = '/home';
  static const CONTENT = '/content';
  static const MOMENT = '/moment';
  static const ADD = '/add';
  static const MOMENT_ADD = '$MOMENT$ADD';
  static const LOGIN = '/login';
  static const SETTINGS = '/settings';
  static const SPLASH = '/splash';
  static const PRODUCT = '/product';
  static const DynamicId = '/:id';
  static const PRODUCT_DETAILS = '$PRODUCT$DynamicId';
  static const MOMENT_DETAILS = '$MOMENT$DynamicId';
  static const WEBVIEW = '/webview';
  static const REMOTE_BROWSE = '/remotebrowse';
  static const NOTFOUND = '/NOTFOUND';

  static String productDetails(String productId) => '$PRODUCT/$productId';
  static String momentDetails(String momentId) => '$MOMENT/$momentId';
  static String contentDetails(ContentType type, Int64 contentId) => '${getContentRoute(type)}/$contentId';

  static String getContentRoute(ContentType type) {
    switch (type) {
      case ContentType.ContentMoment:
        return Routes.MOMENT;
      case ContentType.ContentPlaceholder:
      case ContentType.ContentNote:
      case ContentType.ContentDairy:
      case ContentType.ContentDairyBook:
      case ContentType.ContentFavorites:
      case ContentType.ContentCollection:
      case ContentType.ContentComment:
        return Routes.NOTFOUND;
      default:
        return Routes.NOTFOUND;
    }
  }

  static Widget authCheck(Widget Function() builder) =>
      globalState.authState.userAuth == null ? LoginView() : builder();

  static Route<dynamic> _route(RouteSettings settings, Widget child) =>
      MaterialPageRoute(builder: (_) => child, settings: settings);

  static Route<dynamic>? onGenerateRoute(RouteSettings settings) {
    final name = settings.name ?? START;
    if (name == START) {
      return _route(settings, const StartView());
    }
    if (name == HOME || name == REMOTE_BROWSE) {
      return _route(settings, const RemoteBrowsePairView());
    }
    if (name == LOGIN) {
      return _route(settings, LoginView());
    }
    if (name == SPLASH) {
      return _route(settings, const Splash());
    }
    if (name == WEBVIEW) {
      return _route(settings, const WebViewExample());
    }
    if (name == MOMENT || name == '$MOMENT/') {
      return _route(settings, const MomentListView());
    }
    if (name == MOMENT_ADD) {
      return _route(settings, authCheck(() => MomentAddView()));
    }
    if (name.startsWith('$MOMENT/')) {
      final idStr = name.substring(MOMENT.length + 1);
      if (idStr.isNotEmpty && idStr != ADD.substring(1)) {
        final id = Int64.parseInt(idStr);
        final arg = settings.arguments;
        if (arg != null && arg is $pb.Moment) {
          return _route(settings, MomentDetailView.detail(arg));
        }
        return _route(settings, MomentDetailView.byId(id));
      }
    }
    if (name == NOTFOUND) {
      return _route(settings, const Scaffold(body: Center(child: Text('找不到页面'))));
    }
    return _route(settings, const Scaffold(body: Center(child: Text('找不到页面'))));
  }
}
