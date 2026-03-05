package dev.walnuts.test.prf_example

import android.app.Application
import androidx.credentials.CredentialManager
import dev.walnuts.test.prf_example.api.ApiClient
import timber.log.Timber

class PRFExampleApplication : Application() {
    val credentialManager: CredentialManager by lazy { CredentialManager.create(this) }
    val apiClient: ApiClient by lazy { ApiClient(baseUrl = "https://prfexample.walnuts.dev") }
    val repository: PasskeyRepository by lazy { PasskeyRepository(credentialManager, apiClient, this) }

    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        }
    }
}
