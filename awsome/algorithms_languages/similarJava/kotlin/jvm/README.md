module坑是真的多
inputs.property("moduleName", moduleName) 可以用module-info.java

IDEA 跑kotlin argfile 废的

手动编辑argfile，结果报找不到主类，calss文件放java目录下可以了

然后测试kotlin的方法就是在kotlin写函数然后java调?
```kotlin
compileJava{
        inputs.property("moduleName", moduleName)
        options.encoding = "UTF-8"
        options.compilerArgs = listOf(
          "-Xlint:deprecation",
          "--add-opens=java.base/jdk.internal.misc=jvm",
          "--add-exports=java.base/jdk.internal.misc=jvm",
          "--module-path", classpath.asPath,
          "--patch-module", "$moduleName=${sourceSets["main"].output.asPath}"
        )
        classpath = files()
    }
```
