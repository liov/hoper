package xyz.hoper.entity


import jakarta.persistence.Entity
import jakarta.persistence.GeneratedValue
import jakarta.persistence.Id
import kotlinx.serialization.Serializable


/**
 * Example JPA entity.
 *
 * To use it, get access to a JPA EntityManager via injection.
 *
 * {@code
 *    @Inject
 *    lateinit var em:EntityManager;
 *
 *     fun doSomething() {
 *         val entity1 = MyKotlinEntity();
 *         entity1.field = "field-1"
 *         em.persist(entity1);
 *
 *         val entities:List<MyKotlinEntity>  = em.createQuery("from MyEntity", MyKotlinEntity::class.java).getResultList()
 *     }
 * }
 */
@Entity
@Serializable
open class Moment {
    @get:GeneratedValue
    @get:Id
    var id: Long = 0
    open var name: String = ""
}