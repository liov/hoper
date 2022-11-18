package xyz.hoper.entity

import io.quarkus.hibernate.orm.panache.kotlin.PanacheEntity
import javax.persistence.Cacheable
import javax.persistence.Column
import javax.persistence.Entity


@Entity
@Cacheable
class Fruit: PanacheEntity() {

    @Column(length = 40, unique = true)
    lateinit var name:String
}