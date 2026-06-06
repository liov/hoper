import 'package:app/components/hero.dart';
import 'package:app/util/image_file.dart';
import 'package:app/util/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:app/util/nav.dart';

class SlideImageView extends StatelessWidget {
  SlideImageView(this.url, {super.key});

  final String url;
  final _slidePageKey = GlobalKey<ExtendedImageSlidePageState>();

  @override
  Widget build(BuildContext context) {
    return Material(
      type: MaterialType.transparency,
      child: GestureDetector(
        child: ExtendedImageSlidePage(
          key: _slidePageKey,
          slideAxis: SlideAxis.both,
          slideType: SlideType.onlyImage,

          child: HeroImage(
            tag: url,
            slideType: SlideType.onlyImage,
            slidePageKey: _slidePageKey,
            child: url.startsWith("http")
                    ? ExtendedImage.network(url, enableSlideOutPage: true)
                    : extendedImageFile(url, enableSlideOutPage: true),
          ),
        ),
        onTap: () {
          _slidePageKey.currentState!.popPage();
          Navigator.pop(context);
        },
      ),
    );
  }
}

void slideImageRoute(String url) {
  AppNavigator.nav!.push(SimpleRoute(builder: (context) => SlideImageView(url)));
}
