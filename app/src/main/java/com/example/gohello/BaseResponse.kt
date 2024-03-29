package com.example.gohello

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class BaseResponse<T>(
    @field:Json(name = "Status") val status: Int,
    @field:Json(name = "Data") val data: T?,
    @field:Json(name = "Msg") val msg: String?,
)
