package xyz.hoper.dart

import android.os.Bundle
import android.view.View
import android.widget.TextView
import androidx.appcompat.app.AppCompatActivity
import java.util.*

class NativePageActivity : AppCompatActivity(), View.OnClickListener {
    private var mOpenNative: TextView? = null
    private var mOpenFlutter: TextView? = null
    private var mOpenFlutterFragment: TextView? = null
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.native_page)
        mOpenNative = findViewById(R.id.open_native)
        mOpenFlutter = findViewById(R.id.open_flutter)
        mOpenFlutterFragment = findViewById(R.id.open_flutter_fragment)
        mOpenNative!!.setOnClickListener(this)
        mOpenFlutter!!.setOnClickListener(this)
        mOpenFlutterFragment!!.setOnClickListener(this)
    }

    override fun onClick(v: View) {
        val params = HashMap<Any?, Any?>()
        params["test1"] = "v_test1"
        params["test2"] = "v_test2"
        if (v === mOpenNative) {
            PageRouter.openPageByUrl(this, PageRouter.NATIVE_PAGE_URL, params)
        } else if (v === mOpenFlutter) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_PAGE_URL, params)
        } else if (v === mOpenFlutterFragment) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_FRAGMENT_PAGE_URL, params)
        }
    }
}