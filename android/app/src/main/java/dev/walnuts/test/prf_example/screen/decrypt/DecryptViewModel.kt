package dev.walnuts.test.prf_example.screen.decrypt

import android.util.Base64
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import dev.walnuts.test.prf_example.Crypto
import dev.walnuts.test.prf_example.EncryptedRecord
import dev.walnuts.test.prf_example.LastEncryptResult
import dev.walnuts.test.prf_example.PasskeyRepository
import dev.walnuts.test.prf_example.R
import dev.walnuts.test.prf_example.UiMessage
import dev.walnuts.test.prf_example.toHexString
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharingStarted
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.stateIn
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import timber.log.Timber
import javax.inject.Inject

data class DecryptUiState(
    val isLoading: Boolean = false,
    val isRefreshing: Boolean = false,
    val encryptedInput: String = "",
    val ivInput: String = "",
    val decryptedText: String? = null,
    val errorMessage: UiMessage? = null,
)

@HiltViewModel
class DecryptViewModel
    @Inject
    constructor(
        private val repository: PasskeyRepository,
    ) : ViewModel() {
        private val _uiState = MutableStateFlow(DecryptUiState())
        val uiState: StateFlow<DecryptUiState> = _uiState.asStateFlow()

        val savedRecords: StateFlow<List<EncryptedRecord>> =
            repository.encryptedRecords
                .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5000), emptyList())

        val lastEncryptResult: StateFlow<LastEncryptResult?> =
            repository.lastEncryptResult
                .stateIn(viewModelScope, SharingStarted.WhileSubscribed(5000), null)

        init {
            viewModelScope.launch {
                repository.registrationState.collect { regState ->
                    if (!regState.isRegistered) {
                        _uiState.value = DecryptUiState()
                    }
                }
            }
        }

        fun refreshRecords() {
            viewModelScope.launch {
                _uiState.update { it.copy(isRefreshing = true) }
                try {
                    repository.refreshEncryptedRecords()
                } catch (e: Exception) {
                    Timber.e(e, "Failed to refresh records")
                    _uiState.update {
                        it.copy(
                            errorMessage =
                                UiMessage.StringResource(
                                    R.string.error_refresh_failed,
                                    listOf(e.message ?: ""),
                                ),
                        )
                    }
                } finally {
                    _uiState.update { it.copy(isRefreshing = false) }
                }
            }
        }

        fun updateEncryptedInput(text: String) {
            _uiState.update { it.copy(encryptedInput = text) }
        }

        fun updateIvInput(text: String) {
            _uiState.update { it.copy(ivInput = text) }
        }

        fun clearError() {
            _uiState.update { it.copy(errorMessage = null) }
        }

        fun copyEncryptResultToDecrypt() {
            val result = lastEncryptResult.value ?: return
            val ivHex = Base64.decode(result.ivBase64, Base64.NO_WRAP).toHexString()
            _uiState.update {
                it.copy(
                    encryptedInput = result.ciphertextBase64,
                    ivInput = ivHex,
                )
            }
        }

        fun loadEncryptedRecord(record: EncryptedRecord) {
            val ivHex = Base64.decode(record.ivBase64, Base64.NO_WRAP).toHexString()
            _uiState.update {
                it.copy(
                    encryptedInput = record.dataBase64,
                    ivInput = ivHex,
                    decryptedText = null,
                    errorMessage = null,
                )
            }
        }

        fun decryptText() {
            val state = _uiState.value
            if (state.encryptedInput.isBlank() || state.ivInput.isBlank()) {
                _uiState.update {
                    it.copy(errorMessage = UiMessage.StringResource(R.string.error_empty_decrypt_input))
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
                        decryptedText = null,
                    )
                }

                try {
                    Timber.d("Decrypting with PRF key (${prfOutput.size} bytes)")

                    val ciphertext = Base64.decode(state.encryptedInput, Base64.NO_WRAP)
                    val iv = decodeIvInput(state.ivInput)

                    val decryptedText = Crypto.decrypt(prfOutput, iv, ciphertext)

                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            decryptedText = decryptedText,
                        )
                    }
                } catch (e: Exception) {
                    Timber.e(e, "Decryption failed")
                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            errorMessage =
                                UiMessage.StringResource(
                                    R.string.error_decryption_failed,
                                    listOf(e.message ?: ""),
                                ),
                        )
                    }
                }
            }
        }

        private fun decodeIvInput(input: String): ByteArray = input.hexToByteArray()

        private fun String.hexToByteArray(): ByteArray {
            val hex = this.replace(" ", "")
            require(hex.length % 2 == 0) { "Hex string must have even length" }
            return ByteArray(hex.length / 2) { i ->
                hex.substring(i * 2, i * 2 + 2).toInt(16).toByte()
            }
        }
    }
