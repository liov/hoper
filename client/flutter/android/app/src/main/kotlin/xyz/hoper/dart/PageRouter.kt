package xyz.hoper.dart

import android.content.Context
import android.content.Intent
import android.util.Log
import io.flutter.embedding.android.FlutterActivity


object PageRouter {
    const val FLUTTER_PAGE_URL = "flutterPage"
    const val Native_PAGE_URL = "nativePage"

    @JvmOverloads
    @JvmStatic
    fun openPageByUrl(context: Context, url: String, params: Map<*, *>?, requestCode: Int = 0): Boolean {
        val path = url.split("\\?".toRegex()).toTypedArray()[0]
        Log.i("openPageByUrl", path)
        return try {
            when {
                url.startsWith(FLUTTER_PAGE_URL) -> {
                    context.startActivity(FlutterActivity.withCachedEngine(App.ENGINE_ID)
                            .build(context))
                    return true
                }

                else -> context.startActivity(Intent(context, NativeActivity::class.java))
            }
            return true
        } catch (t: Throwable) {
            false
        }
    }
}