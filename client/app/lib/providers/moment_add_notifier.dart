import 'package:app/gen/pb/common/common.model.pbenum.dart';
import 'package:app/gen/pb/content/moment.service.pb.dart';
import 'package:app/global/service.dart';
import 'package:app/providers/media_actions.dart';
import 'package:app/util/nav.dart';
import 'package:image_picker/image_picker.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'moment_add_notifier.g.dart';

class MomentAddState {
  const MomentAddState({this.content = '', this.revision = 0});

  final String content;
  final int revision;

  MomentAddState copyWith({String? content, int? revision}) {
    return MomentAddState(content: content ?? this.content, revision: revision ?? this.revision);
  }
}

@riverpod
class MomentAdd extends _$MomentAdd {
  final media = MediaActions();

  @override
  MomentAddState build() {
    ref.onDispose(() => media.disposeVideoController());
    return const MomentAddState();
  }

  List<dynamic> get imageFiles => media.imageFiles;
  List<String> get imageUrls => media.imageUrls;

  set content(String value) => state = state.copyWith(content: value);

  Future<void> pickImage(ImageSource source, {bool isCamera = false}) async {
    await media.onImageButtonPressed(source, isCamera: isCamera);
    state = state.copyWith(revision: state.revision + 1);
  }

  Future<void> save() async {
    try {
      await globalService.momentClient.stub.add(
        AddMomentReq(
          type: MediaType.MediaTypeImage,
          content: state.content,
          images: media.imageUrls,
        ),
      );
      AppNavigator.pop();
    } catch (e) {
      globalService.logger.severe(e);
    }
  }
}
