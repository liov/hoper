import 'package:extended_image/extended_image.dart';
import 'package:flutter/widgets.dart';

Widget extendedImageFile(
  String path, {
  AlignmentGeometry alignment = Alignment.center,
  BoxFit fit = BoxFit.contain,
  bool enableSlideOutPage = false,
}) {
  return ExtendedImage.network(
    path,
    alignment: alignment,
    fit: fit,
    enableSlideOutPage: enableSlideOutPage,
  );
}
