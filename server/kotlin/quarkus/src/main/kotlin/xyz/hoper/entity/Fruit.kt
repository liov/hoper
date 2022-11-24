package xyz.hoper.entity





import io.quarkus.hibernate.reactive.panache.PanacheEntity
import javax.persistence.*


@Entity
@Cacheable
class Fruit: PanacheEntity() {

    @Column(length = 40, unique = true)
    lateinit var name:String
}