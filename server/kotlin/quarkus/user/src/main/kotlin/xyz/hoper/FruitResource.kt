package xyz.hoper


import io.quarkus.panache.common.Sort
import io.smallrye.mutiny.Uni
import io.smallrye.mutiny.coroutines.awaitSuspending
import jakarta.enterprise.context.ApplicationScoped
import jakarta.enterprise.inject.Default
import jakarta.inject.Inject
import jakarta.ws.rs.GET
import jakarta.ws.rs.POST
import jakarta.ws.rs.Path
import jakarta.ws.rs.Produces
import jakarta.ws.rs.core.MediaType


import xyz.hoper.dao.FruitRepository
import xyz.hoper.entity.Fruit


/**
 * @Description TODO
 * @Date 2023/6/26 14:35
 * @Created by lbyi
 */
@Path("/fruits")
@ApplicationScoped
@Produces(MediaType.APPLICATION_JSON)
class FruitResource {

    @Inject
    @field: Default
    lateinit var fruitRepository: FruitRepository

    @GET
    @Path("/count")
    fun count() = fruitRepository.count()

    @GET
    fun get(): List<Fruit> {
        return fruitRepository.listAll(Sort.by("name"))
    }

    @GET
    @Path("/{id}")
    fun getSingle(id: Long?): Fruit? {
        return id?.let { fruitRepository.findById(it) }
    }

    @POST
    fun create(fruit: Fruit): Fruit {
        fruitRepository.persist(fruit)
        return fruit
    }
}