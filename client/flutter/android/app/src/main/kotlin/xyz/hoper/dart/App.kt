package xyz.hoper.dart


import android.content.Context
import android.content.Intent
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.Build
import android.os.Environment
import android.util.Log
import com.journeyapps.barcodescanner.CaptureActivity
import hoper.xyz.dart.bridge.LuaEnum
import hoper.xyz.dart.bridge.SILuaBridge
import hoper.xyz.dart.bridge.StaticBridge
import io.flutter.app.FlutterApplication
import kotlinx.coroutines.*

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
        const val TAG: String = "App"

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