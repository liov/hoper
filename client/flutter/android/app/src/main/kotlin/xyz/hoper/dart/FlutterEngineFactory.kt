package xyz.hoper.dart

import android.content.Context
import android.util.Log
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.embedding.engine.FlutterEngineCache
import io.flutter.embedding.engine.dart.DartExecutor


object FlutterEngineFactory {
    const val ENGINE_ID = "cached_engine_id"
    private const val Tag = "FlutterEngineFactory"
    private val flutterEngineCache by lazy { FlutterEngineCache.getInstance() }

    fun createFlutterEngine(context: Context) {
        if(flutterEngineCache.contains(ENGINE_ID)) return
        val engine = FlutterEngine(context)
        engine.navigationChannel.setInitialRoute("/")
        engine.dartExecutor
                .executeDartEntrypoint(DartExecutor.DartEntrypoint.createDefault())
        flutterEngineCache.put(ENGINE_ID, engine)
        Log.d(Tag,"执行了吗")
    }

    fun destroyEngine() {
        flutterEngineCache[ENGINE_ID]?.destroy()
    }

    fun getEngine(): FlutterEngine {
        return flutterEngineCache[ENGINE_ID]!!
    }
}