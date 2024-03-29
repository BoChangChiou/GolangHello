package com.example.gohello.network

import com.example.gohello.BaseResponse
import com.example.gohello.Member
import com.example.gohello.preference.Preference
import kotlinx.coroutines.flow.flow
import retrofit2.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ApiModel @Inject constructor(private val apiService: ApiService) {

    fun login(pathUrl: String, login: String, password: String) = flow {
        emit(apiService.login(pathUrl, login, password))
    }
    fun get(pathUrl: String) = flow {
        emit(checkJwtFail(apiService.get(pathUrl)))
    }

    fun stop(pathUrl: String) = flow {
        emit(checkJwtFail(apiService.stop(pathUrl)))
    }

    fun insertMember(pathUrl: String, member: Member) = flow {
        emit(checkJwtFail(apiService.insertMember(pathUrl, member)))
    }

    fun patch(url: String, data: Map<String, String>) = flow {
        emit(checkJwtFail(apiService.patch(url, data)))
    }

    fun delete(url: String) = flow {
        emit(checkJwtFail(apiService.delete(url)))
    }

    private fun <T>checkJwtFail(response: Response<BaseResponse<T>>): Response<BaseResponse<T>> {
        if (response.code() == 200 && response.body()?.status == 2) {
            Preference.Global.jwt = ""
            println("clear JWT")
        }
        return response
    }
}