package xyz.hoper.dart.bridge

import com.immomo.mls.wrapper.Constant
import com.immomo.mls.wrapper.ConstantClass

/**
 * Created by MLN Templates
 * 注册方法：
 *
 * @see com.immomo.mls.MLSBuilder.registerConstants
 */
@ConstantClass
interface LuaEnum {
    companion object {
        /**
         * Lua可通过 LuaEnum.a 读取
         */
        @Constant
        val a = 1

        /**
         * Lua可通过 LuaEnum.c 读取
         */
        @Constant(alias = "c")
        val b = 2
    }
}