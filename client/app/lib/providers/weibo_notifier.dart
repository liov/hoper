import 'package:app/global/service.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'weibo_notifier.g.dart';

class WeiboState {
  const WeiboState({
    this.userId = 0,
    this.page = 1,
    this.feature = 0,
    this.sinceId = '',
    this.list = const [],
    this.picWidth = 300,
    this.picHeight = 300,
    this.isEnd = false,
    this.isLoading = false,
  });

  final int userId;
  final int page;
  final int feature;
  final String sinceId;
  final List<String> list;
  final int picWidth;
  final int picHeight;
  final bool isEnd;
  final bool isLoading;

  WeiboState copyWith({
    int? userId,
    int? page,
    int? feature,
    String? sinceId,
    List<String>? list,
    int? picWidth,
    int? picHeight,
    bool? isEnd,
    bool? isLoading,
  }) {
    return WeiboState(
      userId: userId ?? this.userId,
      page: page ?? this.page,
      feature: feature ?? this.feature,
      sinceId: sinceId ?? this.sinceId,
      list: list ?? this.list,
      picWidth: picWidth ?? this.picWidth,
      picHeight: picHeight ?? this.picHeight,
      isEnd: isEnd ?? this.isEnd,
      isLoading: isLoading ?? this.isLoading,
    );
  }
}

@riverpod
class Weibo extends _$Weibo {
  @override
  WeiboState build() => const WeiboState();

  Future<void> newList(int userId) async {
    state = state.copyWith(userId: userId, page: 1, isEnd: false, isLoading: false, list: []);
    return getList();
  }

  Future<void> getList() async {
    globalService.logger.fine('getList');
    if (state.isEnd || state.isLoading) return;
    state = state.copyWith(isLoading: true);
    try {
      final response = await globalService.weiboClient.getOriginalList(
        uid: state.userId,
        page: state.page,
        feature: state.feature,
        sinceId: state.sinceId,
      );
      if (response == null || response.list.isEmpty) {
        state = state.copyWith(isEnd: true, isLoading: false);
        return;
      }
      final urls = [...state.list];
      for (final e in response.list) {
        if (e.picInfos != null) {
          urls.addAll(e.picInfos!.values.map((v) => v.mw2000.url));
        }
      }
      globalService.logger.fine('${response.list.length} ${urls.length}');
      state = state.copyWith(list: urls, page: state.page + 1, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false);
      rethrow;
    }
  }
}
