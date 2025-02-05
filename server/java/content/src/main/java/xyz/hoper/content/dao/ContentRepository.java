package xyz.hoper.content.dao;


import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;
import xyz.hoper.content.entity.Content;

/**
 * @Description TODO
 * @Date 2022/11/21 15:49
 * @Created by lbyi
 */
@Repository
public interface ContentRepository extends CrudRepository<Content,Long> {
}