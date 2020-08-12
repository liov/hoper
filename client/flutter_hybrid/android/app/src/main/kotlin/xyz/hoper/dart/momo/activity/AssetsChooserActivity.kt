/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.activity

import android.R
import android.annotation.SuppressLint
import android.app.Activity
import android.app.ListActivity
import android.content.Intent
import android.os.Bundle
import android.text.TextUtils
import android.widget.AdapterView
import android.widget.ArrayAdapter
import com.immomo.mls.Constants
import com.immomo.mls.MLSBundleUtils
import com.immomo.mls.activity.LuaViewActivity
import com.immomo.mls.utils.MLSUtils
import xyz.hoper.dart.momo.activity.AssetsChooserActivity
import java.io.File
import java.io.IOException
import java.util.*

/**
 * Created by Xiong.Fangyu on 2019-08-01
 */
@SuppressLint("Registered")
class AssetsChooserActivity : ListActivity() {
    private var folderName: String? = null
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        val i = intent
        folderName = i.getStringExtra(KEY)
        if (folderName == null) {
            folderName = ""
        }
        initContent()
    }

    private fun initContent() {
        val adapter = ArrayAdapter(this, R.layout.simple_expandable_list_item_1, contentData)
        listView.adapter = adapter
        listView.onItemClickListener = AdapterView.OnItemClickListener { parent, view, position, id ->
            val fileName = getFileName(adapter.getItem(position - listView.headerViewsCount))
            val intent = Intent(this@AssetsChooserActivity, LuaViewActivity::class.java)
            val initData = MLSBundleUtils.createInitData(fileName, false)
            initData.forceNotUseX64()
            intent.putExtras(MLSBundleUtils.createBundle(initData))
            startActivity(intent)
        }
    }

    //                        PreloadUtils.preload(getFileName(name));
    private val contentData: List<String>
        private get() {
            var array: Array<String>? = null
            val result: MutableList<String> = ArrayList()
            try {
                array = resources.assets.list(folderName!!)
                if (array != null) {
                    for (name in array) {
                        if (MLSUtils.isLuaScript(name)) {
                            result.add(name)
                            //                        PreloadUtils.preload(getFileName(name));
                        }
                    }
                }
            } catch (e: IOException) {
                e.printStackTrace()
            }
            return result
        }

    private fun getFileName(n: String?): String {
        return Constants.ASSETS_PREFIX + if (TextUtils.isEmpty(folderName)) n else folderName + File.separator + n
    }

    companion object {
        private const val KEY = "__KEY"
        fun startActivity(a: Activity, dir: String?) {
            val i = Intent(a, AssetsChooserActivity::class.java)
            i.putExtra(KEY, dir)
            a.startActivity(i)
        }
    }
}