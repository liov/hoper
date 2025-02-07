package xyz.hoper.content.service;



import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

import org.springframework.stereotype.Service;
import xyz.hoper.content.api.ApiResponse;
import xyz.hoper.content.dao.ContentRepository;
import xyz.hoper.content.entity.Content;

import java.util.Optional;


@Service
class ContentServiceImpl implements ContentService {

    private static final Logger logger = LoggerFactory.getLogger(ContentServiceImpl.class);

    @Autowired
    private ContentRepository contentRepository ;


    public ApiResponse<Content> info(Long id ) {
        try {
            Optional<Content> contentOptional = contentRepository.findById(id);
            if (contentOptional.isPresent()) {
                logger.info("Content found with id: {}", id);
                return ApiResponse.success(contentOptional.get());
            } else {
                logger.warn("Content not found with id: {}", id);
                return ApiResponse.error(404, "Content not found");
            }
        } catch (Exception e) {
            logger.error("Error retrieving content with id: {}", id, e);
            return ApiResponse.error(500, e.getMessage());
        }
    }
}
