package xyz.hoper.dart

import android.content.Context
import android.view.View
import android.widget.TextView
import io.flutter.plugin.common.MessageCodec
import io.flutter.plugin.platform.PlatformView
import io.flutter.plugin.platform.PlatformViewFactory

class TextPlatformViewFactory(createArgsCodec: MessageCodec<Any?>?) : PlatformViewFactory(createArgsCodec) {
    override fun create(context: Context, i: Int, o: Any): PlatformView {
        return TextPlatformView(context)
    }

    private class TextPlatformView(context: Context?) : PlatformView {
        private val platformTv: TextView
        override fun getView(): View {
            return platformTv
        }

        override fun dispose() {}

        init {
            platformTv = TextView(context)
            platformTv.text = "PlatformView Demo"
        }
    }
}