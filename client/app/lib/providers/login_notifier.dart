import 'package:app/gen/pb/user/user.model.pbenum.dart';
import 'package:app/gen/pb/user/user.service.pb.dart';
import 'package:app/global/state.dart';
import 'package:app/util/dialog.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'login_notifier.g.dart';

class LoginState {
  const LoginState({this.mode = 1, this.gender = Gender.GenderPlaceholder, this.birthDate});

  final int mode;
  final Gender gender;
  final DateTime? birthDate;

  LoginState copyWith({int? mode, Gender? gender, DateTime? birthDate}) {
    return LoginState(
      mode: mode ?? this.mode,
      gender: gender ?? this.gender,
      birthDate: birthDate ?? this.birthDate,
    );
  }
}

@riverpod
class Login extends _$Login {
  String? phone;
  String? mail;
  String? countryCallingCode;
  String? account;
  String? password;
  String? nickname;

  @override
  LoginState build() => LoginState(birthDate: DateTime(2020, 1, 1));

  DateTime get birthDate => state.birthDate ?? DateTime(2020, 1, 1);

  set birthDate(DateTime value) => state = state.copyWith(birthDate: value);

  void setMode(int mode) => state = state.copyWith(mode: mode);

  void setGender(Gender gender) => state = state.copyWith(gender: gender);

  Future<void> login() async {
    return globalState.authState.login(countryCallingCode, account!, password!);
  }

  Future<void> signup() async {
    try {
      final resp = await globalService.userClient.stub.signup(SignupReq(
        name: nickname,
        gender: state.gender,
        password: password,
        mail: mail,
        phone: phone,
        vCode: 'super',
      ));
      if (resp.value != "") toast(resp.value);
      globalState.authState.account = mail!;
      setMode(1);
    } catch (e) {
      globalService.logger.warning('$e');
    }
  }
}
