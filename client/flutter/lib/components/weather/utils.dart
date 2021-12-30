
import 'dart:ui';

import 'package:flutter/services.dart';
import 'package:flutter/widgets.dart';
import 'dart:ui' as ui;

/// 图片相关的工具类
class ImageUtils {
  /// 绘制时需要用到 ui.Image 的对象，通过此方法进行转换
  static Future<ui.Image> getImage(String asset) async {
    ByteData data = await rootBundle.load("assets/images/weather/$asset");
    Codec codec = await instantiateImageCodec(data.buffer.asUint8List());
    FrameInfo fi = await codec.getNextFrame();
    return fi.image;
  }
}

/// 定义打印函数
typedef WeatherPrint = void Function(String message,
    {int wrapWidth, String tag});

const DEBUG = true;

WeatherPrint weatherPrint = debugPrintThrottled;

// 统一方法进行打印
void debugPrintThrottled(String message, {int? wrapWidth, String? tag}) {
  if (DEBUG) {
    debugPrint("flutter-weather: $tag: $message", wrapWidth: wrapWidth);
  }
}
