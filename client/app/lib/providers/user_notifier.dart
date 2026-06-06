import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'user_notifier.g.dart';

@riverpod
class User extends _$User {
  @override
  int build() => 0;

  void touch() => state++;
}

@riverpod
class Follow extends _$Follow {
  @override
  int build() => 0;

  void touch() => state++;
}
