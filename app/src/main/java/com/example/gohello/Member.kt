package com.example.gohello

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass
import java.math.BigDecimal

@JsonClass(generateAdapter = true)
data class Member(
    @field:Json(name = "Id") val id: Int? = null,
    @field:Json(name = "Name") val name: String,
    @field:Json(name = "Age") val age: Int,
    @field:Json(name = "Sex") val sex: Boolean,
    @field:Json(name = "Height") val height: BigDecimal,
)