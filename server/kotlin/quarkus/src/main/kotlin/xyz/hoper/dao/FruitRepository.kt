package xyz.hoper.dao


import io.quarkus.hibernate.reactive.panache.PanacheRepository
import xyz.hoper.entity.Fruit
import javax.enterprise.context.ApplicationScoped

/**
 * @Description TODO
 * @Date 2022/11/23 10:30
 * @Created by lbyi
 */
@ApplicationScoped
class FruitRepository: PanacheRepository<Fruit> {
}