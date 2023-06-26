package xyz.hoper.user.entity

import kotlinx.serialization.Serializable
import org.springframework.data.annotation.Id
import org.springframework.data.relational.core.mapping.Table

@Table("public.user")
@Serializable
class User {
    @Id
    var id: Long = 0

    var name: String = ""

    var phone: String = ""
}

