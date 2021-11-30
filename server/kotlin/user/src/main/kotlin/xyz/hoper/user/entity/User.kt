package xyz.hoper.user.entity

import io.vertx.codegen.annotations.DataObject
import io.vertx.core.json.JsonObject
import lombok.Data


@Data
@DataObject(generateConverter = true)
class User {
    private val serialVersionUID = 1L

    var id: Long = 0

    var name: String = ""

    var phone: String = ""

    constructor(jsonObject: JsonObject){
        UserConverter.fromJson(jsonObject, this);
    }

    constructor()

    fun toJson(): JsonObject? {
        val json = JsonObject()
        UserConverter.toJson(this, json)
        return json
    }
}


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
