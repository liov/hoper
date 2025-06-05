import 'dart:io';

import 'package:app/components/hero.dart';
import 'package:applib/util/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';

class SlideImageView extends StatelessWidget {
  SlideImageView(this.url, {super.key});
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
          child:  HeroImage(
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

slideImageRoute(String url){
  navigator!.push(SimpleRoute(builder:(context)=>SlideImageView(url)));
}