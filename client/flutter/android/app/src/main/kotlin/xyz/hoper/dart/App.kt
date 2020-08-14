package xyz.hoper.dart


import android.util.Log
import io.flutter.app.FlutterApplication
import com.common.luakit.LuaHelper

class App : FlutterApplication(){

    var alreadyRegistered = false

    override fun onCreate() {
        super.onCreate()
        LuaHelper.startLuaKit(this)
        FlutterEngineFactory.createFlutterEngine(this)
    }

    //释放flutter引擎
    override fun onTerminate() {
        FlutterEngineFactory.destroyEngine()
        super.onTerminate()
    }

}