package dev.walnuts.test.prf_example.screen.encrypt

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Lock
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.FilledTonalButton
import androidx.compose.material3.Icon
import androidx.compose.material3.LocalContentColor
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.SnackbarHostState
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalResources
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.unit.dp
import dev.walnuts.test.prf_example.R
import dev.walnuts.test.prf_example.UiMessage
import dev.walnuts.test.prf_example.screen.components.DetailCard

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun EncryptScreen(
    viewModel: EncryptViewModel,
    snackbarHostState: SnackbarHostState,
    modifier: Modifier = Modifier,
) {
    val uiState by viewModel.uiState.collectAsState()
    val resources = LocalResources.current

    LaunchedEffect(uiState.errorMessage) {
        val msg = uiState.errorMessage ?: return@LaunchedEffect
        val text =
            when (msg) {
                is UiMessage.StringResource -> {
                    if (msg.formatArgs.isEmpty()) {
                        resources.getString(msg.resId)
                    } else {
                        resources.getString(msg.resId, *msg.formatArgs.toTypedArray())
                    }
                }
            }
        snackbarHostState.showSnackbar(text)
        viewModel.clearError()
    }
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.nav_encrypt)) },
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
            Column(
                verticalArrangement = Arrangement.spacedBy(8.dp),
                modifier = Modifier.fillMaxWidth(),
            ) {
                OutlinedTextField(
                    value = uiState.plaintext,
                    onValueChange = { viewModel.updatePlaintext(it) },
                    label = { Text(stringResource(R.string.encrypt_input_label)) },
                    modifier = Modifier.fillMaxWidth(),
                    minLines = 3,
                    maxLines = 6,
                    enabled = !uiState.isLoading,
                )

                FilledTonalButton(
                    onClick = { viewModel.encryptText() },
                    enabled = !uiState.isLoading && uiState.plaintext.isNotBlank(),
                    modifier = Modifier.fillMaxWidth(),
                ) {
                    if (uiState.isLoading) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(24.dp),
                            color = LocalContentColor.current,
                            strokeWidth = 2.dp,
                        )
                    } else {
                        Icon(
                            imageVector = Icons.Filled.Lock,
                            contentDescription = null,
                            modifier = Modifier.padding(end = 8.dp),
                        )
                        Text(stringResource(R.string.encrypt_button))
                    }
                }
            }

            Column(
                verticalArrangement = Arrangement.spacedBy(12.dp),
                modifier = Modifier.fillMaxWidth(),
            ) {
                Text(
                    text = stringResource(R.string.encrypt_result_title),
                    style = MaterialTheme.typography.titleMedium,
                )

                DetailCard(
                    label = stringResource(R.string.label_iv_hex),
                    value = uiState.ivHex ?: "",
                    useMonospace = true,
                )

                DetailCard(
                    label = stringResource(R.string.label_ciphertext_base64),
                    value = uiState.encryptedData ?: "",
                    useMonospace = true,
                )
            }
        }
    }
}