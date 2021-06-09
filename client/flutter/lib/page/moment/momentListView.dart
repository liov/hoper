import 'package:app/model/moment.dart';
import 'package:app/model/user.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';

import 'package:flutter_markdown/flutter_markdown.dart';

class MomentListView extends StatefulWidget {
  late final List<Moment> list = List.empty(growable: true);
  final Map<int, User> users = Map();

  MomentListStage createState() => MomentListStage();
}

class MomentListStage extends State<MomentListView>
    with AutomaticKeepAliveClientMixin {
  int pageNo = 1;
  int pageSize = 10;
  int times = 0;

  _getList() async {
    var response = await getMomentList(pageNo, pageSize);
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.
    response!.users.forEach((e) => widget.users[e.id] = e);
    widget.list.addAll(response.list);
    if (!mounted) return; //这是什么大坑啊，切换组件后，这个组件销毁，切换回来重建但是没有mounted
    //不能setState，但是可以更新Widget里的
    //https://github.com/flutter/flutter/issues/27680
    setState(() => {times += 1});
  }
  late final ScrollController _controller = ScrollController()
  ..addListener(() {
  if (this._controller.position.atEdge) {
  _getList();
  }
  });
  initState() {
    print("add initState");
    super.initState();
    _getList();
  }

  @override
  bool get wantKeepAlive => true;

  @override
  void dispose() {
    print("add dispose");

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    var list = widget.list;
    var users = widget.users;
    return ListView.separated(
        controller: this._controller,
        itemCount:  list.length,
        separatorBuilder: (BuildContext context, int index) {
          return Divider();
        },
        itemBuilder: (context, index) {
          return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(children: [
                  Text('${users[list[index].userId]!.name}  ${list[index].createdAt}'),
                ]),
                MarkdownBody(data: '${list[index].content}'),
                Row(children: [
                  Icon(
                    Icons.star,
                    color: Colors.blue,
                  ),
                  Icon(
                    Icons.favorite,
                    color: Colors.red,
                  ),
                ])
              ]);
        });
  }
}
