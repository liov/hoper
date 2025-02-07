package xyz.hoper.content.api;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.Data;

@JsonInclude(JsonInclude.Include.NON_NULL) // 确保空值不被序列化
@Data
public class ApiResponse<T>  {

    private int code;
    private String msg;
    private T data;

    // 默认构造函数
    public ApiResponse() {
    }

    // 带参数的构造函数
    public ApiResponse(int code, String msg) {
        this.code = code;
        this.msg = msg;
    }

    // 带数据的构造函数
    public ApiResponse(int code, String msg, T data) {
        this.code = code;
        this.msg = msg;
        this.data = data;
    }

    public static <T> ApiResponse<T> success(T data) {
        return new ApiResponse<>(0, "ok", data);
    }

    public static <T> ApiResponse<T> error(int code, String msg) {
        return new ApiResponse<>(code, msg);
    }
}