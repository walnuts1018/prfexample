package dev.walnuts.test.prf_example

import java.security.SecureRandom
import javax.crypto.Cipher
import javax.crypto.spec.GCMParameterSpec
import javax.crypto.spec.SecretKeySpec

object Crypto {
    private const val AES_GCM_ALGORITHM = "AES/GCM/NoPadding"
    private const val AES_KEY_ALGORITHM = "AES"
    private const val GCM_IV_LENGTH_BYTES = 12
    private const val GCM_TAG_LENGTH_BITS = 128

    data class EncryptionResult(
        val iv: ByteArray,
        val ciphertext: ByteArray,
    ) {
        override fun equals(other: Any?): Boolean {
            if (this === other) return true
            if (other !is EncryptionResult) return false
            return iv.contentEquals(other.iv) && ciphertext.contentEquals(other.ciphertext)
        }

        override fun hashCode(): Int {
            var result = iv.contentHashCode()
            result = 31 * result + ciphertext.contentHashCode()
            return result
        }
    }

    fun encrypt(
        prfOutput: ByteArray,
        plaintext: String,
    ): EncryptionResult {
        require(prfOutput.size == 32) { "PRF output must be 32 bytes, got ${prfOutput.size}" }

        val secretKey = SecretKeySpec(prfOutput, AES_KEY_ALGORITHM)

        val iv = ByteArray(GCM_IV_LENGTH_BYTES)
        SecureRandom().nextBytes(iv)

        val cipher = Cipher.getInstance(AES_GCM_ALGORITHM)
        val gcmSpec = GCMParameterSpec(GCM_TAG_LENGTH_BITS, iv)
        cipher.init(Cipher.ENCRYPT_MODE, secretKey, gcmSpec)

        val ciphertext = cipher.doFinal(plaintext.toByteArray(Charsets.UTF_8))

        return EncryptionResult(iv = iv, ciphertext = ciphertext)
    }

    fun decrypt(
        prfOutput: ByteArray,
        iv: ByteArray,
        ciphertext: ByteArray,
    ): String {
        require(prfOutput.size == 32) { "PRF output must be 32 bytes, got ${prfOutput.size}" }
        require(iv.size == GCM_IV_LENGTH_BYTES) { "IV must be $GCM_IV_LENGTH_BYTES bytes, got ${iv.size}" }

        val secretKey = SecretKeySpec(prfOutput, AES_KEY_ALGORITHM)

        val cipher = Cipher.getInstance(AES_GCM_ALGORITHM)
        val gcmSpec = GCMParameterSpec(GCM_TAG_LENGTH_BITS, iv)
        cipher.init(Cipher.DECRYPT_MODE, secretKey, gcmSpec)

        val plaintext = cipher.doFinal(ciphertext)
        return String(plaintext, Charsets.UTF_8)
    }
}