import 'package:app/util/image_file.dart';
import 'package:app/global/service.dart';
import 'package:app/pages/image/slide_image.dart';
import 'package:app/providers/providers.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';

class MomentAddView extends ConsumerWidget {
  MomentAddView({super.key});

  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.watch(momentAddProvider);
    final add = ref.read(momentAddProvider.notifier);
    globalService.logger.fine('@'.codeUnits);
    return Scaffold(
        appBar: AppBar(actions: [
          IconButton(
            icon: const Text('保存'),
            onPressed: () {
              _formKey.currentState!.save();
              add.save();
            },
          )
        ]),
        body: Column(
          children: [
            Padding(
                padding: const EdgeInsets.all(16.0),
                child: Form(
                    key: _formKey,
                    child: Column(children: <Widget>[
                      TextFormField(
                        minLines: 5,
                        maxLines: 10,
                        decoration: const InputDecoration(
                          hintText: '记录这一刻,晒给懂你的人',
                        ),
                        onSaved: (value) {
                          add.content = value!;
                        },
                      ),
                    ]))),
            Builder(
              builder: (context) {
                if (add.imageFiles.isEmpty) {
                  return Container();
                }
                final images = add.imageFiles;
                return GridView.builder(
                    gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                        crossAxisCount: 3, childAspectRatio: 1.0),
                    shrinkWrap: true,
                    itemCount: images.length,
                    itemBuilder: (BuildContext context, int index) {
                      return GestureDetector(
                        child: extendedImageFile(images[index].path,
                            alignment: Alignment.centerLeft, fit: BoxFit.fill),
                        onTap: () => slideImageRoute(images[index].path),
                      );
                    });
              },
            ),
          ],
        ),
        bottomSheet: const MomentAddBottomSheet());
  }
}

class MomentAddBottomSheet extends ConsumerWidget {
  const MomentAddBottomSheet({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final add = ref.read(momentAddProvider.notifier);
    return Row(
      textBaseline: TextBaseline.alphabetic,
      children: [
        Expanded(
            flex: 1,
            child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                child: const Icon(Icons.camera_alt),
                onTap: () {
                  add.pickImage(ImageSource.camera, isCamera: true);
                })),
        Expanded(
            flex: 1,
            child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                child: const Icon(Icons.photo),
                onTap: () {
                  add.pickImage(ImageSource.gallery);
                })),
        Expanded(
            flex: 1,
            child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                child: const Icon(Icons.alternate_email),
                onTap: () {
                  add.pickImage(ImageSource.gallery);
                })),
        Expanded(
            flex: 1,
            child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                child: const FaIcon(FontAwesomeIcons.hashtag),
                onTap: () {
                  add.pickImage(ImageSource.gallery);
                })),
        Expanded(
            flex: 1,
            child: GestureDetector(
                behavior: HitTestBehavior.opaque,
                child: const Icon(Icons.mood),
                onTap: () {
                  add.pickImage(ImageSource.gallery);
                })),
      ],
    );
  }
}
