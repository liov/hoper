import 'package:flutter/cupertino.dart';

import 'moment_list_view.dart';

class MomentListViewWrap extends StatefulWidget {
  MomentListViewWrap({this.tag = "default"}) : super();

  final String tag;
  @override
  _MomentListViewWrapState createState() => _MomentListViewWrapState();
}

class _MomentListViewWrapState extends State<MomentListViewWrap> with AutomaticKeepAliveClientMixin {
  @override
  Widget build(BuildContext context) {
    super.build(context);
    return MomentListView(tag:widget.tag);
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  bool get wantKeepAlive => true;
}