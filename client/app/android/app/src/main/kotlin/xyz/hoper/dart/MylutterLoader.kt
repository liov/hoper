package xyz.hoper.dart

import android.content.Context
import io.flutter.FlutterInjector
import io.flutter.Log
import io.flutter.embedding.engine.loader.FlutterApplicationInfo
import io.flutter.embedding.engine.loader.FlutterLoader
import java.io.File

class MylutterLoader : FlutterLoader() {

    companion object {
        ///完成对FlutterInjector单例的重置，使其属性flutterLoader指向我们的子类
        fun activation() {
            //这里直接使用默认构造函数，与FlutterInjector类中初始化flutterLoader效果是一样的（详见FlutterInjector的fillDefaults()方法）
            val flutterLoader: MylutterLoader = MylutterLoader()
            val flutterInjector = FlutterInjector.Builder().setFlutterLoader(flutterLoader).build()
            //重置FlutterInjector单例
            FlutterInjector.reset()
            FlutterInjector.setInstance(flutterInjector)
            Log.i("------", "已重置FlutterInjector单例")
        }
    }
    //返回准备好的热更新包的路径（本文方案是从服务端下载到zip文件并解压放置到这个路径）
    private fun getHotAppBundlePath(applicationContext: Context): String {
        return applicationContext.filesDir.absolutePath + File.separator + "hot/lib/libapp.so";
    }

    override fun ensureInitializationComplete(
        applicationContext: Context,
        args: Array<out String>?
    ) {
        super.ensureInitializationComplete(applicationContext, args)

        val soFile: File = File(getHotAppBundlePath(applicationContext))
        if (soFile.exists()) {
            try {
                //1.拿到flutterApplicationInfo字段
                val flutterApplicationInfoField = FlutterLoader::class.java.getDeclaredField("flutterApplicationInfo")
                flutterApplicationInfoField.isAccessible = true
                val flutterApplicationInfo = flutterApplicationInfoField[this] as FlutterApplicationInfo
                Log.i(
                    "========",
                    "--aot-shared-library-name=" + flutterApplicationInfo.nativeLibraryDir + flutterApplicationInfo.aotSharedLibraryName
                )

                //2.拿到aotSharedLibraryName修改路径
                val aotSharedLibraryNameField =
                    FlutterApplicationInfo::class.java.getDeclaredField("aotSharedLibraryName")
                aotSharedLibraryNameField.isAccessible = true
                aotSharedLibraryNameField[flutterApplicationInfo] = soFile.absolutePath

                Log.i(
                    "========",
                    "--aot-shared-library-name=" + flutterApplicationInfo.nativeLibraryDir + flutterApplicationInfo.aotSharedLibraryName
                )

                super.ensureInitializationComplete(applicationContext, args)

            } catch (e: Exception) {
                e.printStackTrace()
                e.message?.let { Log.e("----", it) }
            }
        } else {
            Log.i("----", "load fail. 补丁不存在")
        }
    }
}