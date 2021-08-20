import 'package:app/generated/protobuf/content/action.enum.pbenum.dart';
import 'package:app/generated/protobuf/content/action.model.pb.dart';
import 'package:app/generated/protobuf/content/action.service.pb.dart';
import 'package:app/generated/protobuf/content/content.service.pb.dart';
import 'package:app/global/controller.dart';
import 'package:app/routes/route.dart';
import 'package:app/service/action.dart';
import 'package:app/service/content.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:get/get.dart';
import 'package:fixnum/fixnum.dart';
import 'package:app/generated/protobuf/utils/request/param.pb.dart' as $param;

class ActionBar extends StatefulWidget {
  ActionBar(this.ext,{this.userAction}) : super();
  final UserAction? userAction;
  final ContentExt ext;
  @override
  ActionBarState createState() =>ActionBarState();
}

class ActionBarState extends State<ActionBar> {
  late ContentExt ext = widget.ext.toBuilder() as ContentExt;
  late UserAction? action = widget.userAction?.toBuilder() as UserAction?;
  static const size = 25.0;
  final ContentClient contentClient =Get.find();
  final ActionClient actionClient =Get.find();

  @override
  Widget build(BuildContext context) {
    return Padding(
        //上下各添加8像素补白
        padding: const EdgeInsets.symmetric(vertical: 10.0),
        child: Row(children: [
          Expanded(
            flex: 3,
            child: GestureDetector(
                onTap: () {
                  print('点击了');
                },
                child: Row(children: [
                  Expanded(
                      flex: 1,
                      child:
                          Icon(Icons.share, color: Colors.green, size: size)),
                  Expanded(
                    flex: 1,
                    child: Text(ext.shareCount.toStringUnsigned()),
                  )
                ])),
          ),
          Expanded(
            flex: 3,
            child: GestureDetector(
                onTap: () {
                  Get.toNamed(Routes.contentDetails(ext.type, ext.refId));
                },
                child: Row(
                  children: [
                    Expanded(
                        flex: 1,
                        child: FaIcon(FontAwesomeIcons.commentAlt, size: size)),
                    Expanded(
                      flex: 1,
                      child: Text(ext.commentCount.toStringUnsigned()),
                    )
                  ],
                )),
          ),
          Expanded(
            flex: 3,
            child: GestureDetector(
              onTap: () async{
                if(!await check()) return;
                  final favs = await contentClient.stub.favList(FavListReq());
                showBottomSheet(
                    context: context,
                    builder: (context) {
                      return Container(height: 200, color: Colors.lightBlue);
                    });

              },
              child: Row(
              children: [
                Expanded(
                    flex: 1,
                    child: Icon(Icons.star,
                        color: action!=null&&action!.collects.length>0? Colors.blueAccent[200]:Colors.white54, size: size)),
                Expanded(
                    flex: 1, child: Text(ext.collectCount.toStringUnsigned()))
              ],
            )),
          ),
          Expanded(
            flex: 3,
            child: GestureDetector(
                onTap: () async{
                  if(!await check()) return;
                  if(action!.likeId == 0){
                    final object = await actionClient.stub.like(LikeReq(type: ext.type,refId: ext.refId,action:ActionType.ActionLike));
                    action!.likeId = object.id;
                    ext.likeCount++;
                  }else{
                    await actionClient.stub.delLike($param.Object(id:action!.likeId));
                    action!.likeId = Int64(0);
                    ext.likeCount--;
                  }

                  setState(() {});
                },
                child: Row(
                children: [
                  Expanded(
                      flex: 1,
                      child: Icon(Icons.favorite, color: action!=null&&action!.likeId!=0? Colors.red :Colors.white54, size: size)),
                  Expanded(
                    flex: 1,
                    child: Text(ext.likeCount.toStringUnsigned()),
                  )
                ],
            )),
          ),
          Expanded(
            flex: 2,
            child: GestureDetector(
                child: Icon(Icons.more_horiz_outlined, size: size)),
          ),
        ]));
  }
  Future<bool> check() async{
    if (globalState.authState.userAuth ==null) {
      Get.toNamed(Routes.LOGIN);
      return false;
    }
    if(action == null){
      try{
        action = await actionClient.stub.getUserAction(ContentReq(type: ext.type,refId: ext.refId));
      }catch(e){
        return false;
      }
    }
    return true;
  }

}
