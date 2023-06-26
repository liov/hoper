package xyz.hoper.entity



import jakarta.persistence.*
import kotlinx.serialization.Serializable


@Entity
@Cacheable
@Serializable
open class Fruit {
    @Id
    @GeneratedValue
    var id: Long? = null;

    @Column(length = 40, unique = true)
    open lateinit var name:String
}