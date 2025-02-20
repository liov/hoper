import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pb.dart';
import 'package:app/global/state/auth.dart';
import 'package:app/global/state.dart';
import 'package:app/routes/route.dart';
import 'package:app/utils/keyboard.dart';
import 'package:grpc/grpc.dart';
import 'package:app/service/user.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';

import 'login_controller.dart';

class LoginView extends StatelessWidget {
  final _formKey = GlobalKey<FormState>();

  final LoginController loginController = Get.put(LoginController());

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        resizeToAvoidBottomInset: false,
        body: GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () {
            hideKeyboard(context);
          },
          child: Center(
            child: Container(
              padding: EdgeInsets.all(60.0),
              child: Form(
                key: _formKey,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Obx(() {
                      if (loginController.mode == 1)
                        return _buildLoginView();
                      else
                        return _buildSignView();
                    }),
                    SizedBox(
                      height: 24,
                    ),
                    Row(
                      children: [
                        Expanded(
                          flex: 1,
                          child: Center(
                              child: ElevatedButton(
                                style: ButtonStyle(
                                    foregroundColor:
                                        ButtonStyleButton.allOrNull<Color>(
                                            Colors.yellow)),
                                child: Text('注册'),
                                onPressed: () {
                                  if (loginController.mode.value == 1) {
                                    loginController.mode.value = 2;
                                    return;
                                  }
                                  if (_formKey.currentState!.validate()) {
                                    _formKey.currentState!.save();
                                    loginController.signup();
                                  }
                                },
                          )),
                        ),
                        Expanded(
                            flex: 1,
                            child: Center(
                                child: ElevatedButton(
                              style: ButtonStyle(
                                  foregroundColor:
                                      ButtonStyleButton.allOrNull<Color>(
                                          Colors.yellow)),
                              child: Text('登录'),
                              onPressed: () {
                                if (loginController.mode.value == 2) {
                                  loginController.mode.value = 1;
                                  return;
                                }
                                _formKey.currentState!.save();
                                loginController.login();
                              },
                            )))
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),
        ));
  }

  Widget _buildLoginView() {
    return Column(mainAxisAlignment: MainAxisAlignment.center, children: [
      TextFormField(
        decoration: InputDecoration(
          labelText: '邮箱/手机',
          hintText: '邮箱/手机',
        ),
        initialValue: globalService.box.get(AuthState.StringAccountKey),
        onSaved: (value) {
          loginController.account = value!;
        },
      ),
      TextFormField(
        decoration: InputDecoration(
          labelText: '密码',
          hintText: '密码',
        ),
        onSaved: (value) {
          loginController.password = value!;
        },
        obscureText: true,
      ),
    ]);
  }

  Widget _buildSignView() {
    return Column(mainAxisAlignment: MainAxisAlignment.center, children: [
      TextFormField(
        decoration: InputDecoration(
          labelText: '邮箱',
          hintText: '邮箱',
        ),
        validator: (String? value) {
          return value != null ? null : '邮箱不能为空';
        },
        onSaved: (value) {
          loginController.mail = value!;
        },
      ),
      TextFormField(
        decoration: InputDecoration(
          labelText: '手机',
          hintText: '手机',
        ),
        initialValue: loginController.phone,
        validator: (value) {
          if (!RegExp(r'^1\d{10}$').hasMatch(value!)) {
            return '请输入正确手机号';
          }
          return null;
        },
        onSaved: (value) {
          loginController.phone = value!;
        },
      ),
      TextFormField(
        decoration: InputDecoration(
          labelText: '密码',
          hintText: '密码',
        ),
        onChanged: (value) {
          loginController.password = value;
        },
        validator: (String? value) {
          return value!.length > 5 ? null : '密码不小于6位';
        },
        obscureText: true,
      ),
      TextFormField(
        decoration: InputDecoration(
          labelText: '重复密码',
          hintText: '重复密码',
        ),
        validator: (String? value) {
          return value == loginController.password ? null : '密码输入不一致';
        },
        obscureText: true,
      ),
      TextFormField(
        decoration: InputDecoration(
          labelText: '昵称',
          hintText: '昵称',
        ),
        validator: (String? value) {
          return value!.length > 2 ? null : '昵称不小于3位';
        },
        onSaved: (value) {
          loginController.nickname = value!;
        },
      ),
      Obx(() => Row(
            children: <Widget>[
              Flexible(child: const Text('性别:')),
              Flexible(
                child: RadioListTile<Gender>(
                  title: const Text('男'),
                  value: Gender.GenderMale,
                  groupValue: loginController.gender.value,
                  onChanged: (value) {
                    loginController.gender.value = value!;
                  },
                ),
              ),
              Flexible(
                  child: RadioListTile<Gender>(
                title: const Text('女'),
                value: Gender.GenderFemale,
                groupValue: loginController.gender.value,
                onChanged: (value) {
                  loginController.gender.value = value!;
                },
              )),
            ],
          )),
      //_birthdayPicker()
    ]);
  }

  Widget _birthdayPicker() {
    return Row(children: <Widget>[
      const Flexible(flex: 1, child: Text('生日:')),
      Flexible(
          flex: 5,
          child: GestureDetector(
              onTap: () {
                Get.dialog(Center(
                    child: Column(
                  children: [
                    Expanded(
                      flex: 5,
                      child: CupertinoDatePicker(
                        mode: CupertinoDatePickerMode.date,
                        initialDateTime: loginController.birthDate,
                        minimumDate: DateTime(1950),
                        maximumDate: DateTime.now(),
                        onDateTimeChanged: (DateTime value) {
                          loginController.birthDate = value;
                        },
                        backgroundColor: Get.theme.scaffoldBackgroundColor,
                      ),
                    ),
                    Expanded(
                        flex: 1,
                        child: Center(
                            child: ElevatedButton(
                          style: ButtonStyle(
                              foregroundColor:
                                  ButtonStyleButton.allOrNull<Color>(
                                      Colors.yellow)),
                          child: Text('确定'),
                          onPressed: () => navigator!.pop(),
                        )))
                  ],
                )));
              },
              child: Center(
                child: Text(DateFormat('yyyy年MM月dd日')
                    .format(loginController.birthDate)),
              )))
    ]);
  }
}
