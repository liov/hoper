package xyz.hoper.dart

import io.flutter.embedding.android.FlutterActivity
import com.common.luakit.LuaHelper;
import android.os.Bundle
import androidx.annotation.NonNull
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

class MainActivity: FlutterActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        //LuaHelper.startLuaKit(this)
    }
    private val CHANNEL = "xyz.hoper.native/view"

    override fun configureFlutterEngine(@NonNull flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            if (call.method == "toNative") {
               val success=  PageRouter.openPageByUrl(this, PageRouter.Native_PAGE_URL, null)

                if (success) {
                    result.success("跳转成功")
                } else {
                    result.error("UNAVAILABLE", "Battery level not available.", null)
                }
            } else {
                result.notImplemented()
            }
        }
    }
}
