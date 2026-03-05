package dev.walnuts.test.prf_example.api

import io.ktor.client.HttpClient
import io.ktor.client.engine.okhttp.OkHttp
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.plugins.defaultRequest
import io.ktor.client.plugins.logging.LogLevel
import io.ktor.client.plugins.logging.Logger
import io.ktor.client.plugins.logging.Logging
import io.ktor.client.request.get
import io.ktor.client.request.header
import io.ktor.client.request.post
import io.ktor.client.request.setBody
import io.ktor.client.statement.HttpResponse
import io.ktor.client.statement.bodyAsText
import io.ktor.http.ContentType
import io.ktor.http.contentType
import io.ktor.serialization.kotlinx.json.json
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.jsonObject
import timber.log.Timber

class ApiClient(
    baseUrl: String,
) {
    private val json = Json { ignoreUnknownKeys = true }

    private val client =
        HttpClient(OkHttp) {
            install(ContentNegotiation) {
                json(json)
            }
            install(Logging) {
                logger =
                    object : Logger {
                        override fun log(message: String) {
                            Timber.tag("Ktor").d(message)
                        }
                    }
                level = LogLevel.BODY
            }
            defaultRequest {
                url(baseUrl)
            }
            engine {
                config {
                    connectTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                    readTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                    writeTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
                }
            }
        }

    var sessionId: String? = null
        private set

    suspend fun getRegistrationCreation(): String {
        val response =
            client.get("/api/v1/webauthn/registration/creation") {
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("Registration creation failed: ${response.status.value} - $body")
        }
        Timber.d("getRegistrationCreation response: $body")
        return body
    }

    suspend fun createWebAuthnCredential(responseJson: String): RegistrationResult {
        Timber.d("createWebAuthnCredential request JSON: $responseJson")
        val response =
            client.post("/api/v1/webauthn/registration/create") {
                contentType(ContentType.Application.Json)
                setBody(responseJson)
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("Registration create failed: ${response.status.value} - $body")
        }
        Timber.d("createWebAuthnCredential response: $body")
        val jsonObj = Json.parseToJsonElement(body).jsonObject
        val userId =
            jsonObj["user_id"]?.toString()?.trim('"')
                ?: throw ApiException("Missing user_id in response")
        return RegistrationResult(userId = userId, responseJson = body)
    }

    suspend fun getVerificationAssertion(userId: String): String {
        val response =
            client.get("/api/v1/webauthn/verification/assertion") {
                url { parameters.append("user_id", userId) }
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("Verification assertion failed: ${response.status.value} - $body")
        }
        Timber.d("getVerificationAssertion response: $body")
        return body
    }

    suspend fun verifyWebAuthnAssertion(responseJson: String): String {
        val response =
            client.post("/api/v1/webauthn/verification/verify") {
                contentType(ContentType.Application.Json)
                setBody(responseJson)
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("Verification verify failed: ${response.status.value} - $body")
        }
        Timber.d("verifyWebAuthnAssertion response: $body")
        return body
    }

    suspend fun saveEncryptedData(
        dataBase64: String,
        ivBase64: String,
    ): ServerEncryptedData {
        val response =
            client.post("/api/v1/data") {
                contentType(ContentType.Application.Json)
                setBody("""{"data":"$dataBase64","iv":"$ivBase64"}""")
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("Save encrypted data failed: ${response.status.value} - $body")
        }
        Timber.d("saveEncryptedData response: $body")
        return parseEncryptedData(body)
    }

    suspend fun listEncryptedData(): List<ServerEncryptedData> {
        val response =
            client.get("/api/v1/data") {
                applySessionHeader()
            }
        updateSessionId(response)
        val body = response.bodyAsText()
        if (!response.status.isSuccess()) {
            throw ApiException("List encrypted data failed: ${response.status.value} - $body")
        }
        Timber.d("listEncryptedData response: $body")
        return parseEncryptedDataMap(body)
    }

    private fun io.ktor.client.request.HttpRequestBuilder.applySessionHeader() {
        sessionId?.let { header("X-Session-ID", it) }
    }

    private fun updateSessionId(response: HttpResponse) {
        val newSessionId = response.headers["X-Session-ID"]
        if (newSessionId != null) {
            sessionId = newSessionId
            Timber.d("Session ID updated: $newSessionId")
        }
    }

    private fun parseEncryptedData(jsonStr: String): ServerEncryptedData {
        val obj = Json.parseToJsonElement(jsonStr).jsonObject
        return ServerEncryptedData.fromJson(obj)
    }

    private fun parseEncryptedDataMap(jsonStr: String): List<ServerEncryptedData> {
        val obj = Json.parseToJsonElement(jsonStr).jsonObject
        return obj.values.map { ServerEncryptedData.fromJson(it.jsonObject) }
    }
}

private fun io.ktor.http.HttpStatusCode.isSuccess(): Boolean = value in 200..299

data class RegistrationResult(
    val userId: String,
    val responseJson: String,
)

data class ServerEncryptedData(
    val id: String,
    val userId: String,
    val dataBase64: String,
    val ivBase64: String,
    val updatedAt: String,
) {
    companion object {
        fun fromJson(obj: JsonObject): ServerEncryptedData =
            ServerEncryptedData(
                id = obj["ID"]?.toString()?.trim('"') ?: "",
                userId = obj["UserID"]?.toString()?.trim('"') ?: "",
                dataBase64 = obj["Data"]?.toString()?.trim('"') ?: "",
                ivBase64 = obj["IV"]?.toString()?.trim('"') ?: "",
                updatedAt = obj["UpdatedAt"]?.toString()?.trim('"') ?: "",
            )
    }
}

class ApiException(
    message: String,
) : Exception(message)