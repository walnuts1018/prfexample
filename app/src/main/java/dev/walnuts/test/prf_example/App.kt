package dev.walnuts.test.prf_example

import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import dev.walnuts.test.prf_example.ui.theme.PRFExampleTheme

@Composable
internal fun App(
    viewModel: AppViewModel,
    modifier: Modifier = Modifier,
) {
    Greeting("test", modifier)
}


@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Hello $name!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    PRFExampleTheme {
        Greeting("Android")
    }
}