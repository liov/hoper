package xyz.hoper.content.service.impl;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import xyz.hoper.content.dao.ContentRepository;
import xyz.hoper.content.service.ContentService;
import xyz.hoper.content.entity.Content;


@Component
class ContentServiceImpl implements ContentService {

    @Autowired
    private ContentRepository contentRepository ;


    public Content info(Long id ) {
        try {
            return contentRepository.findById(id).get();
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }
}
