package dev.walnuts.test.prf_example

import android.content.Context
import android.util.Base64
import androidx.core.content.edit
import androidx.credentials.CreatePublicKeyCredentialRequest
import androidx.credentials.CreatePublicKeyCredentialResponse
import androidx.credentials.CredentialManager
import androidx.credentials.GetCredentialRequest
import androidx.credentials.GetPublicKeyCredentialOption
import androidx.credentials.PublicKeyCredential
import dev.walnuts.test.prf_example.api.ApiClient
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive
import timber.log.Timber

data class EncryptedRecord(
    val id: String,
    val dataBase64: String,
    val ivBase64: String,
    val updatedAt: String,
)

data class PasskeyInfo(
    val userId: String,
    val credentialId: String?,
    val prfSupported: Boolean?,
    val registrationResponseJson: String?,
)

data class RegistrationState(
    val isRegistered: Boolean = false,
    val passkeyInfo: PasskeyInfo? = null,
)

data class LastEncryptResult(
    val ciphertextBase64: String,
    val ivBase64: String,
)

class PasskeyRepository(
    private val credentialManager: CredentialManager,
    private val apiClient: ApiClient,
    context: Context,
) {
    private val prefs = context.getSharedPreferences("prf_example_prefs", Context.MODE_PRIVATE)

    companion object {
        private const val KEY_LAST_USER_ID = "last_user_id"
    }

    private val _registrationState = MutableStateFlow(RegistrationState())
    val registrationState: StateFlow<RegistrationState> = _registrationState.asStateFlow()

    private val _encryptedRecords = MutableStateFlow<List<EncryptedRecord>>(emptyList())
    val encryptedRecords: StateFlow<List<EncryptedRecord>> = _encryptedRecords.asStateFlow()

    private val _lastEncryptResult = MutableStateFlow<LastEncryptResult?>(null)
    val lastEncryptResult: StateFlow<LastEncryptResult?> = _lastEncryptResult.asStateFlow()

    private val _prfOutput = MutableStateFlow<ByteArray?>(null)
    val prfOutput: StateFlow<ByteArray?> = _prfOutput.asStateFlow()

    private var registeredUserId: String? = null

    fun getSavedUserId(): String? = prefs.getString(KEY_LAST_USER_ID, null)

    private fun saveUserId(userId: String) {
        prefs.edit { putString(KEY_LAST_USER_ID, userId) }
        Timber.d("Saved user ID to storage: $userId")
    }

    private fun clearSavedUserId() {
        prefs.edit { remove(KEY_LAST_USER_ID) }
        Timber.d("Cleared saved user ID from storage")
    }

    fun getPrfOutputOrNull(): ByteArray? = _prfOutput.value

    suspend fun registerPasskey(activityContext: android.app.Activity): PasskeyInfo {
        val creationJson = apiClient.getRegistrationCreation()
        Timber.d("Server creation options: $creationJson")

        val creationObj = Json.parseToJsonElement(creationJson).jsonObject
        val publicKeyJson =
            creationObj["publicKey"]?.toString()
                ?: throw IllegalStateException("Missing publicKey in server response")

        Timber.d("Registration request JSON (publicKey): $publicKeyJson")

        val createRequest = CreatePublicKeyCredentialRequest(requestJson = publicKeyJson)

        val result =
            credentialManager.createCredential(
                context = activityContext,
                request = createRequest,
            )

        check(result is CreatePublicKeyCredentialResponse) {
            "Unexpected credential response type: ${result::class}"
        }

        val responseJson = result.registrationResponseJson
        Timber.d("Registration response JSON: $responseJson")

        val registrationResult = apiClient.createWebAuthnCredential(responseJson)
        Timber.d("Server registration result: ${registrationResult.responseJson}")

        val jsonObj = Json.parseToJsonElement(responseJson).jsonObject
        val credentialId = jsonObj["id"]?.jsonPrimitive?.content

        val clientExtensionResults = jsonObj["clientExtensionResults"]?.jsonObject
        val prfResults = clientExtensionResults?.get("prf")?.jsonObject
        val prfEnabled =
            prfResults
                ?.get("enabled")
                ?.jsonPrimitive
                ?.content
                ?.toBoolean()

        registeredUserId = registrationResult.userId
        saveUserId(registrationResult.userId)

        val info =
            PasskeyInfo(
                userId = registrationResult.userId,
                credentialId = credentialId,
                prfSupported = prfEnabled,
                registrationResponseJson = responseJson,
            )

        _registrationState.update {
            it.copy(isRegistered = true, passkeyInfo = info)
        }

        val prfFromRegistration = extractPrfOutput(responseJson)
        if (prfFromRegistration != null) {
            _prfOutput.value = prfFromRegistration
            Timber.d(
                "PRF key obtained from registration response (${prfFromRegistration.size} bytes)",
            )
        } else {
            Timber.w("PRF output not found in registration response")
        }

        return info
    }

    suspend fun loginWithPasskey(
        activityContext: android.app.Activity,
        userId: String,
    ): ByteArray? {
        registeredUserId = userId
        saveUserId(userId)

        val prfBytes = performAuthentication(activityContext)
        if (prfBytes != null) {
            _prfOutput.value = prfBytes
            Timber.d("PRF key obtained via login (${prfBytes.size} bytes)")

            val info =
                PasskeyInfo(
                    userId = userId,
                    credentialId = null,
                    prfSupported = true,
                    registrationResponseJson = null,
                )
            _registrationState.update {
                it.copy(isRegistered = true, passkeyInfo = info)
            }
        } else {
            Timber.w("PRF output not available from login authentication")
        }

        return prfBytes
    }

    private suspend fun performAuthentication(activityContext: android.app.Activity): ByteArray? {
        val userId =
            registeredUserId
                ?: throw IllegalStateException("User not registered")

        val assertionJson = apiClient.getVerificationAssertion(userId)
        Timber.d("Server assertion options: $assertionJson")

        val assertionObj = Json.parseToJsonElement(assertionJson).jsonObject
        val publicKeyJson =
            assertionObj["publicKey"]?.toString()
                ?: throw IllegalStateException("Missing publicKey in server response")

        Timber.d("Authentication request JSON (publicKey): $publicKeyJson")

        val getCredentialRequest =
            GetCredentialRequest(
                credentialOptions =
                    listOf(
                        GetPublicKeyCredentialOption(requestJson = publicKeyJson),
                    ),
            )

        val result =
            credentialManager.getCredential(
                context = activityContext,
                request = getCredentialRequest,
            )

        val credential = result.credential
        if (credential !is PublicKeyCredential) return null

        val responseJson = credential.authenticationResponseJson
        Timber.d("Authentication response: $responseJson")

        val verifyResponse = apiClient.verifyWebAuthnAssertion(responseJson)
        Timber.d("Server verification result: $verifyResponse")

        return extractPrfOutput(responseJson)
    }

    suspend fun saveEncryptedData(
        ciphertextBytes: ByteArray,
        ivBytes: ByteArray,
    ): EncryptedRecord {
        val dataBase64 = Base64.encodeToString(ciphertextBytes, Base64.NO_WRAP)
        val ivBase64 = Base64.encodeToString(ivBytes, Base64.NO_WRAP)

        val serverData = apiClient.saveEncryptedData(dataBase64, ivBase64)

        val record =
            EncryptedRecord(
                id = serverData.id,
                dataBase64 = serverData.dataBase64,
                ivBase64 = serverData.ivBase64,
                updatedAt = serverData.updatedAt,
            )

        _encryptedRecords.update { it + record }

        return record
    }

    suspend fun refreshEncryptedRecords() {
        val serverDataList = apiClient.listEncryptedData()
        val records =
            serverDataList.map { serverData ->
                EncryptedRecord(
                    id = serverData.id,
                    dataBase64 = serverData.dataBase64,
                    ivBase64 = serverData.ivBase64,
                    updatedAt = serverData.updatedAt,
                )
            }
        _encryptedRecords.value = records
    }

    fun updateLastEncryptResult(result: LastEncryptResult) {
        _lastEncryptResult.value = result
    }

    fun signOut() {
        _prfOutput.value = null
        _encryptedRecords.value = emptyList()
        _lastEncryptResult.value = null
        _registrationState.value = RegistrationState()
        Timber.d("Signed out: PRF key cleared from memory, saved User ID preserved")
    }

    fun deleteAll() {
        registeredUserId = null
        _prfOutput.value = null
        _encryptedRecords.value = emptyList()
        _lastEncryptResult.value = null
        _registrationState.value = RegistrationState()
        clearSavedUserId()
    }

    private fun extractPrfOutput(responseJson: String): ByteArray? =
        try {
            val jsonObj = Json.parseToJsonElement(responseJson).jsonObject
            val clientExtensionResults = jsonObj["clientExtensionResults"]?.jsonObject
            val prfResults = clientExtensionResults?.get("prf")?.jsonObject
            val results = prfResults?.get("results")?.jsonObject
            val firstBase64 = results?.get("first")?.jsonPrimitive?.content

            if (firstBase64 != null) {
                Base64.decode(firstBase64, Base64.URL_SAFE or Base64.NO_PADDING or Base64.NO_WRAP)
            } else {
                Timber.w("PRF results.first not found in response")
                Timber.d("clientExtensionResults: $clientExtensionResults")
                null
            }
        } catch (e: Exception) {
            Timber.e(e, "Failed to extract PRF output")
            null
        }
}

fun ByteArray.toHexString(): String = joinToString("") { "%02x".format(it) }
