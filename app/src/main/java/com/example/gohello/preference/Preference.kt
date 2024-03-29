package com.example.gohello.preference

import android.annotation.SuppressLint
import android.content.Context
import android.content.SharedPreferences
import com.example.gohello.MyApplication

abstract class Preference {
    abstract fun getPreferenceName(): String

    val preferences: SharedPreferences by lazy {
        getPreference()
    }

    open fun getPreference(): SharedPreferences =
        MyApplication.instance.getSharedPreferences(
            getPreferenceName(),
            Context.MODE_PRIVATE,
        )

    @SuppressLint("ApplySharedPref")
    fun commit() {
        preferences.edit().commit()
    }

    object Global : Preference() {
        const val TAG = "Global"

        override fun getPreferenceName() = TAG

        var host by Delegates.string("")
        var jwt by Delegates.string("")
    }
}