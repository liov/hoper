package xyz.hoper.dart

import io.flutter.embedding.android.FlutterActivity
import com.common.luakit.LuaHelper;
import android.os.Bundle

class MainActivity: FlutterActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        //LuaHelper.startLuaKit(this)
    }
}
