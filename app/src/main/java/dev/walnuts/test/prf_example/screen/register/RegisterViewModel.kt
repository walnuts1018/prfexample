package dev.walnuts.test.prf_example.screen.register

import androidx.credentials.exceptions.CreateCredentialException
import androidx.credentials.exceptions.GetCredentialException
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import dev.walnuts.test.prf_example.PasskeyRepository
import dev.walnuts.test.prf_example.R
import dev.walnuts.test.prf_example.UiMessage
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import timber.log.Timber

data class RegisterUiState(
    val isLoading: Boolean = false,
    val errorMessage: UiMessage? = null,
    val isLoginMode: Boolean = false,
    val userIdInput: String = "",
)

class RegisterViewModel(
    private val repository: PasskeyRepository,
) : ViewModel() {
    companion object {
        fun provideFactory(repository: PasskeyRepository): ViewModelProvider.Factory =
            object : ViewModelProvider.Factory {
                @Suppress("UNCHECKED_CAST")
                override fun <T : ViewModel> create(modelClass: Class<T>): T =
                    RegisterViewModel(repository) as T
            }
    }

    private val _uiState = MutableStateFlow(RegisterUiState())
    val uiState: StateFlow<RegisterUiState> = _uiState.asStateFlow()

    init {
        val savedUserId = repository.getSavedUserId()
        if (savedUserId != null) {
            _uiState.update {
                it.copy(isLoginMode = true, userIdInput = savedUserId)
            }
            Timber.d("Loaded saved user ID: $savedUserId")
        }

        viewModelScope.launch {
            repository.registrationState.collect { regState ->
                if (!regState.isRegistered) {
                    val currentSavedUserId = repository.getSavedUserId()
                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            errorMessage = null,
                            isLoginMode = currentSavedUserId != null,
                            userIdInput = currentSavedUserId ?: "",
                        )
                    }
                }
            }
        }
    }

    fun toggleLoginMode() {
        _uiState.update {
            it.copy(isLoginMode = !it.isLoginMode, errorMessage = null)
        }
    }

    fun updateUserIdInput(text: String) {
        _uiState.update { it.copy(userIdInput = text) }
    }

    fun clearError() {
        _uiState.update { it.copy(errorMessage = null) }
    }

    fun registerPasskey(activityContext: android.app.Activity) {
        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, errorMessage = null) }
            try {
                repository.registerPasskey(activityContext)

                val prfObtained = repository.getPrfOutputOrNull() != null
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage =
                            if (prfObtained) {
                                null
                            } else {
                                UiMessage.StringResource(R.string.error_register_success_no_prf)
                            },
                    )
                }
            } catch (e: CreateCredentialException) {
                Timber.e(e, "Create credential failed")
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage =
                            UiMessage.StringResource(
                                R.string.error_register_credential_failed,
                                listOf(e.type, e.errorMessage ?: ""),
                            ),
                    )
                }
            } catch (e: Exception) {
                Timber.e(e, "Registration failed")
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage =
                            UiMessage.StringResource(
                                R.string.error_register_failed,
                                listOf(e.message ?: ""),
                            ),
                    )
                }
            }
        }
    }

    fun loginWithPasskey(activityContext: android.app.Activity) {
        val userId = _uiState.value.userIdInput.trim()
        if (userId.isBlank()) {
            _uiState.update {
                it.copy(errorMessage = UiMessage.StringResource(R.string.error_user_id_required))
            }
            return
        }

        viewModelScope.launch {
            _uiState.update { it.copy(isLoading = true, errorMessage = null) }
            try {
                val prfOutput = repository.loginWithPasskey(activityContext, userId)
                if (prfOutput == null) {
                    _uiState.update {
                        it.copy(
                            isLoading = false,
                            errorMessage = UiMessage.StringResource(R.string.error_prf_not_supported),
                        )
                    }
                    return@launch
                }

                _uiState.update { it.copy(isLoading = false) }
            } catch (e: GetCredentialException) {
                Timber.e(e, "Login authentication failed")
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage =
                            UiMessage.StringResource(
                                R.string.error_login_credential_failed,
                                listOf(e.type, e.errorMessage ?: ""),
                            ),
                    )
                }
            } catch (e: Exception) {
                Timber.e(e, "Login failed")
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        errorMessage =
                            UiMessage.StringResource(
                                R.string.error_login_failed,
                                listOf(e.message ?: ""),
                            ),
                    )
                }
            }
        }
    }
}
