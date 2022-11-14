import 'dart:io';

import 'package:app/components/camrea/camera_view.dart';
import 'package:app/generated/protobuf/content/content.enum.pb.dart';
import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/global_state.dart';
import 'package:app/service/moment.dart';
import 'package:camera/camera.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';

class MediaController {

  List<XFile> imageFiles = List<XFile>.empty(growable: true);
  List<String> imageUrls = List.empty(growable: true);
  dynamic pickImageError;

  VideoPlayerController? videoController;

  String? retrieveDataError;

  final ImagePicker picker = ImagePicker();

  void onImageButtonPressed(ImageSource source,
      {BuildContext? context, bool isMultiImage = false,bool isCamera = false}) async {
    if (videoController != null) {
      await videoController!.setVolume(0.0);
    }
    if (isCamera) {
      final XFile? file = await getPhoto();
      if(file==null){
        return;
      }
      final url = await globalService.uploadClient.upload(File(file.path));
      imageUrls.add(url);
      imageFiles.add(file);
    } else  {
      final XFile? file = await getPhoto2();
      if(file==null){
        return;
      }
      final url = await globalService.uploadClient.upload(File(file.path));
      imageUrls.add(url);
      imageFiles.add(file);
/*            try {
              final pickedFileList = await picker.pickMultiImage();
                imageFileList = pickedFileList;
            } catch (e) {
                pickImageError = e;
            }*/
    }
  }

  Future<void> playVideo(XFile? file) async {
    if (file != null) {
      await disposeVideoController();
      late VideoPlayerController controller;
      if (kIsWeb) {
        controller = VideoPlayerController.network(file.path);
      } else {
        controller = VideoPlayerController.file(File(file.path));
      }
      videoController = controller;
      // In web, most browsers won't honor a programmatic call to .play
      // if the video has a sound track (and is not muted).
      // Mute the video so it auto-plays in web!
      // This is not needed if the call to .play is the result of user
      // interaction (clicking on a "play" button, for example).
      const double volume = kIsWeb ? 0.0 : 1.0;
      await controller.setVolume(volume);
      await controller.initialize();
      await controller.setLooping(true);
      await controller.play();
    }
  }

  Future<void> disposeVideoController() async {
    if (videoController != null) {
      await videoController!.dispose();
    }
    videoController = null;
  }


}