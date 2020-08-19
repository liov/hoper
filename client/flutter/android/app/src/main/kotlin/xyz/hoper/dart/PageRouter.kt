package xyz.hoper.dart

import android.content.Context
import android.content.Intent
import android.util.Log
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngineCache


object PageRouter {
    const val FLUTTER_PAGE_URL = "flutterPage"
    const val Native_PAGE_URL = "nativePage"
    const val Lua_PAGE_URL = "luaPage"

    @JvmOverloads
    @JvmStatic
    fun openPageByUrl(context: Context, url: String, params: Map<*, *>?, requestCode: Int = 0): Boolean {
        val path = url.split("\\?".toRegex()).toTypedArray()[0]
        Log.i("openPageByUrl", path)
        return try {
            when {
                url.startsWith(FLUTTER_PAGE_URL) -> {
                    context.startActivity(FlutterActivity.CachedEngineIntentBuilder(MainActivity::class.java, FlutterEngineFactory.ENGINE_ID)
                            .build(context))
                    return true
                }
                url.startsWith(Lua_PAGE_URL) -> {
                    context.startActivity(Intent(context, LuaActivity::class.java))
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