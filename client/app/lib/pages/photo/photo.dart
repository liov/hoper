import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:get/get_core/src/get_main.dart';

class PhotoView extends StatelessWidget {
  PhotoView(this.urls,this.initialPage, {super.key});
  final List<String> urls;
  final int initialPage;


  @override
  Widget build(BuildContext context) {
    return ExtendedImageGesturePageView.builder(
      itemBuilder: (BuildContext context, int index) {
        var url = urls[index];
        Widget image = ExtendedImage.network(
          url,
          fit: BoxFit.contain,
          mode: ExtendedImageMode.gesture,
          extendedImageGestureKey: Key(url),
          initGestureConfigHandler: (state) {
            return GestureConfig(
                minScale: 0.9,
                animationMinScale: 0.7,
                maxScale: 4.0,
                animationMaxScale: 4.5,
                speed: 1.0,
                inertialSpeed: 100.0,
                initialScale: 1.0,
                inPageView: true,
                initialAlignment: InitialAlignment.center,
                //you can cache gesture state even though page view page change.
                //remember call clearGestureDetailsCache() method at the right time.(for example,this page dispose)
                cacheGesture: false
            );
          },
        );
        if (index == initialPage) {
          return Hero(
            tag: url,
            child: image,
          );
        } else {
          return Container(
            padding: const EdgeInsets.all(5.0),
            child: image,
          );
        }
      },
      itemCount: urls.length,
      onPageChanged: (int index) {
      },
      controller: ExtendedPageController(
        initialPage: initialPage,
      ),
      scrollDirection: Axis.horizontal,
    );
  }

}