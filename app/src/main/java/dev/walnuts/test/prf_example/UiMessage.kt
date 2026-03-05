package dev.walnuts.test.prf_example

import androidx.annotation.StringRes

sealed interface UiMessage {
    data class StringResource(
        @param:StringRes val resId: Int,
        val formatArgs: List<Any> = emptyList(),
    ) : UiMessage
}
