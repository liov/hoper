import 'package:app/page/ffi/ffi.dart';
import 'package:app/page/index/fourth.dart';
import 'package:app/page/index/index.dart';
import 'package:app/page/lua/lua.dart';
import 'package:app/page/moment/momentListView.dart';

import 'package:app/page/webview/webview.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'model/state/global.dart';
import 'model/state/user.dart';

void main() async {
  runApp(
      MultiProvider(
        providers: [
          ChangeNotifierProvider(create: (context) => UserInfo()),
          Provider(create: (context) => GlobalState()),
        ],
        child:MyApp(),
      ));
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'hoper',
      theme: ThemeData(
        // This is the theme of your application.
        //
        // Try running your application with "flutter run". You'll see the
        // application has a blue toolbar. Then, without quitting the app, try
        // changing the primarySwatch below to Colors.green and then invoke
        // "hot reload" (press "r" in the console where you ran "flutter run",
        // or simply save your changes to "hot reload" in a Flutter IDE).
        // Notice that the counter didn't reset back to zero; the application
        // is not restarted.
        primarySwatch: Colors.blue,
        primaryColor: Colors.blueAccent[700],
        // This makes the visual density adapt to the platform that you run
        // the app on. For desktop platforms, the controls will be smaller and
        // closer together (more dense) than on mobile platforms.
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: MyHomePage(),
      routes: {
        '/moment': (context) => MomentListView(),
      },
    );
  }
}

class MyHomePage extends StatefulWidget {
  MyHomePage({Key key}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {

  int _selectedIndex = 0;
  static const TextStyle optionStyle =
      TextStyle(fontSize: 30, fontWeight: FontWeight.bold);

  void _onItemTapped(int index) {
    setState(() {
      _selectedIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    List<Widget> _widgetOptions = <Widget>[
      IndexPage(title:'🍀'),
      WebViewExample(),
      Text(
        greeting(),
        style: optionStyle,
      ),
      FourthPage(),
    ];
    // This method is rerun every time setState is called, for instance as done
    // by the _incrementCounter method above.
    //
    // The Flutter framework has been optimized to make rerunning build methods
    // fast, so that you can just rebuild anything that needs updating rather
    // than having to individually change instances of widgets.
    return Scaffold(
      body: Center(
        // Center is a layout widget. It takes a single child and positions it
        // in the middle of the parent.
        child: _widgetOptions.elementAt(_selectedIndex),
      ),
      bottomNavigationBar: Bottom(onTap: _onItemTapped),
    );
  }
}

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
      items: <BottomNavigationBarItem>[
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