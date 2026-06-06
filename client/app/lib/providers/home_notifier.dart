import 'package:app/components/bottom/bottom.dart';
import 'package:app/util/dialog.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'home_notifier.g.dart';

class HomeState {
  const HomeState({this.selectedIndex = 0, this.pageIndex = 0});

  final int selectedIndex;
  final int pageIndex;

  HomeState copyWith({int? selectedIndex, int? pageIndex}) {
    return HomeState(
      selectedIndex: selectedIndex ?? this.selectedIndex,
      pageIndex: pageIndex ?? this.pageIndex,
    );
  }
}

@Riverpod(keepAlive: true)
class Home extends _$Home {
  var scrollNum = 0;

  static final bottomNavigationBarList = [
    Bottom.icon(Icons.home, label: "瞬间", pageIndex: 0),
    Bottom.icon(Icons.gamepad, label: "Profile", pageIndex: 1),
    Bottom.icon(Icons.add, onTap: () => toast("测试")),
    Bottom.icon(Icons.account_box_rounded, label: "Moments", pageIndex: 2),
    Bottom.fa(FontAwesomeIcons.user, label: "我的", pageIndex: 3),
  ];
  static const pageBottomIdx = [0, 1, 3, 4];

  @override
  HomeState build() => const HomeState();

  void onPageChanged(int index) {
    if (index < 0 || index >= pageBottomIdx.length) return;
    state = state.copyWith(selectedIndex: pageBottomIdx[index], pageIndex: index);
  }

  void onItemTapped(int index) {
    final bottom = bottomNavigationBarList[index];
    var pageIndex = state.pageIndex;
    if (bottom.pageIndex != null) {
      pageIndex = bottom.pageIndex!;
    }
    bottom.onTap?.call();
    state = state.copyWith(selectedIndex: index, pageIndex: pageIndex);
  }

  void continueScroll() {
    final index = state.pageIndex + (scrollNum < 0 ? -1 : 1);
    if (index < 0 || index >= pageBottomIdx.length) return;
    scrollNum = 0;
    state = state.copyWith(selectedIndex: pageBottomIdx[index], pageIndex: index);
  }
}
