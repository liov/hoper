package xyz.hoper.content.api;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class BusinessException extends RuntimeException {
    private int code;
    private String msg;

    public BusinessException() {
    }

    public BusinessException(ErrorCode errorCode) {
        this(errorCode.getCode(),errorCode.getMsg());
    }

    public BusinessException(int code, String msg) {
        super(msg);
        this.code = code;
        this.msg = msg;
    }
}
