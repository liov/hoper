import 'package:app/components/hero.dart';
import 'package:app/components/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';

class SlidePhotoView extends StatelessWidget {
  SlidePhotoView(this.url) :super();
  final String url;
  final _slidePageKey  = GlobalKey<ExtendedImageSlidePageState>();

  @override
  Widget build(BuildContext context) {
    return Material(
      type: MaterialType.transparency,
      child: ExtendedImageSlidePage(
        key: _slidePageKey,
        child: GestureDetector(
          child:  HeroWidget(
            child: url.startsWith("http")?ExtendedImage.network(
              url,
              enableSlideOutPage: true,
            ):ExtendedImage.file(
              File(url),
              enableSlideOutPage: true,
            ),
            tag: url,
            slideType: SlideType.onlyImage,
            slidePageKey: _slidePageKey,
          ),
          onTap: () {
            _slidePageKey.currentState!.popPage();
            Navigator.pop(context);
          },
        ),
        slideAxis: SlideAxis.both,
        slideType: SlideType.onlyImage,
      ),
    );
  }
}

slidePhotoRoute(String url){
  navigator!.push(SimpleRoute(builder:(context)=>SlidePhotoView(url)));
}