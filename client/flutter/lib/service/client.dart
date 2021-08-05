

import 'package:get/get_state_manager/src/rx_flutter/rx_disposable.dart';
import 'package:app/utils/dio.dart' as $dio;

class Client extends GetxService {
  final httpClient = $dio.httpClient;
}