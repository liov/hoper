package hoper.xyz.dart.bridge;

import android.content.Context;

import androidx.annotation.NonNull;

import com.immomo.mls.LuaViewManager;
import com.immomo.mls.annotation.LuaClass;
import com.immomo.mls.annotation.LuaBridge;
import com.immomo.mls.wrapper.callback.IVoidCallback;

import org.luaj.vm2.Globals;
import org.luaj.vm2.LuaValue;
import xyz.hoper.dart.PageRouter;

/**
 * Created by MLN Template
 * 注册方法：
 * new MLSBuilder.SIHolder(LuaBridge.LUA_CLASS_NAME, LuaBridge.class)
 */
@LuaClass
public class SILuaBridge {
    /**
     * Lua类型，Lua调用方法（和静态调用相同）：
     * ${LuaClassName}:method()
     */
    public static final String LUA_CLASS_NAME = "LuaBridge";

    /**
     * Lua构造函数，不关心初始化参数
     * @param g 虚拟机
     */
    //public SILuaBridge(@NonNull Globals g) {}
    /**
     * Lua构造函数，不需要虚拟机及上下文环境
     * @param init 初始化参数
     */
    //public SILuaBridge(@NonNull LuaValue[] init) {}
    /**
     * Lua构造函数，不需要虚拟机，不关心初始化参数
     */
    //public SILuaBridge() {}

    protected Globals globals;
    /**
     * Lua构造函数
     *
     * @param g    虚拟机
     * @param init 构造方法中传入的参数
     */
    public SILuaBridge(@NonNull Globals g , @NonNull LuaValue[] init) {
        globals = g;
    }

    /**
     * 直接在属性中增加注解，可让Lua有相关属性
     * eg:
     * LuaBridge:property()      --获取相关值
     * LuaBridge:property(10)    --设置相关值
     */
    @LuaBridge
    public int property;

    /**
     * 获取上下文，一般情况，此上下文为Activity
     * @param globals 虚拟机，可通过构造函数存储
     */
    protected Context getContext(@NonNull Globals globals) {
        LuaViewManager m = (LuaViewManager) globals.getJavaUserdata();
        return m != null ? m.context : null;
    }

    /**
     * Lua可通过对象方法调用此方法
     * eg:
     * LuaBridge:openPage()
     */
    @LuaBridge
    public void openPage(String url) {
        PageRouter.openPageByUrl(getContext(globals),url,null);
    }

    /**
     * 通过[LuaBridge.alias]设置别名，使Lua通过别名调用此方法
     * Lua调用方法：
     * LuaBridge:methodC() --不可使用methodB()调用
     *
     *
     * 参数类型可选择:
     * 1. 基本数据类型，及其数组类型
     * 2. String，及其数组类型
     * 3. Callback [IVoidCallback]
     * [com.immomo.mls.wrapper.callback.IBoolCallback]
     * [com.immomo.mls.wrapper.callback.IIntCallback]
     * [com.immomo.mls.wrapper.callback.IStringCallback]
     * [com.immomo.mls.utils.LVCallback]
     * 4. 任意Lua原始类型
     * 5. 已注册自动转换的类型，如[java.util.Map] [java.util.List]
     *
     *
     * 返回类型可选择:
     * 1. 基本数据类型，及其数组类型
     * 2. String，及其数组类型
     * 3. 任意Lua原始类型
     * 4。已注册自动转换的类型，如[java.util.Map] [java.util.List]
     */
    @LuaBridge(alias = "methodC")
    public String[] methodB(int a, boolean b, String c, IVoidCallback d, LuaValue e) {
        return null;
    }

    /**
     * Lua GC当前对象时调用，可不实现
     */
    //void __onLuaGc() {}
}
