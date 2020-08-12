package xyz.hoper.dart

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.util.Log
import com.idlefish.flutterboost.containers.BoostFlutterActivity
import java.util.*

object PageRouter {
    private val pageName = object : HashMap<String, String>() {
        init {
            put("first", "first")
            put("second", "second")
            put("tab", "tab")
            put("sample://flutterPage", "flutterPage")
        }
    }
    const val NATIVE_PAGE_URL = "sample://nativePage"
    const val FLUTTER_PAGE_URL = "sample://flutterPage"
    const val FLUTTER_FRAGMENT_PAGE_URL = "sample://flutterFragmentPage"

    @JvmOverloads
    fun openPageByUrl(context: Context, url: String, params: Map<*, *>?, requestCode: Int = 0): Boolean {
        val path = url.split("\\?".toRegex()).toTypedArray()[0]
        Log.i("openPageByUrl", path)
        return try {
            if (pageName.containsKey(path)) {
                val intent = BoostFlutterActivity.withNewEngine().url(pageName[path]!!).params(params as MutableMap<String, Any>)
                        .backgroundMode(BoostFlutterActivity.BackgroundMode.opaque).build(context)
                if (context is Activity) {
                    context.startActivityForResult(intent, requestCode)
                } else {
                    context.startActivity(intent)
                }
                return true
            } else if (url.startsWith(FLUTTER_FRAGMENT_PAGE_URL)) {
                context.startActivity(Intent(context, FlutterFragmentPageActivity::class.java))
                return true
            } else if (url.startsWith(NATIVE_PAGE_URL)) {
                context.startActivity(Intent(context, NativePageActivity::class.java))
                return true
            }
            false
        } catch (t: Throwable) {
            false
        }
    }
}