package xyz.hoper.content.controller;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import xyz.hoper.content.api.ApiResponse;
import xyz.hoper.content.entity.Content;
import xyz.hoper.content.service.ContentService;

/**
 * @Description TODO
 * @Date 2022/11/21 15:36
 * @Created by lbyi
 */
@RestController
@RequestMapping("/api/content")
class ContentController {
    Logger log = LoggerFactory.getLogger(ContentController.class);

    @Autowired
    private ContentService contentService;

    @GetMapping("{id}")
    Content content(@PathVariable Long id)  {
        return contentService.info(id);
    }
}