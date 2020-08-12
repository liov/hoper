package xyz.hoper.dart

import android.content.pm.PackageManager
import android.graphics.Color
import android.graphics.drawable.Drawable
import android.os.Build
import android.os.Bundle
import android.view.View
import android.view.Window
import android.view.WindowManager
import android.widget.ImageView
import androidx.appcompat.app.AppCompatActivity
import com.idlefish.flutterboost.containers.FlutterFragment
import io.flutter.embedding.android.DrawableSplashScreen
import io.flutter.embedding.android.SplashScreen
import io.flutter.embedding.android.SplashScreenProvider
import io.flutter.plugin.platform.PlatformPlugin

class FlutterFragmentPageActivity : AppCompatActivity(), View.OnClickListener, SplashScreenProvider {
    private var mFragment: FlutterFragment? = null
    private lateinit var mTab1: View
    private lateinit var mTab2: View
    private lateinit var mTab3: View
    private lateinit var mTab4: View
    override fun onCreate(savedInstanceState: Bundle?) {
        supportRequestWindowFeature(Window.FEATURE_NO_TITLE)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            val window = window
            window.addFlags(WindowManager.LayoutParams.FLAG_DRAWS_SYSTEM_BAR_BACKGROUNDS)
            window.statusBarColor = 0x40000000
            window.decorView.systemUiVisibility = PlatformPlugin.DEFAULT_SYSTEM_UI
        }
        super.onCreate(savedInstanceState)
        val actionBar = supportActionBar
        actionBar?.hide()
        setContentView(R.layout.flutter_fragment_page)
        mTab1 = findViewById(R.id.tab1)
        mTab2 = findViewById(R.id.tab2)
        mTab3 = findViewById(R.id.tab3)
        mTab4 = findViewById(R.id.tab4)
        mTab1.setOnClickListener(this)
        mTab2.setOnClickListener(this)
        mTab3.setOnClickListener(this)
        mTab4.setOnClickListener(this)
    }

    override fun onClick(v: View) {
        mTab1.setBackgroundColor(Color.WHITE)
        mTab2.setBackgroundColor(Color.WHITE)
        mTab3.setBackgroundColor(Color.WHITE)
        mTab4.setBackgroundColor(Color.WHITE)
        mFragment = if (mTab1 === v) {
            mTab1.setBackgroundColor(Color.YELLOW)
            FlutterFragment.NewEngineFragmentBuilder().url("flutterFragment").build()
        } else if (mTab2 === v) {
            mTab2.setBackgroundColor(Color.YELLOW)
            FlutterFragment.NewEngineFragmentBuilder().url("flutterFragment").build()
        } else if (mTab3 === v) {
            mTab3.setBackgroundColor(Color.YELLOW)
            FlutterFragment.NewEngineFragmentBuilder().url("flutterFragment").build()
        } else {
            mTab4.setBackgroundColor(Color.YELLOW)
            FlutterFragment.NewEngineFragmentBuilder().url("flutterFragment").build()
        }
        supportFragmentManager
                .beginTransaction()
                .replace(R.id.fragment_stub, mFragment!!)
                .commit()
    }

    override fun onResume() {
        super.onResume()
        mTab1!!.performClick()
    }

    override fun provideSplashScreen(): SplashScreen? {
        val manifestSplashDrawable = splashScreenFromManifest
        return if (manifestSplashDrawable != null) {
            DrawableSplashScreen(manifestSplashDrawable, ImageView.ScaleType.CENTER, 500L)
        } else {
            null
        }
    }

    // This is never expected to happen.
    private val splashScreenFromManifest: Drawable?
        private get() = try {
            val activityInfo = packageManager.getActivityInfo(
                    componentName,
                    PackageManager.GET_META_DATA or PackageManager.GET_ACTIVITIES
            )
            val metadata = activityInfo.metaData
            val splashScreenId = metadata?.getInt(SPLASH_SCREEN_META_DATA_KEY)
            if (splashScreenId != null) if (Build.VERSION.SDK_INT > Build.VERSION_CODES.LOLLIPOP) resources.getDrawable(splashScreenId, theme) else resources.getDrawable(splashScreenId) else null
        } catch (e: PackageManager.NameNotFoundException) {
            // This is never expected to happen.
            null
        }

    companion object {
        protected const val SPLASH_SCREEN_META_DATA_KEY = "io.flutter.embedding.android.SplashScreenDrawable"
    }
}