import 'dart:io';

import 'package:app/components/camrea/camera_view.dart';
import 'package:app/generated/protobuf/content/content.enum.pb.dart';
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/global_state.dart';
import 'package:app/service/moment.dart';
import 'package:app/utils/media.dart';
import 'package:camera/camera.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';


class MomentAddController extends GetxController with MediaController{

  String content = '';
  final MomentClient momentClient = Get.find();



  Future<void> save() async {
    try {
      await momentClient.stub.add(AddMomentReq(
        type: MomentType.MomentTypeImage,
        content: content,
        images: imageUrls.join(','),
      ));
      navigator!.pop();
    }catch (e) {
      print(e);
    }
  }



  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
