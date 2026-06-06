import 'dart:io';

import 'package:app/components/camrea/camera_view.dart';
import 'package:app/global/service.dart';
import 'package:app/global/state.dart';
import 'package:app/util/app_permission.dart';
import 'package:flutter/foundation.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';

class MediaActions {
  List<XFile> imageFiles = List<XFile>.empty(growable: true);
  List<String> imageUrls = List.empty(growable: true);
  dynamic pickImageError;
  VideoPlayerController? videoController;
  String? retrieveDataError;
  final ImagePicker picker = ImagePicker();

  Future<void> onImageButtonPressed(ImageSource source, {bool isCamera = false}) async {
    if (videoController != null) {
      await videoController!.setVolume(0.0);
    }
    if (isCamera || source == ImageSource.camera) {
      if (!await ensureCameraPermission()) return;
    }
    if (isCamera) {
      final file = await getPhoto();
      if (file == null) return;
      final url = await globalService.uploadClient.upload(File(file.path));
      imageUrls.add(url);
      imageFiles.add(file);
    } else {
      final file = await getPhoto2();
      if (file == null) return;
      final url = await globalService.uploadClient.upload(File(file.path));
      imageUrls.add(url);
      imageFiles.add(file);
    }
  }

  Future<void> playVideo(XFile? file) async {
    if (file == null) return;
    await disposeVideoController();
    late VideoPlayerController controller;
    if (kIsWeb) {
      controller = VideoPlayerController.network(file.path);
    } else {
      controller = VideoPlayerController.file(File(file.path));
    }
    videoController = controller;
    const double volume = kIsWeb ? 0.0 : 1.0;
    await controller.setVolume(volume);
    await controller.initialize();
    await controller.setLooping(true);
    await controller.play();
  }

  Future<void> disposeVideoController() async {
    if (videoController != null) {
      await videoController!.dispose();
    }
    videoController = null;
  }
}
