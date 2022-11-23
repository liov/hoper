package xyz.hoper.entity


import javax.persistence.*


@Entity
@Cacheable
class Fruit {
    @get:GeneratedValue
    @get:Id
    var id: Long = 0
    @Column(length = 40, unique = true)
     var name:String =""
}