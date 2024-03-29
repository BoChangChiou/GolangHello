package com.example.gohello.preference

import kotlin.properties.ReadWriteProperty
import kotlin.reflect.KProperty

object Delegates {
    fun int(defaultValue: Int = 0) =
        object : ReadWriteProperty<Preference, Int> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Int {
                return thisRef.preferences.getInt(property.name, defaultValue)
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Int,
            ) {
                thisRef.preferences.edit().putInt(property.name, value).apply()
            }
        }

    fun long(defaultValue: Long = 0L) =
        object : ReadWriteProperty<Preference, Long> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Long {
                return thisRef.preferences.getLong(property.name, defaultValue)
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Long,
            ) {
                thisRef.preferences.edit().putLong(property.name, value).apply()
            }
        }

    fun boolean(defaultValue: Boolean = false) =
        object : ReadWriteProperty<Preference, Boolean> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Boolean {
                return thisRef.preferences.getBoolean(property.name, defaultValue)
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Boolean,
            ) {
                thisRef.preferences.edit().putBoolean(property.name, value).apply()
            }
        }

    fun float(defaultValue: Float = 0.0f) =
        object : ReadWriteProperty<Preference, Float> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Float {
                return thisRef.preferences.getFloat(property.name, defaultValue)
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Float,
            ) {
                thisRef.preferences.edit().putFloat(property.name, value).apply()
            }
        }

    fun floatOrNull() =
        object : ReadWriteProperty<Preference, Float?> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Float? {
                return if (thisRef.preferences.contains(property.name)) {
                    thisRef.preferences.getFloat(property.name, 0.0f)
                } else {
                    null
                }
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Float?,
            ) {
                if (value == null) {
                    thisRef.preferences.edit().remove(property.name).apply()
                } else {
                    thisRef.preferences.edit().putFloat(property.name, value).apply()
                }
            }
        }

    fun string(defaultValue: String? = null) =
        object : ReadWriteProperty<Preference, String?> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): String? {
                return thisRef.preferences.getString(property.name, defaultValue)
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: String?,
            ) {
                thisRef.preferences.edit().putString(property.name, value).apply()
            }
        }

    fun stringSet(defaultValue: Set<String>? = null) =
        object : ReadWriteProperty<Preference, Set<String>> {
            override fun getValue(
                thisRef: Preference,
                property: KProperty<*>,
            ): Set<String> {
                return thisRef.preferences.getStringSet(property.name, defaultValue)?.toMutableSet() ?: mutableSetOf()
            }

            override fun setValue(
                thisRef: Preference,
                property: KProperty<*>,
                value: Set<String>,
            ) {
                thisRef.preferences.edit().apply {
                    putStringSet(property.name, value).apply()
                }
            }
        }
}
