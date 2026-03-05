package dev.walnuts.test.prf_example.screen.settings

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import dagger.hilt.android.lifecycle.HiltViewModel
import dev.walnuts.test.prf_example.PasskeyRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class SettingsViewModel
    @Inject
    constructor(
        private val repository: PasskeyRepository,
    ) : ViewModel() {
        private val _uiState = MutableStateFlow(SettingsUiState())
        val uiState: StateFlow<SettingsUiState> = _uiState.asStateFlow()

        init {
            viewModelScope.launch {
                repository.registrationState.collect { regState ->
                    _uiState.update {
                        it.copy(
                            userId = regState.passkeyInfo?.userId,
                            credentialId = regState.passkeyInfo?.credentialId,
                            prfSupported = regState.passkeyInfo?.prfSupported,
                            registrationResponseJson = regState.passkeyInfo?.registrationResponseJson,
                        )
                    }
                }
            }
            viewModelScope.launch {
                repository.prfOutput.collect { prf ->
                    _uiState.update { it.copy(prfOutput = prf) }
                }
            }
        }

        // TODO: サーバー側のcredential削除APIを呼び出す
        // TODO: Credential Manager経由でプラットフォーム上のパスキーも削除する
        fun deletePasskey() {
            viewModelScope.launch {
                _uiState.update { it.copy(isDeleting = true) }
                try {
                    repository.deleteAll()
                    _uiState.value = SettingsUiState()
                } catch (e: Exception) {
                    _uiState.update { it.copy(isDeleting = false) }
                }
            }
        }

        fun signOut() {
            repository.signOut()
            _uiState.value = SettingsUiState()
        }

        fun onDeleteClick() {
            _uiState.update { it.copy(dialog = SettingsUiState.DialogState.ConfirmDelete) }
        }

        fun onSignOutClick() {
            _uiState.update { it.copy(dialog = SettingsUiState.DialogState.ConfirmSignOut) }
        }

        fun closeDialog() {
            _uiState.update { it.copy(dialog = SettingsUiState.DialogState.Closed) }
        }

        data class SettingsUiState(
            val userId: String? = null,
            val credentialId: String? = null,
            val prfSupported: Boolean? = null,
            val prfOutput: ByteArray? = null,
            val registrationResponseJson: String? = null,
            val isDeleting: Boolean = false,
            val dialog: DialogState = DialogState.Closed,
        ) {
            sealed interface DialogState {
                data object Closed : DialogState

                data object ConfirmDelete : DialogState

                data object ConfirmSignOut : DialogState
            }

            override fun equals(other: Any?): Boolean {
                if (this === other) return true
                if (javaClass != other?.javaClass) return false

                other as SettingsUiState

                if (prfSupported != other.prfSupported) return false
                if (isDeleting != other.isDeleting) return false
                if (userId != other.userId) return false
                if (credentialId != other.credentialId) return false
                if (!prfOutput.contentEquals(other.prfOutput)) return false
                if (registrationResponseJson != other.registrationResponseJson) return false
                if (dialog != other.dialog) return false

                return true
            }

            override fun hashCode(): Int {
                var result = prfSupported?.hashCode() ?: 0
                result = 31 * result + isDeleting.hashCode()
                result = 31 * result + (userId?.hashCode() ?: 0)
                result = 31 * result + (credentialId?.hashCode() ?: 0)
                result = 31 * result + (prfOutput?.contentHashCode() ?: 0)
                result = 31 * result + (registrationResponseJson?.hashCode() ?: 0)
                result = 31 * result + dialog.hashCode()
                return result
            }
        }
    }
