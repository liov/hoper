import 'package:app/components/async/async.dart';
import 'package:app/global/global_state.dart';
import 'package:app/global/const.dart';
import 'package:app/pages/user/user_controller.dart';
import 'package:app/routes/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:getwidget/components/avatar/gf_avatar.dart';
import 'package:getwidget/components/list_tile/gf_list_tile.dart';


class UserView extends StatelessWidget {

final UserController userController = Get.put(UserController());

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Center(
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
               return GetBuilder<UserController>(builder: (_) {
                 if(globalState.authState.self == null){
                   return _buildNoLogin();
                 }else return Column(
                   children: [
                     _buildHeader(),
                     _buildSignOut()
                   ],
                 );
               },
               );
            }
          },
        ),
      ),
    );
  }
  Widget _buildNoLogin(){
    return  ElevatedButton(
        style: ElevatedButton.styleFrom(
          minimumSize: const Size(100, 45),
          side: BorderSide(),
          shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(4))),
        ),
        child: const Text('立即登录'),
        onPressed: () {
          Get.toNamed(Routes.LOGIN);
        },
      );
    }
  Widget _buildSignOut(){
    return  ElevatedButton(
      style: ElevatedButton.styleFrom(
        minimumSize: const Size(100, 45),
        side: BorderSide(),
        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(4))),
      ),
      child: const Text('退出登录'),
      onPressed: () {
        globalState.authState.logout();
      },
    );
  }

  Widget _buildHeader(){
    return GFListTile(
        avatar:CircleAvatar(
          child:ExtendedImage.network(
            BASE_STATIC_URL + globalState.authState.self!.avatar,
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
