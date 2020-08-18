#指定压缩级别
-optimizationpasses 5

#不跳过非公共的库的类成员
-dontskipnonpubliclibraryclassmembers

#混淆时采用的算法
-optimizations !code/simplification/arithmetic,!field/*,!class/merging/*

#把混淆类中的方法名也混淆了
-useuniqueclassmembernames

#优化时允许访问并修改有修饰符的类和类的成员
-allowaccessmodification

#将文件来源重命名为“SourceFile”字符串
-renamesourcefileattribute SourceFile
#保留行号
-keepattributes SourceFile,LineNumberTable
#保持泛型
-keepattributes Signature

#保持所有实现 Serializable 接口的类成员
-keepclassmembers class * implements java.io.Serializable {
    static final long serialVersionUID;
    private static final java.io.ObjectStreamField[] serialPersistentFields;
    private void writeObject(java.io.ObjectOutputStream);
    private void readObject(java.io.ObjectInputStream);
    java.lang.Object writeReplace();
    java.lang.Object readResolve();
}

#Fragment不需要在AndroidManifest.xml中注册，需要额外保护下
-keep public class * extends androidx.fragment.app.Fragment
-keep class org.chromium.base.**

# 保留所有的本地native方法不被混淆
-keepclasseswithmembernames class * {
    native <methods>;
}

# 保持测试相关的代码
-dontnote junit.framework.**
-dontnote junit.runner.**
-dontwarn android.test.**
-dontwarn androidx.support.test.**
-dontwarn org.junit.**

# 注解和被注解类不混淆
-keepattributes *Annotation*
-keepattributes Exceptions
-keep class com.immomo.mls.annotation.* { *; }

-keep @com.immomo.mls.annotation.LuaClass class * {
    @com.immomo.mls.annotation.LuaBridge <methods>;
}
-keep @com.immomo.mls.wrapper.ConstantClass class * {
    @com.immomo.mls.wrapper.Constant <fields>;
}
-keep,allowobfuscation @interface org.luaj.vm2.utils.LuaApiUsed
-keep @com.immomo.mls.annotation.CreatedByApt class * { *; }
-keep @org.luaj.vm2.utils.LuaApiUsed class *
-keep @org.luaj.vm2.utils.LuaApiUsed class * {
    native <methods>;
    @org.luaj.vm2.utils.LuaApiUsed <methods>;
    @org.luaj.vm2.utils.LuaApiUsed <fields>;
}