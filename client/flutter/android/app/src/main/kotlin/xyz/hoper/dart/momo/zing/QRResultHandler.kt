/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.zing

import android.app.Activity
import android.content.Intent
import android.graphics.Bitmap
import android.net.Uri
import android.text.TextUtils
import android.widget.Toast
import com.google.zxing.Result
import com.immomo.mls.HotReloadHelper
import com.immomo.mls.MLSAdapterContainer
import com.immomo.mls.MLSBundleUtils
import com.immomo.mls.MLSEngine
import com.immomo.mls.activity.LuaViewActivity
import com.immomo.mls.util.FileUtil
import com.immomo.mls.utils.MainThreadExecutor
import com.immomo.mls.utils.ScriptLoadException

/**
 * Created by Xiong.Fangyu on 2019/4/22
 */
class QRResultHandler : OuterResultHandler.IResultHandler {
    override fun handle(activity: Activity, rawResult: Result, barcode: Bitmap): Boolean {
        val code: String = rawResult.text
        if (TextUtils.isEmpty(code)) return false
        if (HotReloadHelper.isIPPortString(code)) {
            val r = HotReloadHelper.setUseWifi(code)
            if (!r) toast("connect with wifi failed") else toast("connecting...")
            activity.finish()
            return true
        }
        val uri = Uri.parse(code)
        val intent = Intent(activity, LuaViewActivity::class.java)
        val initData = MLSBundleUtils.createInitData(code).forceNotUseX64()
        if (isDebugScript(uri)) {
            handleDebugScript(code)
            activity.finish()
            return true
        }
        //            initData.doAutoPreload = !uri.getHost().startsWith("172.16") || uri.getPath().endsWith(".zip");
        intent.putExtras(MLSBundleUtils.createBundle(initData))
        activity.startActivity(intent)
        return true
    }

    companion object {
        private const val DEBUG_SCRIPT = "debug.lua"
        private fun handleDebugScript(code: String) {
            Thread(Runnable {
                try {
                    MLSAdapterContainer.getHttpAdapter().downloadLuaFileSync(code, FileUtil.getLuaDir().absolutePath, DEBUG_SCRIPT, null, null, null, 0)
                    toast("下载debug脚本成功")
                } catch (e: ScriptLoadException) {
                    e.printStackTrace()
                    toast("下载debug脚本失败：")
                }
            }).start()
        }

        private fun toast(msg: String) {
            MainThreadExecutor.post { Toast.makeText(MLSEngine.getContext(), msg, Toast.LENGTH_LONG).show() }
        }

        private fun isDebugScript(uri: Uri): Boolean {
            var path = uri.path
            val index = path!!.lastIndexOf('/')
            if (index >= 0) {
                path = path.substring(index + 1)
            }
            return DEBUG_SCRIPT == path
        }
    }
}