package xyz.hoper.dao





import io.quarkus.hibernate.orm.panache.kotlin.PanacheRepository
import jakarta.enterprise.context.ApplicationScoped
import xyz.hoper.entity.Fruit

/**
 * @Description TODO
 * @Date 2022/11/23 10:30
 * @Created by lbyi
 */
@ApplicationScoped
class FruitRepository: PanacheRepository<Fruit> {
    fun findByName(name: String) = find("name", name).firstResult()
    fun deleteByName(name: String) = delete("name", name)
}