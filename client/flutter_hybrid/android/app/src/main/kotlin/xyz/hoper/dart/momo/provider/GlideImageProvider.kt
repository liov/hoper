/**
 * Created by MomoLuaNative.
 * Copyright (c) 2019, Momo Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
/*
 * Created by LuaView.
 * Copyright (c) 2017, Alibaba Group. All rights reserved.
 *
 * This source code is licensed under the MIT.
 * For the full copyright and license information,please view the LICENSE file in the root directory of this source tree.
 */
package xyz.hoper.dart.momo.provider

import android.app.Activity
import android.content.Context
import android.graphics.RectF
import android.graphics.drawable.Drawable
import android.text.TextUtils
import android.view.ViewGroup
import android.widget.ImageView
import androidx.collection.LruCache
import com.bumptech.glide.Glide
import com.bumptech.glide.RequestBuilder
import com.bumptech.glide.load.DataSource
import com.bumptech.glide.load.engine.GlideException
import com.bumptech.glide.request.RequestListener
import com.bumptech.glide.request.RequestOptions
import com.bumptech.glide.request.target.Target
import com.immomo.mls.provider.DrawableLoadCallback
import com.immomo.mls.provider.ImageProvider
import xyz.hoper.dart.momo.provider.ResourcesUtils.getResourceIdByName
import xyz.hoper.dart.momo.provider.ResourcesUtils.getResourceIdByUrl
import java.lang.ref.WeakReference

/**
 * XXX
 *
 * @author song
 * @date 16/4/11
 * 主要功能描述
 * 修改描述
 * 下午5:42 song XXX
 */
class GlideImageProvider : ImageProvider {
    override fun pauseRequests(view: ViewGroup, context: Context) {
        Glide.with(context).pauseRequests()
    }

    override fun resumeRequests(view: ViewGroup, context: Context) {
        if (context is Activity && (context.isFinishing || context.isDestroyed)) {
            return
        }
        if (Glide.with(context).isPaused) {
            Glide.with(context).resumeRequests()
        }
    }

    /**
     * load url
     * @param url
     * @param placeHolder
     * @param callback
     */
    override fun load(context: Context, imageView: ImageView, url: String,
                      placeHolder: String?, radius: RectF?, callback: DrawableLoadCallback?) {
        var builder: RequestBuilder<*>
        builder = if (callback != null) {
            val cf = WeakReference(callback)
            Glide.with(context).load(url).listener(object : RequestListener<Drawable?> {
                override fun onLoadFailed(e: GlideException?, model: Any, target: Target<Drawable?>, isFirstResource: Boolean): Boolean {
                    e?.printStackTrace()
                    if (cf.get() != null) {
                        cf.get()!!.onLoadResult(null, e?.message)
                    }
                    return false
                }

                override fun onResourceReady(resource: Drawable?, model: Any, target: Target<Drawable?>, dataSource: DataSource, isFirstResource: Boolean): Boolean {
                    if (cf.get() != null) {
                        cf.get()!!.onLoadResult(resource, null)
                    }
                    return false
                }
            })
        } else {
            Glide.with(context).load(url)
        }
        if (placeHolder != null) {
            val id = getResourceIdByUrl(placeHolder, null, ResourcesUtils.TYPE.DRAWABLE)
            if (id > 0) {
                builder = builder.apply(RequestOptions().placeholder(id))
            }
        }
        builder.into(imageView)
    }

    override fun loadWithoutInterrupt(context: Context, iv: ImageView, url: String,
                                      placeHolder: String?, radius: RectF?, callback: DrawableLoadCallback?) {
        load(context.applicationContext, iv, url, placeHolder, radius, callback)
    }

    override fun loadProjectImage(context: Context, name: String): Drawable? {
        if (TextUtils.isEmpty(name)) return null
        val id = getProjectImageId(name)
        return if (id > 0) {
            context.resources.getDrawable(id)
        } else null
    }

    override fun preload(context: Context, url: String, radius: RectF?, callback: DrawableLoadCallback?) {
        val builder: RequestBuilder<*>
        builder = if (callback != null) {
            Glide.with(context).load(url).listener(object : RequestListener<Drawable?> {
                override fun onLoadFailed(e: GlideException?, model: Any, target: Target<Drawable?>, isFirstResource: Boolean): Boolean {
                    callback.onLoadResult(null, e?.message)
                    return false
                }

                override fun onResourceReady(resource: Drawable?, model: Any, target: Target<Drawable?>, dataSource: DataSource, isFirstResource: Boolean): Boolean {
                    callback.onLoadResult(resource, null)
                    return false
                }
            })
        } else {
            Glide.with(context).load(url)
        }
        builder.preload()
    }

    companion object {
        private val IdCache = LruCache<String, Int>(50)
        private fun getProjectImageId(name: String): Int {
            val result = IdCache[name]
            if (result != null) {
                return result
            }
            val id = getResourceIdByName(name, ResourcesUtils.TYPE.DRAWABLE)
            IdCache.put(name, id)
            return id
        }
    }
}