package com.example.gohello.network

import com.example.gohello.preference.Preference
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import io.reactivex.rxjava3.schedulers.Schedulers
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.adapter.rxjava3.RxJava3CallAdapterFactory
import retrofit2.converter.moshi.MoshiConverterFactory
import retrofit2.converter.scalars.ScalarsConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
class NetworkModule {
    companion object {
        val client: OkHttpClient =
            OkHttpClient.Builder()
                .connectTimeout(10, TimeUnit.SECONDS)
                .readTimeout(10, TimeUnit.SECONDS)
                .writeTimeout(10, TimeUnit.SECONDS)
                .addInterceptor(HttpLoggingInterceptor().setLevel(HttpLoggingInterceptor.Level.BODY))
                .addInterceptor {
                    val newRequest =
                        it.request().newBuilder().apply {
                            val jwt = Preference.Global.jwt
                            if (jwt?.isNotBlank() == true) {
                                addHeader("Authorization", "Bearer $jwt")
                            }
                        }.build()
                    it.proceed(newRequest)
                }
                .build()
    }

    private fun getRetrofitBuilder() =
        Retrofit.Builder()
            .addCallAdapterFactory(RxJava3CallAdapterFactory.createWithScheduler(Schedulers.io()))
            .baseUrl("http://192.168.10.100:9090")
            .client(client)
            .addConverterFactory(ScalarsConverterFactory.create())
            .addConverterFactory(MoshiConverterFactory.create(MoshiParser.moshi))

    @Provides
    @Singleton
    fun provideRetrofitService(): ApiService = getRetrofitBuilder().build().create(ApiService::class.java)
}
