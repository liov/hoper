import 'package:app/global/service.dart';

extension FutureSafe<T> on Future<T> {
  Future<T?> safe() => then<T?>((v) => v).catchError((Object e, StackTrace s) {
        globalService.logger.severe('exception', e, s);
        return null;
      });
}
