import 'package:app/global/state.dart';
import 'package:app/global/const.dart';
import 'package:app/pages/route.dart';
import 'package:extended_image/extended_image.dart';
import 'package:flutter/material.dart';
import 'package:app/util/nav.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:getwidget/components/list_tile/gf_list_tile.dart';

class UserView extends ConsumerWidget {
  const UserView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return SafeArea(
      child: Center(
        child: FutureBuilder(
          future: globalState.authState.getSelf(),
          builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
            switch (snapshot.connectionState) {
              case ConnectionState.none:
              case ConnectionState.active:
                return const Text('ConnectionState.active');
              case ConnectionState.waiting:
                return const Center(
                  child: CircularProgressIndicator(),
                );
              case ConnectionState.done:
                if (globalState.authState.self == null) {
                  return _buildNoLogin();
                }
                return Column(
                  children: [
                    _buildHeader(),
                    _buildSignOut(),
                  ],
                );
            }
          },
        ),
      ),
    );
  }

  Widget _buildNoLogin() {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        minimumSize: const Size(100, 45),
        side: const BorderSide(),
        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(4))),
      ),
      child: const Text('立即登录'),
      onPressed: () {
        AppNavigator.pushNamed(Routes.LOGIN);
      },
    );
  }

  Widget _buildSignOut() {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        minimumSize: const Size(100, 45),
        side: const BorderSide(),
        shape: const RoundedRectangleBorder(borderRadius: BorderRadius.all(Radius.circular(4))),
      ),
      child: const Text('退出登录'),
      onPressed: () {
        globalState.authState.logout();
      },
    );
  }

  Widget _buildHeader() {
    return GFListTile(
      avatar: CircleAvatar(
        child: ExtendedImage.network(
          BASE_STATIC_URL + globalState.authState.self!.avatar,
          alignment: Alignment.centerLeft,
          fit: BoxFit.fill,
          shape: BoxShape.circle,
          cache: true,
        ),
      ),
      titleText: globalState.authState.self!.name,
      subTitleText: globalState.authState.self!.signature,
      icon: const Icon(Icons.qr_code),
    );
  }
}
