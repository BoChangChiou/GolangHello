package com.example.gohello.network

import com.example.gohello.BaseResponse
import com.example.gohello.Member
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.DELETE
import retrofit2.http.Field
import retrofit2.http.FieldMap
import retrofit2.http.FormUrlEncoded
import retrofit2.http.GET
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.Url

interface ApiService {

    @FormUrlEncoded
    @POST
    suspend fun login(
        @Url url: String,
        @Field("Account") account: String,
        @Field("Password") password: String,
    ): Response<BaseResponse<String>>

    @GET
    suspend fun get(
        @Url url: String,
    ): Response<BaseResponse<List<Member>>>

    @GET
    suspend fun stop(
        @Url url: String,
    ): Response<BaseResponse<String>>

    @POST
    suspend fun insertMember(
        @Url url: String,
        @Body member: Member,
    ): Response<BaseResponse<String>>

    @FormUrlEncoded
    @PATCH
    suspend fun patch(
        @Url url: String,
        @FieldMap map: Map<String, String>
    ): Response<BaseResponse<String>>

    @DELETE
    suspend fun delete(
        @Url url: String,
    ): Response<BaseResponse<String>>
}