package xyz.hoper.dart

import io.flutter.app.FlutterApplication

class App : FlutterApplication() {

    var alreadyRegistered = false

    override fun onCreate() {
        super.onCreate()
        FlutterEngineFactory.createFlutterEngine(this)
        instance = this
        //GlobalScope.launch { luaOpen() }// 在后台启动一个新的协程并继续
    }


    //释放flutter引擎
    override fun onTerminate() {
        FlutterEngineFactory.destroyEngine()
        super.onTerminate()
    }

    companion object {
        private lateinit var instance: App
        fun getInstance(): App {
            return instance
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