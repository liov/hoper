package xyz.hoper.content.entity


import jakarta.persistence.*
import kotlinx.serialization.Serializable



/**
 * @Description TODO
 * @Date 2022/11/21 15:18
 * @Created by lbyi
 */
@Table(name="content",schema="public")
@Serializable
@Entity
class Content {
    @Id
    var id: Long = 0

    var name:String = ""
}