import 'package:app/components/async/async.dart';
import 'package:app/generated/protobuf/content/content.model.pb.dart' as $moment;

import 'package:app/global/state.dart';

import 'package:app/pages/moment/item/moment_item_view.dart';
import 'package:app/pages/moment/detail/moment_detail_view.dart';

import 'package:flutter/material.dart';


class SampleView extends StatefulWidget {
  SampleView({this.tag = "default"}) : super();

  final String tag;

  // State 生命周期的起点，Flutter 会通过调用 StatefulWidget.createState() 来创建一个 State。可以通过构造方法，来接收父 Widget 传递的初始化 UI 配置数据，而这些配置数据，决定了 Widget 最初的呈现状态
  @override
  _SampleState createState() => _SampleState();
}

class _SampleState extends State<SampleView> with AutomaticKeepAliveClientMixin,WidgetsBindingObserver {

  var times = 0;
  // 在 State 对象被插入视图树时调用。在 State 的生命周期中只会被调用一次，因此可以在 initState 函数中做一些初始化操作
  @override
  void initState() {
    super.initState();
    times = 0;
  }

  // 专门用来处理 State 对象依赖关系变化，会在 initState() 调用结束后调用
  // State 对象的依赖关系发生变化后，Flutter 会回调该方法，随后触发组件构建。State 对象依赖关系发生变化的典型场景：系统语言 Locale 或应用主题改变时，系统会通知 State 执行 didChangeDependencies 回调方法
  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
  }


  Future<void> addTimes() async {
    // 当状态数据发生变化时，可以通过调用 setState 方法告诉 Flutter 使用更新后数据重建 UI
    setState(() {times++;});
  }

  // 构建视图。经过构造方法、initState、didChangeDependencies 后，Framework 认为 State 已经准备就绪，于是便调用 build。在 build 中，需要根据父 Widget 传递过来的初始化配置数据及 State 的当前状态，创建一个 Widget 然后返回
  @override
  Widget build(BuildContext context) {
    super.build(context);
    final _future = addTimes();
    return FutureBuilder<void>(
        future: _future,
        builder: (BuildContext context, AsyncSnapshot<void> snapshot) {
          return snapshot.handle() ??  Center(child: Text('$times'),);
  });
  }


  //在widget重新构建时，Flutter framework会调用Widget.canUpdate来检测Widget树中同一位置的新旧节点，然后决定是否需要更新，如果Widget.canUpdate返回true则会调用此回调。正如之前所述，Widget.canUpdate会在新旧widget的key和runtimeType同时相等时会返回true，也就是说在在新旧widget的key和runtimeType同时相等时didUpdateWidget()就会被调用。
  @override
  void didUpdateWidget(covariant SampleView oldWidget) {
    super.didUpdateWidget(oldWidget);
    times = 0;
  }

  // 此回调是专门为了开发调试而提供的，在热重载(hot reload)时会被调用，此回调在Release模式下永远不会被调用。
  @override
  void reassemble() {
    super.reassemble();
  }

  // 当State对象从树中被移除时，会调用此回调。在一些场景下，Flutter framework会将State对象重新插到树中，如包含此State对象的子树在树的一个位置移动到另一个位置时（可以通过GlobalKey来实现）。如果移除后没有重新插入到树中则紧接着会调用dispose()方法。
  @override
  void deactivate() {
    super.deactivate();
  }

  @override
  void activate() {
    super.activate();
  }

  // 当 State 被永久地从视图树中移除时，Flutter 会调用 dispose 方法，而一旦 dispose 方法被调用，组件就要被销毁了，因此可以在 dispose 方法中进行最终的资源释放、移除监听、清理环境等工作
  @override
  void dispose() {
    super.dispose();
  }

  @override
  bool get wantKeepAlive => true;
}