/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo

import android.app.Activity
import android.app.Application
import android.os.Bundle
import com.immomo.mls.`fun`.lt.SIApplication


/**
 * Author       :   wu.tianlong@immomo.com
 * Date         :   2019/1/10
 * Time         :   下午4:59
 * Description  :
 */
class ActivityLifecycleMonitor : Application.ActivityLifecycleCallbacks {
    private var mCount = 0
    override fun onActivityCreated(activity: Activity, bundle: Bundle?) {}
    override fun onActivityStarted(activity: Activity) {
        mCount++
        if (mCount == 1) {
            SIApplication.setIsForeground(true)
        }
    }

    override fun onActivityResumed(activity: Activity) {}
    override fun onActivityPaused(activity: Activity) {}
    override fun onActivityStopped(activity: Activity) {
        mCount--
        if (mCount == 0) {
            SIApplication.setIsForeground(false)
        }
    }

    override fun onActivitySaveInstanceState(activity: Activity, bundle: Bundle) {}
    override fun onActivityDestroyed(activity: Activity) {}
}