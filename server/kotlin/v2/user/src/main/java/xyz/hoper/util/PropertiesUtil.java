package xyz.hoper.util;

import java.io.BufferedInputStream;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

/**
 * 资源文件读取工具
 */
public class PropertiesUtil {

    private final Properties pps;

    public PropertiesUtil(String filePath) throws IOException {
        pps = new Properties();
        InputStream in = new BufferedInputStream(new FileInputStream(filePath));
        pps.load(in);
    }

    public String getString(String key) {
        return pps.getProperty(key);
    }


}
