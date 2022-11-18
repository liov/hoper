package xyz.hoper

/**
 * @Description TODO
 * @Date 2022/11/18 16:02
 * @Created by lbyi
 */

import org.eclipse.microprofile.graphql.DefaultValue
import org.eclipse.microprofile.graphql.Description
import org.eclipse.microprofile.graphql.GraphQLApi
import org.eclipse.microprofile.graphql.Query


@GraphQLApi
class HelloGraphQLResource {

    @Query
    @Description("Say hello")
    fun sayHello(@DefaultValue("World") name: String): String = "Hello $name"

}