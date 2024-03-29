package com.example.gohello.network

import com.example.gohello.MyApplication
import com.squareup.moshi.FromJson
import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import com.squareup.moshi.ToJson
import com.squareup.moshi.Types
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import java.io.BufferedReader
import java.io.InputStreamReader
import java.math.BigDecimal

object MoshiParser {
    inline fun <reified T> get(): JsonAdapter<List<T>> = moshi.adapter(Types.newParameterizedType(List::class.java, T::class.java), emptySet())

    val moshi: Moshi =
        Moshi.Builder()
            .add(BigDecimalAdapter)
            .add(KotlinJsonAdapterFactory())
            .build()

    fun <T : Any> toJson(data: T): String = moshi.adapter(data.javaClass).toJson(data)

    inline fun <reified T> listToJson(data: List<T>): String {
        val type = Types.newParameterizedType(List::class.java, T::class.java)
        return moshi.adapter<List<T>>(type).toJson(data)
    }

    inline fun <reified T> fromJson(jsonString: String?): T? =
        jsonString.takeIf { !it.isNullOrBlank() }?.let {
            moshi.adapter(T::class.java).fromJson(it)
        }

    inline fun <reified T> fromJsonArray(jsonString: String?): List<T>? =
        jsonString.takeIf { !it.isNullOrBlank() }?.let {
            val type = Types.newParameterizedType(List::class.java, T::class.java)
            moshi.adapter<List<T>>(type).fromJson(it)
        }

    inline fun <reified T> getFromLocalJson(rawID: Int): T? {
        MyApplication.instance.resources.openRawResource(rawID).use { inputSteam ->
            try {
                BufferedReader(InputStreamReader(inputSteam, "UTF8")).use { reader ->
                    val jsonString = reader.use { it.readText() }
                    fromJson<T>(jsonString)?.let {
                        return it
                    }
                }
            } catch (throwable: Throwable) {
                throwable.printStackTrace()
            }
        }

        return null
    }
}

object BigDecimalAdapter {
    @FromJson
    fun fromJson(string: String) = BigDecimal(string)

    @ToJson
    fun toJson(value: BigDecimal) = value.toString()
}