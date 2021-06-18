import 'package:app/page/moment/momentListView.dart';
import 'package:app/page/webview/webview.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

class FourthPage extends StatefulWidget {
  FourthPage({Key? key, this.title}) : super(key: key);

  final String? title;
  final List<String> _tabValues = [
    '关注',
    '推荐',
    '刚刚',
  ];

  @override
  FourthPageState createState() => FourthPageState();
}

class FourthPageState extends State<FourthPage>  with SingleTickerProviderStateMixin{

  late TabController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TabController(
      length: widget._tabValues.length,
      initialIndex: 1,
      vsync: this,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: TabBar(
          isScrollable: true,
          tabs: widget._tabValues.map((choice) {
            return Tab(
              text: choice,
            );
          }).toList(),
          controller: _controller,
        ),
      ),
      body: TabBarView(
        controller: _controller,
        children: widget._tabValues.map((f) {
          if (f == "推荐") return MomentListView();
          return Center(
            child: Text(f),
          );
        }).toList(),
      ),
      floatingActionButton: FloatingActionButton(
        heroTag: 'login',
        onPressed: ()=>Get.to(WebViewExample()),

        tooltip: 'ToBrowser',
        child: Icon(Icons.send),
      ),
    );
  }
}
