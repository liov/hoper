package xyz.hoper.dart

import io.flutter.embedding.android.FlutterActivity
import com.common.luakit.LuaHelper
import io.flutter.plugins.GeneratedPluginRegistrant
import android.os.Bundle
import io.flutter.embedding.engine.FlutterEngine

class MainActivity: FlutterActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        LuaHelper.startLuaKit(this)
        GeneratedPluginRegistrant.registerWith(FlutterEngine(this))
    }
}
