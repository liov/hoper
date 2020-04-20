package xyz.hoper.entity

import io.vertx.core.json.JsonObject

object UserConverter {
    fun fromJson(json: Iterable<Map.Entry<String?, Any>>, obj: User) {
        for ((key, value) in json) {
            when (key) {
                "id" -> if (value is Long) {
                    obj.id = value
                }
                "name" -> if (value is String) {
                    obj.name = value
                }
                "phone" -> if (value is String) {
                    obj.phone = value
                }
            }
        }
    }

    fun toJson(obj: User?, json: JsonObject) {
        if (obj != null) {
            toJson(obj, json.map)
        }
    }

    fun toJson(obj: User, json: MutableMap<String?, Any?>) {
        json["id"] = obj.id
        json["name"] = obj.name
        json["phone"] = obj.phone

    }
}
