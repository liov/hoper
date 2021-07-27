
import 'package:app/generated/protobuf/user/user.service.pb.dart';
import 'package:app/pages/home/global/global_state/auth.dart';
import 'package:app/pages/home/global/global_state/global_controller.dart';
import 'package:app/model/user.dart';
import 'package:app/service/user.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import 'home/global/dao.dart';

class LoginView extends StatelessWidget {

  final GlobalController globalController =  Get.find();
  final UserClient userClient =  Get.find();
  final Dao dao =  Get.find();
  final _formKey = GlobalKey<FormState>();

  login(String account,password) async{
      final rep = await userClient.stub.login(LoginReq(input: account, password: password));
      print(rep.user);
      globalController.authState.user = rep.user;
      globalController.authState.key = rep.token;
      dao.box.put(AuthKey, rep.token);
      navigator!.pop();
  }

  @override
  Widget build(BuildContext context) {
    var _account = '';
    var _password = '';

    return Scaffold(
        resizeToAvoidBottomInset: false,
        body: Center(
          child: Container(
            padding: EdgeInsets.all(60.0),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'Welcome',
                    style: Theme.of(context).textTheme.headline2,
                  ),
                  TextFormField(
                    decoration: InputDecoration(
                      hintText: '邮箱/手机',
                    ),
                    onSaved: (value) {
                      _account = value!;
                    },
                  ),
                  TextFormField(
                    decoration: InputDecoration(
                      hintText: '密码',
                    ),
                    onSaved: (value) {
                      _password = value!;
                    },
                    obscureText: true,
                  ),
                  SizedBox(
                    height: 24,
                  ),
                  ElevatedButton(
                    style: ButtonStyle(
                        foregroundColor:ButtonStyleButton.allOrNull<Color>(Colors.yellow)
                    ),
                    child: Text('登录'),
                    onPressed: () {
                      _formKey.currentState!.save();
                      login(_account, _password);
                    },
                  ),
                ],
              ),
            ),
          ),
        ),
    );
  }
}
