package xyz.hoper.dart


import androidx.annotation.NonNull
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.android.SplashScreen
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel

class MainActivity: FlutterActivity() {

    companion object{
        const val CHANNEL = "xyz.hoper.native/view"
        const val Tag = "MainActivity"
    }

    override fun configureFlutterEngine(@NonNull flutterEngine: FlutterEngine) {
        if(!(this.applicationContext as App).alreadyRegistered) {
            super.configureFlutterEngine(flutterEngine)
            (this.applicationContext as App).alreadyRegistered = true
        }

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, CHANNEL).setMethodCallHandler { call, result ->
            if (call.method == "toNative") {
               val success=  PageRouter.openPageByUrl(this, PageRouter.Native_PAGE_URL, call.arguments())

                if (success) {
                    result.success("跳转成功")
                } else {
                    result.error("UNAVAILABLE", "跳转失败", null)
                }
            } else {
                result.notImplemented()
            }
        }
    }

    override fun provideSplashScreen(): SplashScreen {
        return SplashScreenWithTransition()
    }
}
