package xyz.hoper.content.service

import org.springframework.stereotype.Component
import org.springframework.stereotype.Service
import xyz.hoper.content.entity.Content

/**
 * @Description TODO
 * @Date 2022/11/21 15:36
 * @Created by lbyi
 */

interface ContentService {
  fun info(id:Long): Content
}