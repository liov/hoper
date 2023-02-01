
import 'package:extended_image/extended_image.dart';
import 'package:flutter/cupertino.dart';
import 'package:get/get.dart';
import 'dart:io';

class PictureView extends StatefulWidget {
  const PictureView({super.key});

  @override
  State<StatefulWidget> createState()=> PictureViewState();
}


class PictureViewState extends State<PictureView> {

  PictureViewState();
  int picWidth = 100;
  Directory dir = Directory("F:/Pictures/pron/weibo");
 List<File> images = [];

  @override
  void initState(){
    super.initState();
    dir.list().forEach((element) {
      if (element is Directory) {
          element.list().forEach((element1) {
            if (element1 is Directory) {
              element1.list().forEach((file) {
                if (file is File) {
                  images.add(file);
                }
              });
            }
          });
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final width = MediaQuery.of(context).size.width;
    final height =  MediaQuery.of(context).size.height;
    final columns = (width/picWidth).floor();
    final rows= (height/picWidth).floor();
    return Column(
       children: [
         for(int i=0; i< rows; i++)Row(
           children: [
             for(int j=0; j< columns; j++) ExtendedImage.file(images[i*columns + j],
               width: picWidth.toDouble(),
               height: picWidth.toDouble(),
               fit: BoxFit.scaleDown,)
           ]
         )
       ],
    );
  }
}