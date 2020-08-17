package xyz.hoper.dart


import android.Manifest
import android.app.Activity
import android.content.pm.PackageManager
import android.os.Bundle
import android.view.WindowManager
import android.widget.FrameLayout
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import com.immomo.mls.*
import com.immomo.mls.activity.LuaViewActivity
import com.immomo.mls.utils.MainThreadExecutor


class LuaActivity : AppCompatActivity(), HotReloadHelper.ConnectListener {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        storageAndCameraPermission()
        val frameLayout = FrameLayout(this)
        setContentView(frameLayout)
        val instance = MLSInstance(this, true, true)
        instance.setContainer(frameLayout)
        val initData = InitData(Constants.ASSETS_PREFIX +"/lua/view/demo.lua") //MLSBundleUtils.parseFromBundle(bundle);MLSBundleUtils.createBundle(url)
        instance.setData(initData)
        if (!instance.isValid) {
            //非法url
            Toast.makeText(this, "something wrong", Toast.LENGTH_SHORT).show()
        }

    }

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

    private fun startTeach(usb: Boolean) {
        LuaViewActivity.startHotReload(this, usb)
    }

    override fun onDisConnected() {
        MainThreadExecutor.post {
            window.clearFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON)
            HotReloadHelper.setConnectListener(null)
        }
    }

    companion object {
        private const val REQUEST_EXTERNAL_STORAGE = 1
        private val PERMISSIONS_STORAGE = arrayOf(
                "android.permission.READ_EXTERNAL_STORAGE",
                "android.permission.WRITE_EXTERNAL_STORAGE",
                Manifest.permission.CAMERA)
    }
    fun storageAndCameraPermission() {
        try {
            //检测是否有写的权限
            val permission: Int = ActivityCompat.checkSelfPermission(this,
                    "android.permission.WRITE_EXTERNAL_STORAGE")
            if (permission != PackageManager.PERMISSION_GRANTED) {
                // 没有写的权限，去申请写的权限，会弹出对话框
                ActivityCompat.requestPermissions(this, PERMISSIONS_STORAGE, REQUEST_EXTERNAL_STORAGE)
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

}
