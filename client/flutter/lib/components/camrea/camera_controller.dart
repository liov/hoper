import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:video_player/video_player.dart';

class Controller extends GetxController {

    CameraController? controller;
    XFile? imageFile;
    XFile? videoFile;
    VideoPlayerController? videoController;
    VoidCallback? videoPlayerListener;
    bool enableAudio = true;
    double _minAvailableExposureOffset = 0.0;
    double _maxAvailableExposureOffset = 0.0;
    double _currentExposureOffset = 0.0;
    late AnimationController _flashModeControlRowAnimationController;
    late Animation<double> _flashModeControlRowAnimation;
    late AnimationController _exposureModeControlRowAnimationController;
    late Animation<double> _exposureModeControlRowAnimation;
    late AnimationController _focusModeControlRowAnimationController;
    late Animation<double> _focusModeControlRowAnimation;
    double _minAvailableZoom = 1.0;
    double _maxAvailableZoom = 1.0;
    double _currentScale = 1.0;
    double _baseScale = 1.0;
    int _pointers = 0;
    setController (CameraController? c) async{
     controller?.dispose();
     controller = c;
     await controller?.initialize();
     update();
    }

    /// Returns a suitable camera icon for [direction].
    IconData getCameraLensIcon(CameraLensDirection direction) {
      switch (direction) {
        case CameraLensDirection.back:
          return Icons.camera_rear;
        case CameraLensDirection.front:
          return Icons.camera_front;
        case CameraLensDirection.external:
          return Icons.camera;
        default:
          throw ArgumentError('Unknown lens direction');
      }
    }

    void logError(String code, String? message) {
      if (message != null) {
        print('Error: $code\nError Message: $message');
      } else {
        print('Error: $code');
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
