import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_luakit_plugin/flutter_luakit_plugin.dart';



class Lua extends StatefulWidget {
  Lua({Key key}) : super(key: key);

  // This widget is the home page of your application. It is stateful, meaning
  // that it has a State object (defined below) that contains fields that affect
  // how it looks.

  // This class is the configuration for the state. It holds the values (in this
  // case the title) provided by the parent (in this case the App widget) and
  // used by the build method of the State. Fields in a Widget subclass are
  // always marked "final".

  @override
  LuaState createState() => LuaState();
}

class LuaState extends State<Lua> {
  List<dynamic> weathers;
  String _platformVersion = 'Unknown';
  @override
  void initState() {
    super.initState();
    FlutterLuakitPlugin.callLuaFun("WeatherManager", "loadWeather")
        .then((dynamic d) {
      setState(() {
        if (d != null) {
          weathers = d;
        }
      });
    });
    testLua();
  }

  Future<void> testLua() async {
    String platformVersion;
    // Platform messages may fail, so we use a try/catch PlatformException.
    try {
      platformVersion = await FlutterLuakitPlugin.callLuaFun("Lua", "getVersion");
    } on PlatformException {
      platformVersion = 'Failed to get platform version.';
    }

    // If the widget was removed from the tree while the asynchronous platform
    // message was in flight, we want to discard the reply rather than calling
    // setState to update our non-existent appearance.
    if (!mounted) return;

    setState(() {
      _platformVersion = platformVersion;
    });
  }

  @override
  Widget build(BuildContext context) {
    Map<dynamic, dynamic> m =
    weathers != null ? weathers[0] : {"city": "无数据", "sun_info": "无数据"};
     return Text(
        '${m["city"]} ' + "日出日落：" + '${m["sun_info"]} ',
        style: Theme.of(context).textTheme.headline4,
      );
  }
}
