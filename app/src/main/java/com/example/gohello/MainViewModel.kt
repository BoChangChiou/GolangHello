package com.example.gohello

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.gohello.network.ApiModel
import com.example.gohello.preference.Preference
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.launch
import java.lang.StringBuilder
import java.security.KeyFactory
import java.security.spec.X509EncodedKeySpec
import java.util.Base64
import javax.crypto.Cipher
import javax.inject.Inject

@HiltViewModel
class MainViewModel @Inject constructor(
    private val apiModel: ApiModel
) : ViewModel() {
    var host = ""
    var data = ""

    private var pKey =
        "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLaxSRB24P7C4pzeSReaQ8IdIb7os/7d9rJADxw63B89tMWn4B6PNkvlWhHgOPqWY/5IeGRRfWdvrb8qnQuFCoHbMQvpxYhVxni/z1RCvg3CTXuW38CXyDBBPMU1G9cEJ7bR+m8o5X6zh8YC/1siSHKINjcT9ZnK4dNvW7jN1c1wIDAQAB"
    private val publicKey by lazy { KeyFactory.getInstance("RSA").generatePublic(X509EncodedKeySpec(Base64.getDecoder().decode(pKey))) }

    private val _resultState = MutableStateFlow("")
    val resultState: StateFlow<String> = _resultState.asStateFlow()
    val method = "RSA/ECB/PKCS1Padding"

    fun login() = viewModelScope.launch {
        val map = getDataMap()
        val account = map["Account"]
        if (account == null) {
            _resultState.emit("no input account")
            return@launch
        }

        val password = map["Password"]
        if (password == null) {
            _resultState.emit("no input password")
            return@launch
        }

        println("account ${encryptByPublicKey(method, account)}")
        println("password ${encryptByPublicKey(method, password)}")
        apiModel.login("$host/login", encryptByPublicKey(method, account), encryptByPublicKey(method, password))
            .catch {
                _resultState.emit("get exception: $it")
                it.printStackTrace()
            }
            .collect {
                if (it.code() != 200) {
                    _resultState.emit("login response: http fail ${it.code()}")
                    return@collect
                }
                if (it.body()?.status != 1) {
                    _resultState.emit("login response: login fail ${it.body()?.msg}")
                    return@collect
                }

                Preference.Global.jwt = it.body()?.data ?: ""
                _resultState.emit("login success")
            }
    }

    fun get() = viewModelScope.launch {
        val url = StringBuilder("$host/sqlGet")
        val map = getDataMap()
        map.keys.forEach { key ->
            if (!url.contains("?")) {
                url.append("?")
            }
            url.append("$key=${map[key]}")
        }
        apiModel.get(url.toString())
            .catch {
                _resultState.emit("get exception: $it")
                it.printStackTrace()
            }.collect {
                _resultState.emit("get response: ${it.body()}")
            }
    }

    fun insertMember() = viewModelScope.launch {
        val map = getDataMap()

        val name = map["Name"] ?: run {
            _resultState.emit("no name")
            return@launch
        }
        val sex = map["Sex"]?.lowercase()?.toBooleanStrictOrNull() ?: run {
            _resultState.emit("no sex")
            return@launch
        }
        val height = map["Height"]?.toBigDecimalOrNull() ?: run {
            _resultState.emit("no height")
            return@launch
        }
        val age = map["Age"]?.toIntOrNull() ?: run {
            _resultState.emit("no age")
            return@launch
        }

        apiModel.insertMember("$host/sqlInsert", Member(name = name, sex = sex, height = height, age = age))
            .catch {
                _resultState.emit("insert exception: $it")
                it.printStackTrace()
            }
            .collect {
                _resultState.emit("insert response: ${it.body()}")
            }
    }

    fun updateMember() = viewModelScope.launch {
        val map = getDataMap()
        apiModel.patch("$host/sqlUpdate/$data", map)
            .catch {
                _resultState.emit("update exception: $it")
                it.printStackTrace()
            }
            .collect {
                _resultState.emit("patch ${it.body()}")
            }
    }

    fun deleteMember() = viewModelScope.launch {
        apiModel.delete("$host/sqlDelete/$data")
            .catch {
                _resultState.emit("delete exception: $it")
                it.printStackTrace()
            }
            .collect {
                _resultState.emit("delete response: ${it.body()}")
            }
    }

    private fun getDataMap(): Map<String, String> {
        val map = mutableMapOf<String, String>()
        if (data.isNotEmpty()) {
            data.split(",").forEach {
                val arrays = it.split("=")
                map[arrays[0]] = arrays[1]
            }
        }
        return map
    }

    fun closeServer() = viewModelScope.launch {
        apiModel.stop("$host/stop")
            .catch {
                _resultState.emit("closeServer exception: $it")
                it.printStackTrace()
            }
            .collect {
                _resultState.emit("closeServer response: ${it.body()}")
            }
    }

    private fun encryptByPublicKey(transformation: String, input: String): String {
//****非对称加密****
        //创建cipher对象
        val cipher = Cipher.getInstance(transformation)
        //初始化cipher
        cipher.init(Cipher.ENCRYPT_MODE, publicKey)
        //加密
        val encrypt = cipher.doFinal(input.toByteArray())

        return String(Base64.getEncoder().encode(encrypt))
    }
}