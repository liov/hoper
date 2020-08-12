/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.activity

import android.annotation.SuppressLint
import android.app.Activity
import android.view.WindowManager
import android.widget.Toast
import com.immomo.mls.HotReloadHelper
import com.immomo.mls.HotReloadHelper.ConnectListener
import com.immomo.mls.MLSEngine
import com.immomo.mls.activity.LuaViewActivity
import com.immomo.mls.utils.MainThreadExecutor

@SuppressLint("Registered")
open class BaseActivity : Activity(), ConnectListener {
    override fun onConnected(hasCallback: Boolean) {
        MainThreadExecutor.post {
            window.addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON)
            if (!hasCallback) {
                Toast.makeText(MLSEngine.getContext(), "connect with wifi success", Toast.LENGTH_LONG).show()
                startTeach(false)
            }
            HotReloadHelper.setConnectListener(null)
        }
    }

    protected fun startTeach(usb: Boolean) {
        LuaViewActivity.startHotReload(this, usb)
    }

    override fun onDisConnected() {
        MainThreadExecutor.post {
            window.clearFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON)
            HotReloadHelper.setConnectListener(null)
        }
    }
}