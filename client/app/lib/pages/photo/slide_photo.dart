import 'dart:io';

import 'package:app/components/hero.dart';
import 'package:app/components/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';

class SlidePhotoView extends StatelessWidget {
  SlidePhotoView(this.url, {super.key});
  final String url;
  final _slidePageKey  = GlobalKey<ExtendedImageSlidePageState>();

  @override
  Widget build(BuildContext context) {
    return Material(
      type: MaterialType.transparency,
      child: ExtendedImageSlidePage(
        key: _slidePageKey,
        slideAxis: SlideAxis.both,
        slideType: SlideType.onlyImage,
        child: GestureDetector(
          child:  HeroWidget(
            tag: url,
            slideType: SlideType.onlyImage,
            slidePageKey: _slidePageKey,
            child: url.startsWith("http")?ExtendedImage.network(
              url,
              enableSlideOutPage: true,
            ):ExtendedImage.file(
              File(url),
              enableSlideOutPage: true,
            ),
          ),
          onTap: () {
            _slidePageKey.currentState!.popPage();
            Navigator.pop(context);
          },
        ),
      ),
    );
  }
}

slidePhotoRoute(String url){
  navigator!.push(SimpleRoute(builder:(context)=>SlidePhotoView(url)));
}