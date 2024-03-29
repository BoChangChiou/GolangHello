package com.example.gohello

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.example.gohello.preference.Preference
import com.example.gohello.ui.theme.GoHelloTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            GoHelloTheme {
                // A surface container using the 'background' color from the theme
                Surface(modifier = Modifier
                    .fillMaxSize()
                    .padding(12.dp), color = MaterialTheme.colorScheme.background) {
                    MainLayout()
                }
            }
        }
    }
}

@Composable
fun MainLayout() {
    val viewModel: MainViewModel = hiltViewModel()
    viewModel.host = Preference.Global.host ?: ""

    val resultState by viewModel.resultState.collectAsState()
    Column {
        LabelInput("Host", Preference.Global.host ?: "") {
            viewModel.host = it
            Preference.Global.host = it
        }
//        LabelInput(label = "Path", "") {
//            viewModel.path = it
//        }
        LabelInput(label = "Data", "") {
            viewModel.data = it
        }
        Button(
            onClick = {
                viewModel.login()
            }
        ) {
            Text(text = "Login")
        }
        Button(
            onClick = {
                viewModel.get()
            }
        ) {
            Text(text = "Get")
        }
        Button(
            onClick = {
                viewModel.insertMember()
            }
        ) {
            Text(text = "Insert")
        }
        Button(
            onClick = {
                viewModel.updateMember()
            }
        ) {
            Text(text = "Update")
        }
        Button(
            onClick = {
                viewModel.deleteMember()
            }
        ) {
            Text(text = "Delete")
        }
        Button(
            onClick = {
                viewModel.closeServer()
            }
        ) {
            Text(text = "Close Server")
        }
        Text(text = "Result:\n$resultState")
    }
}

@Composable
fun LabelInput(label: String, value: String, modifier: Modifier = Modifier, onChange: (String) -> Unit) {
    val valueState = remember { mutableStateOf(value) }
    Column {
        Text(
            text = "$label:",
            modifier = modifier
        )
        TextField(
            value = valueState.value,
            onValueChange = {
                println("LabelInput Host onValueChange $it")
                valueState.value = it
                onChange.invoke(it)
            }
        )
    }
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    GoHelloTheme {
        LabelInput("Android", "6.0") {
        }
    }
}