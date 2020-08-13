package xyz.hoper.dart

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.util.Log


object PageRouter {
    const val FLUTTER_PAGE_URL = "flutterPage"
    const val Native_PAGE_URL = "nativePage"

    @JvmOverloads
    fun openPageByUrl(context: Context, url: String, params: Map<*, *>?, requestCode: Int = 0): Boolean {
        val path = url.split("\\?".toRegex()).toTypedArray()[0]
        Log.i("openPageByUrl", path)
        return try {
           if (url.startsWith(FLUTTER_PAGE_URL)) {
                context.startActivity(Intent(context, MainActivity::class.java))
                return true
            } else
                context.startActivity(Intent(context, NativeActivity::class.java))
                return true
        } catch (t: Throwable) {
            false
        }
    }
}