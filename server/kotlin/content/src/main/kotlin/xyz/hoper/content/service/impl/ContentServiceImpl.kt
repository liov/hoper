package xyz.hoper.content.service.impl

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Component
import xyz.hoper.content.dao.ContentRepository
import xyz.hoper.content.entity.Content
import xyz.hoper.content.service.ContentService

/**
 * @Description TODO
 * @Date 2022/11/21 15:37
 * @Created by lbyi
 */
@Component
class ContentServiceImpl: ContentService {

    @Autowired
    lateinit var repository: ContentRepository

    override fun info(id:Long): Content {
        return repository.getReferenceById(id)
    }
}