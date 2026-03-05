package dev.walnuts.test.prf_example

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.viewModels
import dev.walnuts.test.prf_example.screen.decrypt.DecryptViewModel
import dev.walnuts.test.prf_example.screen.encrypt.EncryptViewModel
import dev.walnuts.test.prf_example.screen.register.RegisterViewModel
import dev.walnuts.test.prf_example.screen.settings.SettingsViewModel
import dev.walnuts.test.prf_example.ui.theme.PRFExampleTheme

class MainActivity : ComponentActivity() {
    private val app by lazy { application as PRFExampleApplication }

    private val registerViewModel: RegisterViewModel by viewModels {
        RegisterViewModel.provideFactory(app.repository)
    }
    private val encryptViewModel: EncryptViewModel by viewModels {
        EncryptViewModel.provideFactory(app.repository)
    }
    private val decryptViewModel: DecryptViewModel by viewModels {
        DecryptViewModel.provideFactory(app.repository)
    }
    private val settingsViewModel: SettingsViewModel by viewModels {
        SettingsViewModel.provideFactory(app.repository)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        setContent {
            PRFExampleTheme {
                App(
                    repository = app.repository,
                    registerViewModel = registerViewModel,
                    encryptViewModel = encryptViewModel,
                    decryptViewModel = decryptViewModel,
                    settingsViewModel = settingsViewModel,
                )
            }
        }
    }
}
