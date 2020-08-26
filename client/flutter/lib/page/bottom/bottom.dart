import 'package:flutter/material.dart';

class Bottom extends StatefulWidget {
  final ValueChanged<int> onTap;

  Bottom({key,@required this.onTap}): super(key: key);
  BottomState createState() => BottomState();
}

class BottomState extends State<Bottom> {

  int _selectedIndex = 0;

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
    widget.onTap(_selectedIndex);
  }
  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
      type: BottomNavigationBarType.fixed,
      backgroundColor: Theme.of(context).primaryColor.withAlpha(127),
      items: const <BottomNavigationBarItem>[
        BottomNavigationBarItem(
          icon: Icon(Icons.movie),
          title: Text('flutter'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.language),
          title: Text('webview'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.search),
          title: Text('rustffi'),
        ),
        BottomNavigationBarItem(
          icon: Icon(Icons.chat_bubble_outline),
          title: Text('lua业务逻辑'),
        ),
      ],
      currentIndex: _selectedIndex,
      selectedItemColor: Theme.of(context).canvasColor,
      onTap: _onItemTapped,
    );
  }
}