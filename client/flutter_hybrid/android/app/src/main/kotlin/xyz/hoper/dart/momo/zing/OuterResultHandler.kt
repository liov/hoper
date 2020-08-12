package xyz.hoper.dart.momo.zing

import  android.app.Activity
import android.graphics.Bitmap
import com.google.zxing.Result
import java.util.*

/**
 * Created by XiongFangyu on 2018/4/11.
 */
object OuterResultHandler {
    private val resultHandlers: MutableList<IResultHandler> = ArrayList(2)

    /**
     * @hide
     */
    fun handleDecodeInternally(activity: Activity, rawResult: Result, barcode: Bitmap): Boolean {
        val l = resultHandlers.size
        for (i in 0 until l) {
            val handler = resultHandlers[i]
            if (handler.handle(activity, rawResult, barcode)) {
                return true
            }
        }
        return false
    }

    fun registerResultHandler(handler: IResultHandler) {
        if (!resultHandlers.contains(handler)) {
            resultHandlers.add(0, handler)
        }
    }

    fun unregisterResultHandler(handler: IResultHandler?) {
        resultHandlers.remove(handler)
    }

    interface IResultHandler {
        fun handle(activity: Activity, rawResult: Result, barcode: Bitmap): Boolean
    }
}