package xyz.hoper.web.entity

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