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
import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.webkit.WebChromeClient
import android.webkit.WebView
import android.widget.Toast
import xyz.hoper.dart.R


/**
 * Created by fanqiang on 2018/12/25.
 */
@SuppressLint("Registered")
class WebActivity : Activity() {
    private lateinit var webView: WebView
    private var url: String? = null
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        url = intent.getStringExtra(INTENT_URL)
        setContentView(R.layout.activity_web)
        webView = findViewById(R.id.wv)
        webView.webChromeClient = client
        webView.loadUrl(url!!)
    }

    private val client: WebChromeClient = object : WebChromeClient() {
        private val TOAST_DURATION_LEAST = 500
        private var tip = false
        private var toast: Toast? = null
        private var toastBegin: Long = 0
        override fun onProgressChanged(view: WebView, newProgress: Int) {
            if (newProgress != 100 && !tip) {
                toast = Toast.makeText(view.context, view.context.getText(R.string.wait), Toast.LENGTH_SHORT)
                toast!!.show()
                toastBegin = System.currentTimeMillis()
                tip = true
            } else if (newProgress == 100 && tip) {
                val current = System.currentTimeMillis()
                if (current - toastBegin >= TOAST_DURATION_LEAST) toast!!.cancel() else {
                    view.postDelayed({ toast!!.cancel() }, TOAST_DURATION_LEAST - (current - toastBegin))
                }
            }
        }
    }

    override fun onDestroy() {
        webView.webChromeClient = null
        super.onDestroy()
    }

    companion object {
        const val INTENT_URL = "INTENT_URL"
        @JvmStatic
        fun startActivity(context: Context, url: String?) {
            val intent = Intent(context, WebActivity::class.java)
            intent.putExtra(INTENT_URL, url)
            context.startActivity(intent)
        }
    }
}