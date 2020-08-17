package xyz.hoper.dart

import android.app.Activity
import android.os.Bundle
import android.view.View
import android.widget.TextView
import androidx.annotation.Nullable
import androidx.appcompat.app.AppCompatActivity
import java.lang.ref.WeakReference


class NativeActivity: AppCompatActivity(), View.OnClickListener {

    private lateinit var mOpenNative: TextView
    private lateinit var mOpenFlutter: TextView
    override fun onCreate(@Nullable savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        sRef = WeakReference(this)
        setContentView(R.layout.native_page)
        mOpenNative = findViewById(R.id.open_native)
        mOpenNative.setOnClickListener(this)
        mOpenFlutter = findViewById(R.id.open_flutter)
        mOpenFlutter.setOnClickListener(this)
    }

    override fun onDestroy() {
        super.onDestroy()
        sRef?.clear()
        sRef = null
    }
    override fun onClick(v: View?) {
        val params= HashMap<Any?, Any?>()
        params["test1"] = "v_test1"
        params["test2"] = "v_test2"
        //Add some params if needed.
       if (v === mOpenFlutter) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_PAGE_URL, params)
        }
    }


    companion object {
        var sRef: WeakReference<NativeActivity?>? = null
    }
}