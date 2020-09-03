package xyz.hoper.dart

import android.content.Context
import android.os.Bundle
import android.util.Log
import android.view.View
import io.flutter.embedding.android.SplashScreen


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
