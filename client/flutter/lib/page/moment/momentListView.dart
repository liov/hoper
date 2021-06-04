import 'package:app/model/moment.dart';
import 'package:app/service/moment.dart';
import 'package:flutter/material.dart';

import 'package:flutter_markdown/flutter_markdown.dart';

class MomentListView extends StatefulWidget {
  List<Moment>? list;
  ScrollController? _controller;
  MomentListStage createState() => MomentListStage();
}

class MomentListStage extends State<MomentListView> with AutomaticKeepAliveClientMixin {
  int pageNo = 0;
  int pageSize = 10;
  int times = 0;

  _getList() async {
    var response = await getMomentList(pageNo, pageSize);
    // If the widget was removed from the tree while the message was in flight,
    // we want to discard the reply rather than calling setState to update our
    // non-existent appearance.


    if (widget.list == null)
      widget.list = response.data;
    else
      widget.list!.addAll(response.data);
    if (!mounted) return;//这是什么大坑啊，切换组件后，这个组件销毁，切换回来重建但是没有mounted
    //不能setState，但是可以更新Widget里的
    //https://github.com/flutter/flutter/issues/27680
    setState(() => {times+=1});
  }

  initState() {
    print("add initState");
    super.initState();
    if (widget._controller == null) {
      widget._controller = ScrollController()
        ..addListener(() {
          if(widget._controller!.position.atEdge){
            _getList();
          }
        });
    }
    if (widget.list == null) _getList();
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
    return ListView.separated(
        controller:widget._controller,
        itemCount: list != null ? list.length : 0,
        separatorBuilder: (BuildContext context, int index) {
          return Divider();
        },
        itemBuilder: (context, index) {
          return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(children: [
                  Text('${list![index].user.name}  ${list![index].createdAt}'),
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
