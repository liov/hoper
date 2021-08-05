import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:app/components/async/async.dart';
import 'camera_controller.dart';
import 'camera_example.dart';

class CameraView extends StatefulWidget {
  CameraView(this.cameras);

  final List<CameraDescription> cameras;

  @override
  CameraViewState createState() => CameraViewState();
}

class CameraViewState extends State<CameraView> with WidgetsBindingObserver {
  CameraController? controller;

  final GlobalKey<ScaffoldState> _scaffoldKey = GlobalKey<ScaffoldState>();
  final _aspectRatio = Get.width / Get.height;
  final _scale = 1.0;
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: _scaffoldKey,
      body: Container(
        color: Colors.black,
        child: Transform.scale(
          scale: _scale,
          child: Center(
              child: FutureBuilder(
              future: onNewCameraSelected(widget.cameras[0]),
              builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
                return snapshot.handle() ??
                    AspectRatio(
                      aspectRatio: _aspectRatio ,
                      child: CameraPreview(controller!,
                          child: Align(
                              alignment: Alignment.bottomCenter,
                              child: Padding(
                                padding: const EdgeInsets.all(10.0),
                                child:ElevatedButton(
                                    style: ButtonStyle(
                                        backgroundColor:
                                        MaterialStateProperty.resolveWith(
                                                (states) => (Colors.red)),
                                        shape: MaterialStateProperty.resolveWith(
                                                (states) => CircleBorder()),
                                        minimumSize:
                                        MaterialStateProperty.resolveWith(
                                                (states) => Size(100, 100))),
                                    child: Text('按'),
                                    onPressed: () {
                                      controller!
                                          .takePicture()
                                          .then((value) => navigator!.pop(value));
                                    })
                              ))),
                  );
            },
          )),
        ),
      ),
    );
  }

  double _getImageZoom(MediaQueryData data) {
    final double logicalWidth = data.size.width;
    final double logicalHeight = _aspectRatio * logicalWidth;

    final EdgeInsets padding = data.padding;
    final double maxLogicalHeight =
        data.size.height - padding.top - padding.bottom;
    return maxLogicalHeight / logicalHeight;
  }

  onNewCameraSelected(CameraDescription cameraDescription) async {
    if (controller != null) {
      await controller!.dispose();
    }
    final CameraController cameraController = CameraController(
      cameraDescription,
      ResolutionPreset.veryHigh,
      enableAudio: true,
      imageFormatGroup: ImageFormatGroup.jpeg,
    );
    controller = cameraController;

    // If the controller is updated then update the UI.
/*    cameraController.addListener(() {
      if (mounted) setState(() {});
      if (cameraController.value.hasError) {
        Get.rawSnackbar(message:'Camera error ${cameraController.value.errorDescription}');
      }
    });*/

    try {
      await cameraController.initialize();
    } on CameraException catch (e) {
      _showCameraException(e);
    }

    /*if (mounted) {
      setState(() {});
    }*/
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance?.addObserver(this);
  }

  @override
  void dispose() {
    print('调用了');
    WidgetsBinding.instance?.removeObserver(this);
    controller?.dispose();
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    print("--" + state.toString());
    // App state changed before we got the chance to initialize.
    if (controller == null || !controller!.value.isInitialized) {
      return;
    }
    if (state == AppLifecycleState.inactive) {
      controller!.dispose();
    } else if (state == AppLifecycleState.resumed) {
      onNewCameraSelected(controller!.description);
      if (mounted) {
        setState(() {});
      }
    }
  }
}

void _showCameraException(CameraException e) {
  logError(e.code, e.description);
  Get.rawSnackbar(message: 'Error: ${e.code}\n${e.description}');
}

Future<XFile?> getPhoto() async {
  final cameras = await availableCameras();
  print(cameras);
  if (cameras.length == 0) return null;
  XFile? file = await Get.to(() => CameraView(cameras));
  return file;
}

Future<XFile?> getPhoto2() async {
  final cameras = await availableCameras();
  print(cameras);
  if (cameras.length == 0) return null;
  XFile? file = await Get.to(() => CameraExample(cameras));
  return file;
}
