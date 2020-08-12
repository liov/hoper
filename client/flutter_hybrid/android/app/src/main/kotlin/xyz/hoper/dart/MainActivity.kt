package xyz.hoper.dart

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import android.view.View
import android.widget.TextView
import androidx.annotation.Nullable
import java.lang.ref.WeakReference


class MainActivity : AppCompatActivity(), View.OnClickListener {
    private lateinit var mOpenNative: TextView
    private lateinit var mOpenFlutter: TextView
    private lateinit var mOpenFlutterFragment: TextView
    override fun onCreate(@Nullable savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        sRef = WeakReference<MainActivity?>(this)
        setContentView(R.layout.native_page)
        mOpenNative = findViewById(R.id.open_native)
        mOpenFlutter = findViewById(R.id.open_flutter)
        mOpenFlutterFragment = findViewById(R.id.open_flutter_fragment)
        mOpenNative.setOnClickListener(this)
        mOpenFlutter.setOnClickListener(this)
        mOpenFlutterFragment.setOnClickListener(this)
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
        if (v === mOpenNative) {
            PageRouter.openPageByUrl(this, PageRouter.NATIVE_PAGE_URL, params)
        } else if (v === mOpenFlutter) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_PAGE_URL, params)
        } else if (v === mOpenFlutterFragment) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_FRAGMENT_PAGE_URL, params)
        }
    }

    companion object {
        var sRef: WeakReference<MainActivity?>? = null
    }
}
