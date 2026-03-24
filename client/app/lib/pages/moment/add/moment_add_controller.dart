import 'package:app/generated/protobuf/content/moment.service.pb.dart';
import 'package:app/global/state.dart';
import 'package:app/rpc/moment.dart';
import 'package:app/components/media/media.dart';
import 'package:get/get.dart';

import 'package:app/generated/protobuf/common/common.model.pbenum.dart';

class MomentAddController extends GetxController with MediaController {
  String content = '';
  final MomentGrpcClient momentClient = Get.find();

  Future<void> save() async {
    try {
      await momentClient.stub.add(
        AddMomentReq(
          type: MediaType.MediaTypeImage,
          content: content,
          images: imageUrls,
        ),
      );
      navigator!.pop();
    } catch (e) {
      globalService.logger.e(e);
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
