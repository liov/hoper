import 'dart:io';

import 'package:camera/camera.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';


class MomentAddController extends GetxController {

  List<XFile>? imageFileList;

  set imageFile(XFile? value) {
    imageFileList = value == null ? null : [value];
  }

  dynamic pickImageError;

  VideoPlayerController? controller;
  VideoPlayerController? toBeDisposed;
  String? retrieveDataError;

  final ImagePicker picker = ImagePicker();

  void onImageButtonPressed(ImageSource source,
      {BuildContext? context, bool isMultiImage = false,bool isVideo = false}) async {
    if (controller != null) {
      await controller!.setVolume(0.0);
    }
    if (isVideo) {
      final XFile? file = await picker.pickVideo(
          source: source, maxDuration: const Duration(seconds: 10));
      await playVideo(file);
    } else  {
            try {
              final pickedFileList = await picker.pickMultiImage();
                imageFileList = pickedFileList;

            } catch (e) {
                pickImageError = e;
            }

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
      controller = controller;
      // In web, most browsers won't honor a programmatic call to .play
      // if the video has a sound track (and is not muted).
      // Mute the video so it auto-plays in web!
      // This is not needed if the call to .play is the result of user
      // interaction (clicking on a "play" button, for example).
      final double volume = kIsWeb ? 0.0 : 1.0;
      await controller.setVolume(volume);
      await controller.initialize();
      await controller.setLooping(true);
      await controller.play();
    }
  }

  Future<void> disposeVideoController() async {
    if (toBeDisposed != null) {
      await toBeDisposed!.dispose();
    }
    toBeDisposed = controller;
    controller = null;
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
