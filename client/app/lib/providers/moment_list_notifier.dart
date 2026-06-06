import 'package:app/gen/pb/content/moment.model.pb.dart';
import 'package:app/gen/pb/content/moment.service.pb.dart';
import 'package:app/global/state.dart';
import 'package:app/rpc/moment.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'moment_list_notifier.g.dart';

class MomentTagList {
  MomentTagList() : req = MomentListReq(pageNo: 1, pageSize: 10);

  final MomentListReq req;
  final list = List<Moment>.empty(growable: true);
  var revision = 0;
  var loading = false;
  var exhausted = false;

  void resetList() {
    list.clear();
    req.pageNo = 1;
    exhausted = false;
    revision++;
  }

  Future<bool> grpcGetList(MomentGrpcClient momentClient) async {
    if (loading || exhausted) return false;
    loading = true;
    try {
      final response = await momentClient.stub.list(req);
      if (response.list.isEmpty) {
        exhausted = true;
        return false;
      }
      globalState.userState.appendUsers(response.users);
      list.addAll(response.list);
      req.pageNo++;
      revision++;
      return true;
    } finally {
      loading = false;
    }
  }
}

class MomentListState {
  const MomentListState({this.tags = const {}});

  final Map<String, MomentTagList> tags;

  MomentListState copyWith({Map<String, MomentTagList>? tags}) {
    return MomentListState(tags: tags ?? this.tags);
  }
}

@Riverpod(keepAlive: true)
class MomentList extends _$MomentList {
  @override
  MomentListState build() => const MomentListState();

  MomentTagList? tagList(String tag) => state.tags[tag];

  Future<void> newList(String tag) async {
    if (state.tags.containsKey(tag)) return;
    final list = MomentTagList();
    await list.grpcGetList(globalService.momentClient);
    final next = Map<String, MomentTagList>.from(state.tags)..[tag] = list;
    state = state.copyWith(tags: next);
  }

  Future<void> pullList(String tag) async {
    final tagList = state.tags[tag];
    if (tagList == null || tagList.loading || tagList.exhausted) return;
    final changed = await tagList.grpcGetList(globalService.momentClient);
    if (changed) {
      state = state.copyWith(tags: Map<String, MomentTagList>.from(state.tags));
    }
  }

  Future<void> resetList(String tag) async {
    state.tags[tag]?.resetList();
    await pullList(tag);
  }
}
