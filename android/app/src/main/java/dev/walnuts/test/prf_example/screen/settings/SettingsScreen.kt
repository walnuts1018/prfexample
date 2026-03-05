package dev.walnuts.test.prf_example.screen.settings

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Logout
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material3.AlertDialog
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilledTonalButton
import androidx.compose.material3.Icon
import androidx.compose.material3.LocalContentColor
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import dev.walnuts.test.prf_example.R
import dev.walnuts.test.prf_example.screen.components.DetailCard

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    viewModel: SettingsViewModel,
    modifier: Modifier = Modifier,
) {
    val state by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.nav_settings)) },
            )
        },
        modifier = modifier.fillMaxSize(),
    ) { innerPadding ->
        Column(
            modifier =
                Modifier
                    .fillMaxSize()
                    .padding(innerPadding)
                    .verticalScroll(rememberScrollState())
                    .padding(horizontal = 16.dp, vertical = 16.dp),
            verticalArrangement = Arrangement.spacedBy(24.dp),
        ) {
            Text(
                text = stringResource(R.string.passkey_info_title),
                style = MaterialTheme.typography.titleMedium,
                modifier = Modifier.fillMaxWidth(),
            )

            Column(
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                state.prfOutput?.let {
                    DetailCard(
                        label = stringResource(R.string.label_prf_key),
                        value = it.toHexString(),
                    )
                }

                state.userId?.let {
                    DetailCard(
                        label = stringResource(R.string.user_id),
                        value = it,
                        useMonospace = true,
                    )
                }

                DetailCard(
                    label = stringResource(R.string.label_prf_support),
                    value =
                        when (state.prfSupported) {
                            true -> stringResource(R.string.prf_supported)
                            false -> stringResource(R.string.prf_not_supported)
                            null -> stringResource(R.string.prf_unknown)
                        },
                    enableCopyButton = false,
                )

                state.credentialId?.let {
                    DetailCard(
                        label = stringResource(R.string.label_credential_id),
                        value = it,
                        useMonospace = true,
                    )
                }

                state.registrationResponseJson?.let {
                    DetailCard(
                        label = stringResource(R.string.label_registration_response),
                        value = it,
                        useMonospace = true,
                    )
                }
            }

            Column(
                verticalArrangement = Arrangement.spacedBy(8.dp),
            ) {
                FilledTonalButton(
                    onClick = viewModel::onSignOutClick,
                    modifier = Modifier.fillMaxWidth(),
                ) {
                    Icon(
                        imageVector = Icons.AutoMirrored.Filled.Logout,
                        contentDescription = null,
                        modifier = Modifier.padding(end = 8.dp),
                    )
                    Text(stringResource(R.string.sign_out))
                }

                FilledTonalButton(
                    onClick = viewModel::onDeleteClick,
                    enabled = !state.isDeleting,
                    modifier = Modifier.fillMaxWidth(),
                    colors =
                        ButtonDefaults.filledTonalButtonColors(
                            containerColor = MaterialTheme.colorScheme.errorContainer,
                            contentColor = MaterialTheme.colorScheme.onErrorContainer,
                        ),
                ) {
                    if (state.isDeleting) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(24.dp),
                            color = LocalContentColor.current,
                            strokeWidth = 2.dp,
                        )
                    } else {
                        Icon(
                            imageVector = Icons.Filled.Delete,
                            contentDescription = null,
                            modifier = Modifier.padding(end = 8.dp),
                        )
                        Text(stringResource(R.string.delete_passkey))
                    }
                }
            }
        }
    }

    when (state.dialog) {
        is SettingsViewModel.SettingsUiState.DialogState.Closed -> {}

        is SettingsViewModel.SettingsUiState.DialogState.ConfirmSignOut -> {
            AlertDialog(
                onDismissRequest = {
                    viewModel.closeDialog()
                },
                title = { Text(stringResource(R.string.sign_out)) },
                text = { Text(stringResource(R.string.sign_out_confirm_message)) },
                confirmButton = {
                    TextButton(
                        onClick = {
                            viewModel.closeDialog()
                            viewModel.signOut()
                        },
                    ) {
                        Text(stringResource(R.string.sign_out))
                    }
                },
                dismissButton = {
                    TextButton(onClick = {
                        viewModel.closeDialog()
                    }) {
                        Text(stringResource(R.string.cancel))
                    }
                },
            )
        }

        is SettingsViewModel.SettingsUiState.DialogState.ConfirmDelete -> {
            AlertDialog(
                onDismissRequest = { viewModel.closeDialog() },
                title = { Text(stringResource(R.string.delete_passkey)) },
                text = {
                    Text(stringResource(R.string.delete_confirm_message))
                },
                confirmButton = {
                    TextButton(
                        onClick = {
                            viewModel.closeDialog()
                            viewModel.deletePasskey()
                        },
                    ) {
                        Text(
                            stringResource(R.string.delete),
                            color = MaterialTheme.colorScheme.error,
                        )
                    }
                },
                dismissButton = {
                    TextButton(onClick = {
                        viewModel.closeDialog()
                    }) {
                        Text(stringResource(R.string.cancel))
                    }
                },
            )
        }
    }
}