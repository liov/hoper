import java.lang.System.getenv

plugins {
    id("org.jetbrains.kotlin.jvm")
    id("com.squareup.wire") version "3.1.0"
}

wire {
    sourcePath {
        //srcDir ("${rootDir}\\..\\..\\..\\proto")
        println(srcDirs)
        //val gopath = getenv("GOPATH")
        //include ("${gopath}/src/**")
    }
    kotlin {
        out = "src/main/kotlin"
    }
}
