import 'package:app/gen/pb/user/user.model.pb.dart';
import 'package:app/global/service.dart';
import 'package:app/global/state/auth.dart';
import 'package:app/util/keyboard.dart';
import 'package:app/providers/providers.dart';
import 'package:app/providers/login_notifier.dart';
import 'package:flutter/cupertino.dart';
import 'package:app/util/nav.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class LoginView extends ConsumerWidget {
  LoginView({super.key});

  final _formKey = GlobalKey<FormState>();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final loginState = ref.watch(loginProvider);
    final login = ref.read(loginProvider.notifier);
    return Scaffold(
        resizeToAvoidBottomInset: false,
        body: GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () {
            hideKeyboard(context);
          },
          child: Center(
            child: Container(
              padding: const EdgeInsets.all(60.0),
              child: Form(
                key: _formKey,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    loginState.mode == 1 ? _buildLoginView(login) : _buildSignView(context, loginState, login),
                    const SizedBox(height: 24),
                    Row(
                      children: [
                        Expanded(
                          flex: 1,
                          child: Center(
                              child: ElevatedButton(
                            style: ButtonStyle(
                                foregroundColor: ButtonStyleButton.allOrNull<Color>(Colors.yellow)),
                            child: const Text('注册'),
                            onPressed: () {
                              if (loginState.mode == 1) {
                                login.setMode(2);
                                return;
                              }
                              if (_formKey.currentState!.validate()) {
                                _formKey.currentState!.save();
                                login.signup();
                              }
                            },
                          )),
                        ),
                        Expanded(
                            flex: 1,
                            child: Center(
                                child: ElevatedButton(
                              style: ButtonStyle(
                                  foregroundColor: ButtonStyleButton.allOrNull<Color>(Colors.yellow)),
                              child: const Text('登录'),
                              onPressed: () {
                                if (loginState.mode == 2) {
                                  login.setMode(1);
                                  return;
                                }
                                _formKey.currentState!.save();
                                login.login();
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

  Widget _buildLoginView(Login login) {
    return Column(mainAxisAlignment: MainAxisAlignment.center, children: [
      TextFormField(
        decoration: const InputDecoration(
          labelText: '邮箱/手机',
          hintText: '邮箱/手机',
        ),
        initialValue: globalService.box.get(AuthState.StringAccountKey),
        onSaved: (value) {
          login.account = value!;
        },
      ),
      TextFormField(
        decoration: const InputDecoration(
          labelText: '密码',
          hintText: '密码',
        ),
        onSaved: (value) {
          login.password = value!;
        },
        obscureText: true,
      ),
    ]);
  }

  Widget _buildSignView(BuildContext context, LoginState loginState, Login login) {
    return Column(mainAxisAlignment: MainAxisAlignment.center, children: [
      TextFormField(
        decoration: const InputDecoration(
          labelText: '邮箱',
          hintText: '邮箱',
        ),
        validator: (String? value) {
          return value != null ? null : '邮箱不能为空';
        },
        onSaved: (value) {
          login.mail = value!;
        },
      ),
      TextFormField(
        decoration: const InputDecoration(
          labelText: '手机',
          hintText: '手机',
        ),
        initialValue: login.phone,
        validator: (value) {
          if (!RegExp(r'^1\d{10}$').hasMatch(value!)) {
            return '请输入正确手机号';
          }
          return null;
        },
        onSaved: (value) {
          login.phone = value!;
        },
      ),
      TextFormField(
        decoration: const InputDecoration(
          labelText: '密码',
          hintText: '密码',
        ),
        onChanged: (value) {
          login.password = value;
        },
        validator: (String? value) {
          return value!.length > 5 ? null : '密码不小于6位';
        },
        obscureText: true,
      ),
      TextFormField(
        decoration: const InputDecoration(
          labelText: '重复密码',
          hintText: '重复密码',
        ),
        validator: (String? value) {
          return value == login.password ? null : '密码输入不一致';
        },
        obscureText: true,
      ),
      TextFormField(
        decoration: const InputDecoration(
          labelText: '昵称',
          hintText: '昵称',
        ),
        validator: (String? value) {
          return value!.length > 2 ? null : '昵称不小于3位';
        },
        onSaved: (value) {
          login.nickname = value!;
        },
      ),
      Row(
        children: <Widget>[
          const Flexible(child: Text('性别:')),
          Flexible(
            child: RadioListTile<Gender>(
              title: const Text('男'),
              value: Gender.GenderMale,
              groupValue: loginState.gender,
              onChanged: (value) {
                login.setGender(value!);
              },
            ),
          ),
          Flexible(
            child: RadioListTile<Gender>(
              title: const Text('女'),
              value: Gender.GenderFemale,
              groupValue: loginState.gender,
              onChanged: (value) {
                login.setGender(value!);
              },
            ),
          ),
        ],
      ),
    ]);
  }

  Widget _birthdayPicker(BuildContext context, Login login) {
    return Row(children: <Widget>[
      const Flexible(flex: 1, child: Text('生日:')),
      Flexible(
          flex: 5,
          child: GestureDetector(
              onTap: () {
                AppNavigator.dialog(Center(
                    child: Column(
                  children: [
                    Expanded(
                      flex: 5,
                      child: CupertinoDatePicker(
                        mode: CupertinoDatePickerMode.date,
                        initialDateTime: login.birthDate,
                        minimumDate: DateTime(1950),
                        maximumDate: DateTime.now(),
                        onDateTimeChanged: (DateTime value) {
                          login.birthDate = value;
                        },
                        backgroundColor: Theme.of(context).scaffoldBackgroundColor,
                      ),
                    ),
                    Expanded(
                        flex: 1,
                        child: Center(
                            child: ElevatedButton(
                          style: ButtonStyle(
                              foregroundColor: ButtonStyleButton.allOrNull<Color>(Colors.yellow)),
                          child: const Text('确定'),
                          onPressed: () => AppNavigator.pop(),
                        )))
                  ],
                )));
              },
              child: Center(
                child: Text(DateFormat('yyyy年MM月dd日').format(login.birthDate)),
              )))
    ]);
  }
}
