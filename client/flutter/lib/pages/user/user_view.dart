import 'package:app/components/async/async.dart';
import 'package:app/global/controller.dart';
import 'package:app/model/const/const.dart';
import 'package:app/routes/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:getwidget/components/avatar/gf_avatar.dart';
import 'package:getwidget/components/list_tile/gf_list_tile.dart';


class UserView extends StatelessWidget {


  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Container(
        child: FutureBuilder(
          future: globalState.authState.getSelf(),
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            switch (snapshot.connectionState) {
              case ConnectionState.none:
              case ConnectionState.active:
                return Text('ConnectionState.active');
              case ConnectionState.waiting:
                return Center(
                  child: CircularProgressIndicator(),
                );
              case ConnectionState.done:
                if(globalState.authState.self == null){
                  return _buildNoLogin();
                }else return Column(
                  children: [
                    _buildHeader()
                  ],
                );
            }
          },
        ),
      ),
    );
  }
  Widget _buildNoLogin(){
      return GestureDetector(
        onTap: (){
          Get.toNamed(Routes.LOGIN);
        },
        child: Text('立即登录'),
      );
    }


  Widget _buildHeader(){
    if (globalState.authState.self==null){
      return GestureDetector(
        onTap: (){
          Get.toNamed(Routes.LOGIN);
        },
        child: Text('立即登录'),
      );
    }
    return GFListTile(
        avatar:CircleAvatar(
          child:ExtendedImage.network(
            BASE_STATIC_URL + globalState.authState.self!.avatarUrl,
            alignment: Alignment.centerLeft,
            fit: BoxFit.fill,
            shape: BoxShape.circle,
            cache: true,
          ),
        ),
        titleText:globalState.authState.self!.name,
        subTitleText:globalState.authState.self!.signature,
        icon: Icon(Icons.qr_code)
    );
  }
}
