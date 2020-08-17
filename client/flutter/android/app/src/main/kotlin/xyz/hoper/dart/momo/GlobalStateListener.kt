/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo

import android.util.Log
import com.immomo.mls.utils.DebugLog
import com.immomo.mls.utils.GlobalStateSDKListener
import java.io.PrintStream

/**
 * Created by XiongFangyu on 2018/8/8.
 */
open class GlobalStateListener : GlobalStateSDKListener() {
    override fun newLog(): DebugLog {
        return D()
    }

    protected class D : DebugLog() {
        override fun log(s: String, ps: PrintStream) {
            super.log(s, ps)
            Log.d(TAG, s)
        }
    }

    companion object {
        private const val TAG = "GlobalStateListener"
    }
}