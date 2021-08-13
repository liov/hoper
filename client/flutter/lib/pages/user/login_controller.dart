import 'package:app/generated/protobuf/user/user.enum.pb.dart';
import 'package:app/generated/protobuf/user/user.service.pb.dart';
import 'package:app/global/controller.dart';
import 'package:app/utils/dialog.dart';
import 'package:get/get.dart';

class LoginController extends GetxController {

  var mode = 1.obs; // 1登录，2注册

  String? phone;
  String? mail ;
  String? account;
  String? password;
  var gender = Gender.GenderUnfilled.obs;
  String? nickname;
  DateTime birthDate = DateTime(2020,1,1,);


  Future<void> login() async{
    return  globalState.authState.login(account!, password!);
  }

  Future<void> signup() async{
    try {
    final resp =  await globalService.userClient.stub.signup(
          SignupReq(name: nickname,gender: gender.value, password: password,mail: mail,phone: phone,vCode: 'super'));
    if(resp.value!="") dialog(resp.value);
      globalState.authState.account = mail!;
      mode.value = 1;
    }catch(e){
      print(e);
    }
  }


  @override
  void onReady() {
    // TODO: implement onReady
    super.onReady();
  }

  @override
  void onClose() {
    // TODO: implement onClose
    super.onClose();
  }
}
