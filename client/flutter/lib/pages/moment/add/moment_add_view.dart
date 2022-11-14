import 'dart:io';


import 'package:app/global/global_service.dart';
import 'package:app/pages/photo/slide_photo.dart';
import 'package:app/components/media/media.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';
import 'package:video_player/video_player.dart';
import 'package:app/pages/moment/add/moment_add_controller.dart';

class MomentAddView extends StatelessWidget {
  final MomentAddController controller = Get.find();
  final _formKey = GlobalKey<FormState>();

  MomentAddView({super.key});


  @override
  Widget build(BuildContext context) {
    globalService.logger.d('@'.codeUnits);
    return Scaffold(
        appBar: AppBar(
            actions: [
              IconButton(icon: const Text('保存'),onPressed: (){
                _formKey.currentState!.save();
                controller.save();
              },)
            ]
        ),
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
                        onSaved: (value){
                          controller.content = value!;
                        },
                      ),
                    ]))),
            GetBuilder<MomentAddController>(builder: (_) {
              if(controller.imageFiles.isEmpty){
                return Container();
              }
              final images = controller.imageFiles;
             return GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                      crossAxisCount: 3, //横轴三个子widget
                      childAspectRatio: 1.0 //宽高比为1时，子widget
                  ),
                  shrinkWrap: true,
                  itemCount:images.length,
                  itemBuilder: (BuildContext context, int index) {
                    return  GestureDetector(
                      child:ExtendedImage.file(
                        File(images[index].path),
                        alignment: Alignment.centerLeft,
                        fit: BoxFit.fill,
                        //cancelToken: cancellationToken,
                      ),
                      onTap:()=>slidePhotoRoute(images[index].path),
                    );});
            },),
          ],
        ),
      bottomSheet: const MomentAddBottomSheet(),
    );
  }
}


class MomentAddBottomSheet extends StatelessWidget{
  const MomentAddBottomSheet({Key? key}):super(key: key);

  @override
  Widget build(BuildContext context) {
    final MediaController controller = MediaController();
   return Row(
     textBaseline: TextBaseline.alphabetic,
     children: [
       Expanded(flex: 1,child:  GestureDetector(
           behavior:HitTestBehavior.opaque,
           child: const Icon(Icons.camera_alt),
           onTap:(){
             controller.onImageButtonPressed(ImageSource.camera,isCamera: true);
           }),),
       Expanded(flex: 1,
           child: GestureDetector(
               behavior:HitTestBehavior.opaque,
               child: const Icon(Icons.photo),
               onTap:(){
                 controller.onImageButtonPressed(ImageSource.gallery);
               }
           )
       ),
       Expanded(flex: 1,child: GestureDetector(
           behavior:HitTestBehavior.opaque,
           child: const Icon(Icons.alternate_email),
           onTap:(){
             controller.onImageButtonPressed(ImageSource.gallery);
           }
       ),),
       Expanded(flex: 1,child: GestureDetector(
           behavior:HitTestBehavior.opaque,
           child: const Icon(FontAwesomeIcons.hashtag),
           onTap:(){
             controller.onImageButtonPressed(ImageSource.gallery);
           }
       ),),
       Expanded(flex: 1,child: GestureDetector(
           behavior:HitTestBehavior.opaque,
           child: const Icon(Icons.mood),
           onTap:(){
             controller.onImageButtonPressed(ImageSource.gallery);
           }
       ),),
     ],
   );
  }
}


