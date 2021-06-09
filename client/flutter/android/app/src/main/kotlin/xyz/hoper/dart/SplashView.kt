package xyz.hoper.dart

import android.content.Context
import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.os.Bundle
import android.os.Handler
import android.os.Message
import android.util.AttributeSet
import android.util.Log
import android.view.View
import android.widget.Toast
import androidx.appcompat.widget.AppCompatImageView
import io.flutter.embedding.android.SplashScreen
import java.io.*
import java.lang.ref.WeakReference
import java.net.HttpURLConnection
import java.net.URL
import kotlin.math.roundToInt

class SplashScreenWithTransition : SplashScreen {
    private lateinit var view: SplashView

    override fun createSplashView(
        context: Context,
        savedInstanceState: Bundle?
    ): View? {
        // A reference to the MySplashView is retained so that it can be told
        // to transition away at the appropriate time.
        view =  SplashView(context)
        Log.i("SplashScreen","SplashScreenWithTransition")
        return view
    }

    override fun transitionToFlutter(onTransitionComplete: Runnable) {
        // Instruct MySplashView to animate away in whatever manner it wants.
        // The onTransitionComplete Runnable is passed to the MySplashView to be
        // invoked when the transition animation is complete.
        view.animateAway(onTransitionComplete)
    }
}


class SplashView: AppCompatImageView {
    companion object{
        const val TAG = "SplashView"
        const val GET_DATA_SUCCESS = 1
        const val NETWORK_ERROR = 2
        const val SERVER_ERROR = 3
        const val ImagePath = "http://liov.xyz/static/images/6cbeb5c8-7160-4b6f-a342-d96d3c00367a.jpg"
    }

    constructor(context: Context):super(context){
        setImageBitmap()
    }
    constructor(context: Context, attrs: AttributeSet?):super(context, attrs)
    constructor(context: Context, attrs: AttributeSet?, defStyleAttr: Int):super(context, attrs, defStyleAttr)

    //子线程不能操作UI，通过Handler设置图片
    class ViewHandler(view:SplashView):Handler(){
        private val viewRef = WeakReference(view)
        override fun handleMessage(msg: Message) {
            super.handleMessage(msg)
            val view = viewRef.get()!!
            when (msg.what) {
                GET_DATA_SUCCESS -> {
                    val bitmap = msg.obj as Bitmap
                    view.setImageBitmap(bitmap)
                }
                NETWORK_ERROR -> Toast.makeText(view.context, "网络连接失败", Toast.LENGTH_SHORT).show()
                SERVER_ERROR -> Toast.makeText(view.context, "服务器发生错误", Toast.LENGTH_SHORT).show()
            }
        }
    }

    private val handler = ViewHandler(this)

    private fun setImageBitmap(){
        object:Thread(){
            override fun run() {
                useCacheImage()
            }
        }.start()
    }

    //设置网络图片
    private fun useNetWorkImage(path: String) {
                try {
                    //获取连接
                    val connection = URL(path).openConnection() as HttpURLConnection
                    //使用GET方法访问网络
                    connection.requestMethod = "GET"
                    //超时时间为10秒
                    connection.connectTimeout = 10000
                    //获取返回码
                    val code = connection.responseCode
                    if (code == 200) {
                        val inputStream = connection.inputStream

                        val baos = ByteArrayOutputStream()
                        try {
                            val buffer = ByteArray(1024)
                            var len: Int
                            while (inputStream.read(buffer).also { len = it } > -1) {
                                baos.write(buffer, 0, len)
                            }
                            baos.flush()
                        } catch (e: IOException) {
                            e.printStackTrace()
                        }

                        //复制新的输入流
                        val inputStream1 = ByteArrayInputStream(baos.toByteArray())
                        val inputStream2 = ByteArrayInputStream(baos.toByteArray())

                        //调用压缩方法显示图片
                        val bitmap = getCompressBitmap(inputStream1)
                        //调用缓存图片方法
                        cacheImage(inputStream2)
                        val msg = Message.obtain().apply {
                            obj = bitmap
                            what = GET_DATA_SUCCESS
                        }
                        handler.sendMessage(msg)
                        inputStream.close()
                    }else {
                        //服务启发生错误
                        Toast.makeText(context, "服务器发生错误", Toast.LENGTH_SHORT).show()
                    }
                } catch (e: IOException) {
                    e.printStackTrace()
                    //网络连接错误
                    Toast.makeText(context, "网络错误", Toast.LENGTH_SHORT).show()
                }
            }


