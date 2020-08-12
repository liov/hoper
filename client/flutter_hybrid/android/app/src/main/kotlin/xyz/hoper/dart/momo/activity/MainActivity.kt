/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.activity

import android.Manifest
import android.annotation.SuppressLint
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Bundle
import android.view.View
import androidx.core.app.ActivityCompat
import com.immomo.mls.*
import com.immomo.mls.activity.LuaViewActivity
import xyz.hoper.dart.momo.activity.WebActivity.Companion.startActivity
import xyz.hoper.dart.R

@SuppressLint("Registered")
class MainActivity : BaseActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        storageAndCameraPermission()
        setContentView(R.layout.activity_main)
        initView()
        findViewById<View>(R.id.tvDevDebug).visibility = if (MLSEngine.DEBUG) View.VISIBLE else View.GONE
    }

    private fun initView() {
        findViewById<View>(R.id.tvInnerDemo).setOnClickListener(this)
        findViewById<View>(R.id.tvDevDebug).setOnClickListener(this)
        findViewById<View>(R.id.tvDemo).setOnClickListener(this)
        findViewById<View>(R.id.tvCourse).setOnClickListener(this)
        findViewById<View>(R.id.tvInstance).setOnClickListener(this)
        findViewById<View>(R.id.tvConsult).setOnClickListener(this)
        findViewById<View>(R.id.tvAbout).setOnClickListener(this)
    }

    override fun onClick(v: View) {
        when (v.id) {
            R.id.tvInnerDemo -> AssetsChooserActivity.startActivity(this, "inner_demo")
            R.id.tvCourse -> startActivity(this, URL_COURSE)
            R.id.tvInstance -> startActivity(this, URL_INSTANCE)
            R.id.tvConsult -> startActivity(this, URL_CONSULT)
            R.id.tvAbout -> startActivity(this, URL_ABOUT)
            R.id.tvDevDebug -> {
                HotReloadHelper.setConnectListener(this)
                startTeach(true)
            }
            R.id.tvDemo -> {
                val intent = Intent(this, LuaViewActivity::class.java)
                val initData = MLSBundleUtils.createInitData(Constants.ASSETS_PREFIX + "gallery/meilishuo.lua")
                intent.putExtras(MLSBundleUtils.createBundle(initData))
                startActivity(intent)
            }
        }
    }

    private fun storageAndCameraPermission() {
        try {
            //检测是否有写的权限
            val permission = ActivityCompat.checkSelfPermission(this,
                    "android.permission.WRITE_EXTERNAL_STORAGE")
            if (permission != PackageManager.PERMISSION_GRANTED) {
                // 没有写的权限，去申请写的权限，会弹出对话框
                ActivityCompat.requestPermissions(this, PERMISSIONS_STORAGE, REQUEST_EXTERNAL_STORAGE)
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

    companion object {
        private const val URL_COURSE = "https://mln.immomo.com/zh-cn/docs/build_dev_environment.html"
        private const val URL_INSTANCE = "https://mln.immomo.com/zh-cn/api/NewListView.lua.html"
        private const val URL_CONSULT = "https://github.com/momotech/MLN"
        private const val URL_ABOUT = "https://mln.immomo.com/zh-cn/"

        // zx add for storage permission 20190806
        private const val REQUEST_EXTERNAL_STORAGE = 1
        private val PERMISSIONS_STORAGE = arrayOf(
                "android.permission.READ_EXTERNAL_STORAGE",
                "android.permission.WRITE_EXTERNAL_STORAGE",
                Manifest.permission.CAMERA)
    }
}