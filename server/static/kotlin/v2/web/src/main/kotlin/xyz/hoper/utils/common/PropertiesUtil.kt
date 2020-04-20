package xyz.hoper.utils.common

import java.io.BufferedInputStream
import java.io.FileInputStream
import java.io.InputStream
import java.util.*

/**
 * 资源文件读取工具
 */
class PropertiesUtil(filePath: String?) {
    private val pps: Properties
    fun getString(key: String?): String {
        return pps.getProperty(key)
    }

    init {
        pps = Properties()
        val `in`: InputStream = BufferedInputStream(FileInputStream(filePath))
        pps.load(`in`)
    }
}