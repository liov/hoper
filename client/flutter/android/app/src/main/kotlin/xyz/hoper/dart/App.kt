package xyz.hoper.dart


import android.content.Intent
import android.os.Environment
import android.util.Log
import com.common.luakit.LuaHelper
import com.immomo.mls.MLSBuilder
import com.immomo.mls.MLSEngine
import com.immomo.mls.global.LVConfigBuilder
import com.immomo.mls.wrapper.Register
import com.journeyapps.barcodescanner.CaptureActivity
import hoper.xyz.dart.bridge.LuaEnum
import hoper.xyz.dart.bridge.SILuaBridge
import hoper.xyz.dart.bridge.StaticBridge
import io.flutter.app.FlutterApplication
import org.luaj.vm2.Globals
import xyz.hoper.dart.momo.GlobalStateListener
import xyz.hoper.dart.momo.zing.QRResultHandler
import xyz.hoper.dart.momo.provider.GlideImageProvider
import xyz.hoper.dart.momo.zing.OuterResultHandler


class App : FlutterApplication() {

    var alreadyRegistered = false

    override fun onCreate() {
        super.onCreate()
        LuaHelper.startLuaKit(this)
        FlutterEngineFactory.createFlutterEngine(this)
        instance = this
        init()
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
                .setQrCaptureAdapter {context->
                    val intent = Intent(context, CaptureActivity::class.java)
                    context.startActivity(intent)
                } //设置二维码工具，可不设置
                .setDefaultLazyLoadImage(false)
                .registerSC(Register.newSHolderWithLuaClass(StaticBridge.LUA_CLASS_NAME, StaticBridge::class.java)) //注册静态Bridge
                .registerUD() //注册Userdata
                .registerSingleInsance(MLSBuilder.SIHolder(SILuaBridge.LUA_CLASS_NAME, SILuaBridge::class.java))//注册单例
                .registerConstants(LuaEnum::class.java) // enum in lua
                .build(true)

        /// 设置二维码扫描结果处理工具
        OuterResultHandler.registerResultHandler(QRResultHandler())
        Log.d(TAG, "onCreate: " + Globals.isInit() + " " + Globals.isIs32bit())
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

        private fun init() {
            try {
                SD_CARD_PATH = Environment.getExternalStorageDirectory().absolutePath

                if (!SD_CARD_PATH.endsWith("/")) {
                    SD_CARD_PATH += "/"
                }
                SD_CARD_PATH += "hoper_lua/"
            } catch (e: Exception) {
            }
        }
    }
}