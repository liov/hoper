/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo

import android.content.Context
import android.content.Intent
import com.google.zxing.client.android.camera.CameraConfigurationUtils
import com.immomo.mls.adapter.MLSQrCaptureAdapter

class MLSQrCaptureImpl : MLSQrCaptureAdapter {
    override fun startQrCapture(context: Context) {
        val intent = Intent(context, CameraConfigurationUtils::class.java)
        context.startActivity(intent)
    }
}