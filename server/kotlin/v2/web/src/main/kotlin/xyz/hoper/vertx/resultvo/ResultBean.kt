package xyz.hoper.vertx.resultvo

import io.vertx.core.json.Json

class ResultBean {
    var code: String? = null
        private set
    var msg: String? = null
        private set
    var data: Any? = null
        private set

    fun setCode(code: String?): ResultBean {
        this.code = code
        return this
    }

    fun setData(data: Any?): ResultBean {
        this.data = data
        return this
    }

    fun setMsg(msg: String?): ResultBean {
        this.msg = msg
        return this
    }

    fun setResultConstant(resultConstant: ResultConstant): ResultBean {
        code = resultConstant.code
        msg = resultConstant.msg
        return this
    }

    fun setResultConstant(resultConstant: ResultConstant, data: Any?): ResultBean {
        code = resultConstant.code
        msg = resultConstant.msg
        this.data = data
        return this
    }

    override fun toString(): String {
        return Json.encode(this)
    }

    companion object {
        fun build(): ResultBean {
            return ResultBean()
        }
    }
}