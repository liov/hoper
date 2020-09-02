package xyz.hoper.dart


import android.content.Intent
import android.os.Build
import android.os.Environment
import android.util.Log
import com.common.luakit.LuaHelper
import com.immomo.mls.MLSBuilder
import com.immomo.mls.MLSEngine
import com.immomo.mls.global.LVConfigBuilder
import com.immomo.mls.wrapper.Register
import com.immomo.mmui.MMUIActivity
import com.immomo.mmui.MMUIEngine
import com.immomo.mmui.MMUIEngine.LinkHolder
import com.journeyapps.barcodescanner.CaptureActivity
import hoper.xyz.dart.bridge.LuaEnum
import hoper.xyz.dart.bridge.SILuaBridge
import hoper.xyz.dart.bridge.StaticBridge
import io.flutter.app.FlutterApplication
import org.luaj.vm2.Globals
import xyz.hoper.dart.momo.GlobalStateListener
import xyz.hoper.dart.momo.provider.GlideImageProvider
import xyz.hoper.dart.momo.zing.OuterResultHandler
import xyz.hoper.dart.momo.zing.QRResultHandler
import kotlinx.coroutines.*

class App : FlutterApplication() {

    var alreadyRegistered = false

    override fun onCreate() {
        super.onCreate()
        FlutterEngineFactory.createFlutterEngine(this)
        instance = this
        //GlobalScope.launch { luaOpen() }// 在后台启动一个新的协程并继续
    }

    private fun luaOpen() {
        LuaHelper.startLuaKit(this)
        MLSEngine.init(this, BuildConfig.DEBUG)
                .setLVConfig(LVConfigBuilder(this)
                        .setSdcardDir(SD_CARD_PATH) //设置sdcard目录
                        .setRootDir(SD_CARD_PATH + "root") //设置lua根目录
                        .setImageDir(SD_CARD_PATH + "image") //设置lua图片根目录
                        .setCacheDir(SD_CARD_PATH + "cache") //设置lua缓存目录
                        .setGlobalResourceDir(SD_CARD_PATH + "g_res") //设置资源文件目录
                        .build())
                .setImageProvider(GlideImageProvider()) //lua加载图片工具，不实现的话，图片无法展示
                .setGlobalStateListener(GlobalStateListener()) //设置全局脚本加载监听，可不设置
                .setQrCaptureAdapter { context ->
                    val intent = Intent(context, CaptureActivity::class.java)
                    context.startActivity(intent)
                } //设置二维码工具，可不设置
                .setDefaultLazyLoadImage(false)
                .registerSC(Register.newSHolderWithLuaClass(StaticBridge.LUA_CLASS_NAME, StaticBridge::class.java)) //注册静态Bridge
                .registerUD() //注册Userdata
                .registerSingleInsance(MLSBuilder.SIHolder(SILuaBridge.LUA_CLASS_NAME, SILuaBridge::class.java))//注册单例
                .registerConstants(LuaEnum::class.java) // enum in lua
                .build(true)

        MMUIEngine.init(applicationContext)
        MMUIEngine.preInit(1)
        /// 设置二维码扫描结果处理工具
        OuterResultHandler.registerResultHandler(QRResultHandler())
        Log.d(TAG, "onCreate: " + Globals.isInit() + " " + Globals.isIs32bit())
    }

    //释放flutter引擎
    override fun onTerminate() {
        FlutterEngineFactory.destroyEngine()
        super.onTerminate()
    }

    init {
        try {
            //安卓Q API变化
            SD_CARD_PATH = getExternalFilesDir(null)!!.absolutePath

            if (!SD_CARD_PATH.endsWith("/")) {
                SD_CARD_PATH += "/"
            }
            SD_CARD_PATH += "hoper_lua/"
        } catch (e: Exception) {
        }
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