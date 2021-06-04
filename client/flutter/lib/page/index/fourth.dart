import 'package:app/page/moment/momentListView.dart';
import 'package:flutter/material.dart';

class FourthPage extends StatefulWidget {
  FourthPage({Key? key, String? title}) : super(key: key);


  @override
  FourthPageState createState() => FourthPageState();
}

class FourthPageState extends State<FourthPage> {
  final List<String> _tabValues = [
    '关注',
    '推荐',
    '刚刚',
  ];

  late TabController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TabController(
      length: _tabValues.length,
      initialIndex: 1,
      vsync: ScrollableState(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: TabBar(
          isScrollable: true,
          tabs: _tabValues.map((choice) {
            return Tab(
              text: choice,
            );
          }).toList(),
          controller: _controller,
        ),
      ),
      body: TabBarView(
        controller: _controller,
        children: _tabValues.map((f) {
          if (f == "推荐") return MomentListView();
          return Center(
            child: Text(f),
          );
        }).toList(),
      ),
    );
  }
}
