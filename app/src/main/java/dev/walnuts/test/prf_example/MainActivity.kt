package dev.walnuts.test.prf_example

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.credentials.CredentialManager
import dev.walnuts.test.prf_example.api.ApiClient
import dev.walnuts.test.prf_example.screen.decrypt.DecryptViewModel
import dev.walnuts.test.prf_example.screen.encrypt.EncryptViewModel
import dev.walnuts.test.prf_example.screen.register.RegisterViewModel
import dev.walnuts.test.prf_example.screen.settings.SettingsViewModel
import dev.walnuts.test.prf_example.ui.theme.PRFExampleTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        val credentialManager = CredentialManager.create(this)
        val apiClient = ApiClient(baseUrl = "https://prfexample.walnuts.dev")
        val repository = PasskeyRepository(credentialManager, apiClient, applicationContext)
        val registerViewModel = RegisterViewModel(repository)
        val encryptViewModel = EncryptViewModel(repository)
        val decryptViewModel = DecryptViewModel(repository)
        val settingsViewModel = SettingsViewModel(repository)

        setContent {
            PRFExampleTheme {
                App(
                    repository = repository,
                    registerViewModel = registerViewModel,
                    encryptViewModel = encryptViewModel,
                    decryptViewModel = decryptViewModel,
                    settingsViewModel = settingsViewModel,
                )
            }
        }
    }
}
