import 'package:flutter/cupertino.dart';

import 'moment_list_view.dart';

class MomentListViewWarp extends StatefulWidget {
  MomentListViewWarp({this.tag = "default"}) : super();

  final String tag;
  @override
  _MomentListViewWarpState createState() => _MomentListViewWarpState();
}

class _MomentListViewWarpState extends State<MomentListViewWarp> with AutomaticKeepAliveClientMixin {
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