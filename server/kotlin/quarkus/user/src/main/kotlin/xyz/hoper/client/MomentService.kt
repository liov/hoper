package xyz.hoper.client

import jakarta.ws.rs.GET
import jakarta.ws.rs.Path
import jakarta.ws.rs.QueryParam
import jakarta.ws.rs.core.MultivaluedMap
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient
import org.jboss.resteasy.reactive.RestQuery
import xyz.hoper.entity.Moment


/**
 * @Description TODO
 * @Date 2023/6/26 14:13
 * @Created by lbyi
 */

@Path("/moment")
@RegisterRestClient(configKey = "extensions-api")
interface MomentService {
    @GET
    fun getById(@QueryParam("id") id: Int?): Set<Moment?>?

    @GET
    fun getByName(@RestQuery name: String?): Set<Moment?>?

    @GET
    fun getByFilter(@RestQuery filter: Map<String?, String?>?): Set<Moment?>?

    @GET
    fun getByFilters(@RestQuery filters: MultivaluedMap<String?, String?>?): Set<Moment?>?
}