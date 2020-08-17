/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.provider

import android.content.res.Resources
import android.text.TextUtils
import androidx.collection.LruCache
import xyz.hoper.dart.App

/**
 * Created by XiongFangyu on 2018/8/9.
 */
object ResourcesUtils {
    private const val NINE_PATCH_PNG = ".9.png"
    private const val NORMAL_PNG = ".png"
    private val IdCache = LruCache<String, Int>(50)
    @JvmStatic
    fun getResourceIdByUrl(url: String, suffix: String?, type: TYPE?): Int {
        if (TextUtils.isEmpty(url) || type == null) return -1
        val key = url + type.toString()
        val result = IdCache[key]
        if (result != null) return result
        val start = url.lastIndexOf('/')
        val end = if (TextUtils.isEmpty(suffix)) url.length else url.lastIndexOf(suffix!!)
        if (start < 0 && end < 0 || end <= start) {
            IdCache.put(key, -1)
            return -1
        }
        val name = url.substring(start + 1, end)
        if (TextUtils.isEmpty(name)) {
            IdCache.put(key, -1)
            return -1
        }
        val id = getResourceIdByName(name, type)
        IdCache.put(key, id)
        return id
    }

    @JvmStatic
    fun getResourceIdByName(name: String?, type: TYPE): Int {
        return resources.getIdentifier(name, type.toString(), App.getPackageNameImpl())
    }

    fun getPngResourceIdByUrl(url: String): Int {
        return getResourceIdByUrl(url, NORMAL_PNG, TYPE.DRAWABLE)
    }

    fun getNinePatchResourceIdByUrl(url: String): Int {
        return getResourceIdByUrl(url, NINE_PATCH_PNG, TYPE.DRAWABLE)
    }

    private val resources: Resources
        get() = App.getInstance().resources

    enum class TYPE(private val s: String) {
        DRAWABLE("drawable"), ID("id"), DIMEN("dimen"), RAW("raw");

        override fun toString(): String {
            return s
        }

    }
}