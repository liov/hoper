import 'package:app/util/dialog.dart';
import 'package:permission_handler/permission_handler.dart';

/// 打开相机/录像前调用；相册选图一般由 image_picker 自行处理（iOS 14+ PHPicker 常无需相册权限）。
Future<bool> ensureCameraPermission({bool openSettingsIfDenied = true}) async {
  var status = await Permission.camera.status;
  if (status.isGranted) return true;
  if (status.isDenied) status = await Permission.camera.request();
  if (status.isGranted) return true;
  if (status.isPermanentlyDenied || status.isRestricted) {
    if (openSettingsIfDenied) {
      toast('请在系统设置中允许访问相机');
      await openAppSettings();
    }
    return false;
  }
  return false;
}
