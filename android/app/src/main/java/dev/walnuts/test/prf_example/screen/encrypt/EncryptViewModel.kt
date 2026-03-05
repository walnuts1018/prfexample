package dev.walnuts.test.prf_example.screen.encrypt

import android.util.Base64
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import dev.walnuts.test.prf_example.Crypto
import dev.walnuts.test.prf_example.LastEncryptResult
import dev.walnuts.test.prf_example.PasskeyRepository
import dev.walnuts.test.prf_example.R
import dev.walnuts.test.prf_example.UiMessage
import dev.walnuts.test.prf_example.toHexString
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

@HiltViewModel
class EncryptViewModel
    @Inject
    constructor(
        private val repository: PasskeyRepository,
    ) : ViewModel() {
        private val _uiState = MutableStateFlow(EncryptUiState())
        val uiState: StateFlow<EncryptUiState> = _uiState.asStateFlow()

        init {
            viewModelScope.launch {
                repository.registrationState.collect { regState ->
                    if (!regState.isRegistered) {
                        _uiState.value = EncryptUiState()
                    }
                }
            }
        }

        fun updatePlaintext(text: String) {
            _uiState.update { it.copy(plaintext = text) }
        }

        fun clearError() {
            _uiState.update { it.copy(errorMessage = null) }
        }

        fun encryptText() {
            val plaintext = _uiState.value.plaintext
            if (plaintext.isBlank()) {
                _uiState.update {
                    it.copy(errorMessage = UiMessage.StringResource(R.string.error_empty_plaintext))
                }
                return
            }

            val prfOutput = repository.getPrfOutputOrNull()
            if (prfOutput == null) {
                _uiState.update {
                    it.copy(errorMessage = UiMessage.StringResource(R.string.error_prf_key_missing))
                }
                return
            }

            viewModelScope.launch {
                _uiState.update {
                    it.copy(
                        isLoading = true,
                        errorMessage = null,
                        encryptedData = null,
                        ivHex = null,
                    )
                }

                try {
                    Timber.d("Encrypting with PRF key (${prfOutput.size} bytes)")

                    val encryptionResult = Crypto.encrypt(prfOutput, plaintext)

                    val ciphertextBase64 =
                        Base64.encodeToString(
                            encryptionResult.ciphertext,
                            Base64.NO_WRAP,
                        )
                    val ivHex = encryptionResult.iv.toHexString()

                    repository.saveEncryptedData(
                        ciphertextBytes = encryptionResult.ciphertext,
                        ivBytes = encryptionResult.iv,
                    )

                    val ivBase64 =
                        Base64.encodeToString(encryptionResult.iv, Base64.NO_WRAP)
                    repository.updateLastEncryptResult(
                        LastEncryptResult(
                            ciphertextBase64 = ciphertextBase64,
                            ivBase64 = ivBase64,
                        ),
                    )

                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            encryptedData = ciphertextBase64,
                            ivHex = ivHex,
                        )
                    }
                } catch (e: Exception) {
                    Timber.e(e, "Encryption failed")
                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            errorMessage =
                                UiMessage.StringResource(
                                    R.string.error_encryption_failed,
                                    listOf(e.message ?: ""),
                                ),
                        )
                    }
                }
            }
        }
    }

data class EncryptUiState(
    val isLoading: Boolean = false,
    val plaintext: String = "",
    val encryptedData: String? = null,
    val ivHex: String? = null,
    val errorMessage: UiMessage? = null,
)
