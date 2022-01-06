

import 'package:app/generated/protobuf/user/user.model.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pb.dart';
import 'package:app/generated/protobuf/empty/empty.pb.dart';
import 'package:app/generated/protobuf/request/param.pb.dart' as request;
import 'package:app/global/global_state.dart';
import 'package:app/global/service.dart';
import 'package:app/model/const/const.dart';
import 'package:app/pages/home/home_controller.dart';
import 'package:app/pages/user/login_view.dart';
import 'package:app/utils/dialog.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get_core/src/get_main.dart';
import 'package:get/get_instance/src/extension_instance.dart';
import 'package:get/get_navigation/src/extension_navigation.dart';
import 'package:grpc/grpc.dart';


class AuthState {
  UserAuthInfo? userAuth = null;
  UserBaseInfo? get userBaseInfo => self!=null?UserBaseInfo(id:self!.id,name:self!.name,gender: self!.gender,avatarUrl: self!.avatarUrl):null;
  User? self = null;

  static const _PRE = "AuthState";
  static const StringAuthKey = _PRE+Authorization;
  static const StringAccountKey = _PRE+"AccountKey";
  static const StringAuthInfoKey = _PRE+"AuthInfoKey";

  set account (String account)=> globalService.box.put(AuthState.StringAccountKey, account);
  String get account => globalService.box.get(AuthState.StringAccountKey);

  Future<void> getAuth() async {
    if (userAuth != null) return;
    final authKey = globalService.box.get(StringAuthKey);
    final authInfo = globalService.box.get(StringAuthInfoKey);
    if (authKey != null) {
      try {
        final user = await globalService.userClient.stub.authInfo(Empty(),options:CallOptions(metadata: {Authorization: authKey}));
        if (user.id == 0) return;
        this.userAuth = user;
        setAuth(authKey);
        getSelf().then((value) => globalState.userState.append(userBaseInfo));
        return null;
      } catch (err) {
        print(err);
      }
    }
  }

  Future<void> getSelf() async {
    if (self != null) return;
    if (userAuth == null) {
     await getAuth();
    }
    try {
      final user = await globalService.userClient.stub.info(request.Object());
      if (user.user.id == 0) return;
      this.self = user.user;
      return null;
    } catch (err) {
      print(err);
    }
  }

  void setAuth(String authKey) {
    globalService.httpClient.options.headers[Authorization] = authKey;
    globalService.subject.setState(CallOptions(metadata: {Authorization: authKey},timeout: Duration(seconds: 5)));
    globalService.box.put(AuthState.StringAuthKey, authKey);
  }

  Future<void> login(String account,String password) async{
    try{
      final rep = await globalService.userClient.stub.login(LoginReq(input: account, password: password,vCode: 'super'));
      final user = rep.user;
      self = rep.user;
      userAuth = UserAuthInfo(id:user.id,name:user.name,role:user.role,status:user.status);
      setAuth(rep.token);
      this.account = account;
      navigator!.pop();
      //Get.forceAppUpdate();
      Get.rootController.restartApp();
    } on GrpcError catch (e) {
      toast(e.message!);
    }catch (e) {
      // No specified type, handles all
      print('Something really unknown: $e');
    }
  }

   Future<void> logout() async{
    userAuth = null;
    globalService.httpClient.options.headers.remove(Authorization);
    globalService.box.delete(AuthState.StringAuthKey);
    self = null;
    try{
      await globalService.userClient.stub.logout(Empty());
      globalService.subject.setState(CallOptions(timeout: Duration(seconds: 5)));
    } on GrpcError catch (e) {
      toast(e.message!);
    }catch (e) {
      // No specified type, handles all
      print('Something really unknown: $e');
    }
    //Get.forceAppUpdate();
    Get.rootController.restartApp();
  }

  void test(void test()){
    test();
  }
}


