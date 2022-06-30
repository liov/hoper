package xyz.hoper.dart

import android.app.Activity
import android.os.Bundle
import android.view.View
import xyz.hoper.dart.databinding.NativePageBinding
import java.lang.ref.WeakReference


class NativeActivity: Activity(), View.OnClickListener {


    private lateinit var binding: NativePageBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = NativePageBinding.inflate(layoutInflater)
        sRef = WeakReference(this)
        setContentView(binding.root)
        binding.openFlutter.setOnClickListener(this)
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
       if (v === binding.openFlutter) {
            PageRouter.openPageByUrl(this, PageRouter.FLUTTER_PAGE_URL, params)
        }
    }


    companion object {
        const val CHANNEL = "xyz.hoper.native/view"
        const val Tag = "NativeActivity"

        var sRef: WeakReference<NativeActivity?>? = null
    }
}