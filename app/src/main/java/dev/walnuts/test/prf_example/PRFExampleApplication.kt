package dev.walnuts.test.prf_example

import android.app.Application
import timber.log.Timber

class PRFExampleApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        }
    }
}
