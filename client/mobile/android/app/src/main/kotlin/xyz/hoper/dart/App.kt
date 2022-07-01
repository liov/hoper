package xyz.hoper.dart

import io.flutter.app.FlutterApplication
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.embedding.engine.FlutterEngineCache

class App : FlutterApplication() {

    var alreadyRegistered = false
    var engineCached = false

    override fun onCreate() {
        super.onCreate()
        instance = this
        //GlobalScope.launch { luaOpen() }// 在后台启动一个新的协程并继续
    }


    //释放flutter引擎
    override fun onTerminate() {
        super.onTerminate()
    }

    companion object {
        private lateinit var instance: App
        const val ENGINE_ID = "cached_engine_id"
        val flutterEngineCache by lazy { FlutterEngineCache.getInstance() }
        fun getInstance(): App {
            return instance
        }
        fun getEngine(): FlutterEngine {
            return flutterEngineCache[ENGINE_ID]!!
        }

        lateinit var SD_CARD_PATH: String
        const val TAG: String = "hoper"

        @JvmStatic
        fun getPackageNameImpl(): String? {
            var sPackageName: String = App::class.java.`package`!!.name
            if (sPackageName.contains(":")) {
                sPackageName = sPackageName.substring(0, sPackageName.lastIndexOf(":"))
            }
            return sPackageName
        }

    }
}