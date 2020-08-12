package xyz.hoper.dart

import android.app.Application
import android.os.Build
import android.os.Environment
import android.util.Log
import com.idlefish.flutterboost.FlutterBoost
import com.idlefish.flutterboost.Utils
import com.idlefish.flutterboost.interfaces.INativeRouter
import com.immomo.mls.MLSEngine
import com.immomo.mls.`fun`.lt.SIApplication
import com.immomo.mls.global.LVConfigBuilder
import io.flutter.embedding.android.FlutterView
import io.flutter.plugin.common.MethodChannel
import io.flutter.plugin.common.StandardMessageCodec
import org.luaj.vm2.Globals
import xyz.hoper.dart.momo.ActivityLifecycleMonitor
import xyz.hoper.dart.momo.GlobalStateListener
import xyz.hoper.dart.momo.MLSQrCaptureImpl
import xyz.hoper.dart.momo.QRResultHandler
import xyz.hoper.dart.momo.provider.GlideImageProvider
import xyz.hoper.dart.momo.zing.OuterResultHandler


class App : Application() {
    override fun onCreate() {
        super.onCreate()
        val router = INativeRouter { context, url, urlParams, _, _ ->
            val assembleUrl: String = Utils.assembleUrl(url, urlParams)
            PageRouter.openPageByUrl(context!!, assembleUrl, urlParams)
        }
        val boostLifecycleListener = object : FlutterBoost.BoostLifecycleListener {
            override fun beforeCreateEngine() {}
            override fun onEngineCreated() {

                // 注册MethodChannel，监听flutter侧的getPlatformVersion调用
                val methodChannel = MethodChannel(FlutterBoost.instance().engineProvider().dartExecutor, "flutter_native_channel")
                methodChannel.setMethodCallHandler { call, result ->
                    if (call.method == "getPlatformVersion") {
                        result.success(Build.VERSION.RELEASE)
                    } else {
                        result.notImplemented()
                    }
                }

                // 注册PlatformView viewTypeId要和flutter中的viewType对应
                FlutterBoost
                        .instance()
                        .engineProvider()
                        .platformViewsController
                        .registry
                        .registerViewFactory("plugins.test/view", TextPlatformViewFactory(StandardMessageCodec.INSTANCE))
            }

            override fun onPluginsRegistered() {}
            override fun onEngineDestroy() {}
        }

        //
        // AndroidManifest.xml 中必须要添加 flutterEmbedding 版本设置
        //
        //   <meta-data android:name="flutterEmbedding"
        //               android:value="2">
        //    </meta-data>
        // GeneratedPluginRegistrant 会自动生成 新的插件方式　
        //
        // 插件注册方式请使用
        // FlutterBoost.instance().engineProvider().getPlugins().add(new FlutterPlugin());
        // GeneratedPluginRegistrant.registerWith()，是在engine 创建后马上执行，放射形式调用
        //
        val platform = FlutterBoost.ConfigBuilder(this, router)
                .isDebug(true)
                .whenEngineStart(FlutterBoost.ConfigBuilder.ANY_ACTIVITY_CREATED)
                .renderMode(FlutterView.RenderMode.texture)
                .lifecycleListener(boostLifecycleListener)
                .build()
        FlutterBoost.instance().init(platform)

        app = this
        init()

        /// -----------配合 Application 使用------------

        /// -----------配合 Application 使用------------
        SIApplication.isColdBoot = true
        registerActivityLifecycleCallbacks(ActivityLifecycleMonitor())
        /// ---------------------END-------------------

        /// ---------------------END-------------------
        MLSEngine.init(this, true) //BuildConfig.DEBUG)
                .setLVConfig(LVConfigBuilder(this)
                        .setRootDir(SD_CARD_PATH)
                        .setCacheDir(SD_CARD_PATH + "cache")
                        .setImageDir(SD_CARD_PATH + "image")
                        .setGlobalResourceDir(SD_CARD_PATH + "g_res")
                        .build())
                .setImageProvider(GlideImageProvider()) //设置图片加载器，若不设置，则不能显示图片
                .setGlobalStateListener(GlobalStateListener()) //设置全局脚本加载监听，可不设置
                .setQrCaptureAdapter(MLSQrCaptureImpl()) //设置二维码工具，可不设置
                .setDefaultLazyLoadImage(false) ///注册静态Bridge
                .registerSC() ///注册Userdata
                .registerUD() ///注册单例
                .registerSingleInsance()
                .build(true)
        /// 设置二维码扫描结果处理工具
        OuterResultHandler.registerResultHandler(QRResultHandler())
        log("onCreate: " + Globals.isInit() + " " + Globals.isIs32bit())
    }


    private fun init() {
        try {
            SD_CARD_PATH = Environment.getExternalStorageDirectory().absolutePath
            if (!SD_CARD_PATH.endsWith("/")) {
                SD_CARD_PATH += "/"
            }
            SD_CARD_PATH += "MLN_Android/"
        } catch (e: Exception) {
        }
    }

    private fun log(s: String) {
        Log.d("app", s)
    }


    companion object {
        lateinit var app: App
        lateinit var SD_CARD_PATH: String

        @JvmStatic
        fun getPackageNameImpl(): String? {
            var sPackageName: String = app.packageName
            if (sPackageName.contains(":")) {
                sPackageName = sPackageName.substring(0, sPackageName.lastIndexOf(":"))
            }
            return sPackageName
        }
    }
}