    //使用缓存图片
    fun useCacheImage() {
        //创建路径一样的文件
        val file = File(context.cacheDir, getURLPath())
        //判断文件是否存在
        if (file.length() > 0) {
            //使用本地图片
            try {
                val inputStream = FileInputStream(file)
                //调用压缩方法显示图片
                val bitmap = getCompressBitmap(inputStream)
                //利用Message把图片发给Handler
                val msg = Message.obtain().apply {
                    obj = bitmap
                    what = GET_DATA_SUCCESS
                }
                handler.sendMessage(msg)
                Log.e(TAG, "使用缓存图片")
            } catch (e: FileNotFoundException) {
                e.printStackTrace()
            }
        }else{
            useNetWorkImage(ImagePath)
        }
    }

    /**
     * 缓存网络的图片
     * @param inputStream 网络的输入流
     */
    fun cacheImage(inputStream: InputStream) {
        try {
            val file = File(context.cacheDir, getURLPath())
            val fos = FileOutputStream(file)
            var len: Int
            val buffer = ByteArray(1024)
            while (inputStream.read(buffer).also { len = it } != -1) {
                fos.write(buffer, 0, len)
            }
            fos.close()
            Log.e(TAG, "缓存成功")
        } catch (e: IOException) {
            e.printStackTrace()
            Log.e(TAG, "缓存失败")
        }
    }

    /**
     * 根据网址生成一个文件名
     * @return 文件名
     */
    private fun getURLPath(): String {
        val urlStr2 = StringBuilder()
        val strings = ImagePath.split("/")
        for (string in strings) {
            urlStr2.append(string, "_")
        }
        Log.e(TAG, "文件名：$urlStr2")
        return urlStr2.toString()
    }

    /**
     * 根据输入流返回一个压缩的图片
     *
     * @param input 图片的输入流
     * @return 压缩的图片
     */
    fun getCompressBitmap(input: InputStream): Bitmap? {
        //因为InputStream要使用两次，但是使用一次就无效了，所以需要复制两个
        val baos = ByteArrayOutputStream()
        try {
            val buffer = ByteArray(1024)
            var len: Int
            while (input.read(buffer).also { len = it } > -1) {
                baos.write(buffer, 0, len)
            }
            baos.flush()
        } catch (e: IOException) {
            e.printStackTrace()
        }

        //复制新的输入流
        val `is`: InputStream = ByteArrayInputStream(baos.toByteArray())
        val is2: InputStream = ByteArrayInputStream(baos.toByteArray())

        //只是获取网络图片的大小，并没有真正获取图片
        val options: BitmapFactory.Options = BitmapFactory.Options()
        options.inJustDecodeBounds = true
        BitmapFactory.decodeStream(`is`, null, options)
        //获取图片并进行压缩
        options.inSampleSize = getInSampleSize(options)
        options.inJustDecodeBounds = false
        return BitmapFactory.decodeStream(is2, null, options)
    }

    /**
     * 获得需要压缩的比率
     *
     * @param options 需要传入已经BitmapFactory.decodeStream(is, null, options);
     * @return 返回压缩的比率，最小为1
     */
    private fun getInSampleSize(options: BitmapFactory.Options): Int {
        var inSampleSize = 1
        val realWith: Int = realImageViewWith()
        val realHeight: Int = realImageViewHeight()
        val outWidth = options.outWidth
        Log.i("网络图片实际的宽度", outWidth.toString())
        val outHeight = options.outHeight
        Log.i("网络图片实际的高度", outHeight.toString())

        //获取比率最大的那个
        if (outWidth > realWith || outHeight > realHeight) {
            val withRadio = (outWidth / realWith.toFloat()).roundToInt()
            val heightRadio = (outHeight / realHeight.toFloat()).roundToInt()
            inSampleSize = if (withRadio > heightRadio) withRadio else heightRadio
        }
        Log.i("压缩比率", inSampleSize.toString())
        return inSampleSize
    }

    /**
     * 获取ImageView实际的宽度
     *
     * @return 返回ImageView实际的宽度
     */
    private fun realImageViewWith(): Int {

        //如果ImageView设置了宽度就可以获取实在宽带
        var width = width
        if (width <= 0)  width = layoutParams.width //如果ImageView没有设置宽度，就获取父级容器的宽度

        if (width <= 0) width = maxWidth //获取ImageView宽度的最大值

        if (width <= 0) width = context.resources.displayMetrics.widthPixels //获取屏幕的宽度


        Log.i("ImageView实际的宽度", width.toString())
        return width
    }

    /**
     * 获取ImageView实际的高度
     *
     * @return 返回ImageView实际的高度
     */
    fun realImageViewHeight(): Int {

        //如果ImageView设置了高度就可以获取实在宽度
        var height = height
        if (height <= 0) height = layoutParams.height //如果ImageView没有设置高度，就获取父级容器的高度


        if (height <= 0) height = maxHeight //获取ImageView高度的最大值


        if (height <= 0) height = context.resources.displayMetrics.heightPixels //获取ImageView高度的最大值


        Log.i("ImageView实际的高度", height.toString())
        return height
    }


    fun animateAway(onTransitionComplete: Runnable){
        onTransitionComplete.run()
    }
}