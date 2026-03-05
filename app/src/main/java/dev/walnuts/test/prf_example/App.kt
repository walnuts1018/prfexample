package dev.walnuts.test.prf_example

import androidx.annotation.StringRes
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.consumeWindowInsets
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material.icons.filled.LockOpen
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHost
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.stringResource
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation3.runtime.NavKey
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.runtime.rememberNavBackStack
import androidx.navigation3.ui.NavDisplay
import dev.walnuts.test.prf_example.screen.decrypt.DecryptScreen
import dev.walnuts.test.prf_example.screen.decrypt.DecryptViewModel
import dev.walnuts.test.prf_example.screen.encrypt.EncryptScreen
import dev.walnuts.test.prf_example.screen.encrypt.EncryptViewModel
import dev.walnuts.test.prf_example.screen.register.RegisterScreen
import dev.walnuts.test.prf_example.screen.register.RegisterViewModel
import dev.walnuts.test.prf_example.screen.settings.SettingsScreen
import dev.walnuts.test.prf_example.screen.settings.SettingsViewModel
import kotlinx.serialization.Serializable

sealed interface TopLevelRoute : NavKey {
    @get:StringRes
    val labelResId: Int
    val icon: ImageVector

    @Serializable
    data object Encrypt : TopLevelRoute {
        override val labelResId = R.string.nav_encrypt
        override val icon = Icons.Filled.Lock
    }

    @Serializable
    data object Decrypt : TopLevelRoute {
        override val labelResId = R.string.nav_decrypt
        override val icon = Icons.Filled.LockOpen
    }

    @Serializable
    data object Settings : TopLevelRoute {
        override val labelResId = R.string.nav_settings
        override val icon = Icons.Filled.Settings
    }
}

val topLevelRoutes = listOf(TopLevelRoute.Encrypt, TopLevelRoute.Decrypt, TopLevelRoute.Settings)

@Composable
internal fun App(modifier: Modifier = Modifier) {
    val registerViewModel: RegisterViewModel = hiltViewModel()
    val registrationState by registerViewModel.registrationState.collectAsStateWithLifecycle()

    if (registrationState.isRegistered) {
        MainApp(modifier = modifier)
    } else {
        RegisterScreen(
            viewModel = registerViewModel,
            modifier = modifier,
        )
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun MainApp(modifier: Modifier = Modifier) {
    val encryptViewModel: EncryptViewModel = hiltViewModel()
    val decryptViewModel: DecryptViewModel = hiltViewModel()
    val settingsViewModel: SettingsViewModel = hiltViewModel()
    val backStack = rememberNavBackStack(TopLevelRoute.Encrypt)
    val snackbarHostState = remember { SnackbarHostState() }
    val currentRoute = backStack.lastOrNull() as? TopLevelRoute

    Scaffold(
        modifier = modifier,
        snackbarHost = { SnackbarHost(hostState = snackbarHostState) },
        bottomBar = {
            NavigationBar {
                topLevelRoutes.forEach { route ->
                    val label = stringResource(route.labelResId)
                    NavigationBarItem(
                        icon = {
                            Icon(
                                imageVector = route.icon,
                                contentDescription = label,
                            )
                        },
                        label = { Text(label) },
                        selected = currentRoute == route,
                        onClick = {
                            if (currentRoute != route) {
                                backStack.removeAll { it != TopLevelRoute.Encrypt }
                                if (route != TopLevelRoute.Encrypt) {
                                    backStack.add(route)
                                }
                            }
                        },
                    )
                }
            }
        },
    ) { innerPadding ->
        val bottomPadding = PaddingValues(bottom = innerPadding.calculateBottomPadding())
        NavDisplay(
            backStack = backStack,
            onBack = { backStack.removeAt(backStack.lastIndex) },
            modifier =
                Modifier
                    .padding(bottomPadding)
                    .consumeWindowInsets(bottomPadding),
            entryProvider =
                entryProvider {
                    entry<TopLevelRoute.Encrypt> {
                        EncryptScreen(
                            viewModel = encryptViewModel,
                            snackbarHostState = snackbarHostState,
                        )
                    }
                    entry<TopLevelRoute.Decrypt> {
                        DecryptScreen(
                            viewModel = decryptViewModel,
                            snackbarHostState = snackbarHostState,
                        )
                    }
                    entry<TopLevelRoute.Settings> {
                        SettingsScreen(viewModel = settingsViewModel)
                    }
                },
        )
    }
}