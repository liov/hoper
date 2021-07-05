import 'package:app/pages/home/home_binding.dart';
import 'package:app/pages/home/home_view.dart';
import 'package:app/pages/loginView.dart';
import 'package:app/pages/moment/list/moment_list_view.dart';
import 'package:app/pages/moment/moment_binding.dart';
import 'package:get/get.dart';

class AppPages {
  static final routes = [
    GetPage(
      name: Routes.HOME,
      page: () => HomeView(),
      bindings: [HomeBinding(),MomentBinding()],
    ),
    GetPage(
      name: Routes.MOMENT,
      page: () => MomentListView(),
      binding: MomentBinding(),
    ),
    GetPage(
      name: Routes.LOGIN,
      page: () => LoginView(),
    ),
  ];
}

abstract class Routes {
  Routes._();

  static const HOME = _Paths.HOME;
  static const LOGIN = _Paths.LOGIN;
  static const CONTENT = _Paths.CONTENT;
  static const MOMENT = _Paths.MOMENT;

  static const SETTINGS = _Paths.SETTINGS;

  static const PRODUCTS = _Paths.HOME + _Paths.PRODUCTS;
  static String PRODUCT_DETAILS(String productId) => '$PRODUCTS/$productId';
}

abstract class _Paths {
  static const HOME = '/home';
  static const CONTENT = '/content';
  static const MOMENT = '/moment';
  static const LOGIN = '/login';
  static const SETTINGS = '/settings';
  static const PRODUCTS = '/products';
  static const PRODUCT_DETAILS = '/:productId';
}