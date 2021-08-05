import 'package:app/global/global_controller.dart';
import 'package:app/pages/home/home_binding.dart';
import 'package:app/pages/home/home_view.dart';
import 'package:app/pages/home/splash_view.dart';
import 'package:app/pages/login_view.dart';
import 'package:app/pages/moment/add/moment_add_controller.dart';
import 'package:app/pages/moment/add/moment_add_view.dart';
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
      children:[
        GetPage(
          name: Routes.ADD,
          page: () =>  globalController.authCheck() ?? MomentAddView(),
          binding: BindingsBuilder.put(() => MomentAddController())
        ),
      ]
    ),
    GetPage(
      name: Routes.LOGIN,
      page: () => LoginView(),
    ),
    GetPage(
      name: Routes.LOGIN,
      page: () => LoginView(),
    ),
    GetPage(
      name: Routes.SPLASH,
      page: () => Splash(),
    ),
  ];
}

abstract class Routes {
  Routes._();

  static const HOME = '/home';
  static const CONTENT = '/content';
  static const MOMENT = '/moment';
  static const ADD = '/add';
  static const MOMENT_ADD = MOMENT + ADD;
  static const LOGIN = '/login';
  static const SETTINGS = '/settings';
  static const SPLASH = '/splash';
  static const PRODUCTS = HOME + '/products';
  static const PRODUCT_DETAILS = '/:productId';

  static String productDetails(String productId) => '$PRODUCTS/$productId';
}