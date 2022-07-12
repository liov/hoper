package xyz.hoper.dart

import android.os.Build
import android.os.Bundle
import androidx.core.view.WindowCompat
import androidx.annotation.NonNull
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.embedding.engine.FlutterEngineCache
import io.flutter.plugin.common.MethodChannel

class MainActivity: FlutterActivity() {

    companion object{
        const val Tag = "MainActivity"
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        // Aligns the Flutter view vertically with the window.
        WindowCompat.setDecorFitsSystemWindows(getWindow(), false)

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            // Disable the Android splash screen fade out animation to avoid
            // a flicker before the similar frame is drawn in Flutter.
            splashScreen.setOnExitAnimationListener { splashScreenView -> splashScreenView.remove() }
        }

        super.onCreate(savedInstanceState)
    }

    override fun configureFlutterEngine(@NonNull flutterEngine: FlutterEngine) {
        if(!(this.applicationContext as App).alreadyRegistered) {
            super.configureFlutterEngine(flutterEngine)
            (this.applicationContext as App).alreadyRegistered = true
        }
        FlutterEngineCache.getInstance().put(App.ENGINE_ID, flutterEngine)
        (this.applicationContext as App).engineCached = true

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, NativeActivity.CHANNEL).setMethodCallHandler { call, result ->
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

}
