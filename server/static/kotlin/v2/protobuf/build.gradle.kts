import java.lang.System.getenv

import com.google.protobuf.gradle.generateProtoTasks
import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.ofSourceSet
import com.google.protobuf.gradle.plugins
import com.google.protobuf.gradle.protobuf
import com.google.protobuf.gradle.protoc

plugins {
    id("org.jetbrains.kotlin.jvm")
    //id("com.squareup.wire") version "3.1.0"
    id("com.google.protobuf") version "0.8.12"
    id("idea")
}


//wire {
//    sourcePath {
//        srcDir ("${rootDir}\\..\\..\\..\\proto")
//        println(srcDirs)
//        //val gopath = getenv("GOPATH")
//        //include ("${gopath}/src/**")
//    }
//    kotlin {
//        out = "src/main/kotlin"
//    }
//}

sourceSets{
    main {
        proto {
            srcDir("${rootDir}/../../../proto")
            println(srcDirs)
        }
    }
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:3.10.0"
    }

    plugins {
        id("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.27.2"
        }
//        id("reactor") {
//            artifact = "com.salesforce.servicelibs:reactor-grpc:1.0.0"
//        }
    }

    generateProtoTasks {
        ofSourceSet("main").forEach { generateProtoTask ->
            generateProtoTask
                    .plugins {
                        id("grpc")
                        //id("reactor")
                    }
        }
    }
}

idea{

}
//卧槽
//https://blog.csdn.net/qq_15807167/article/details/89737226
//
tasks.getByName<Jar>("jar") {
    enabled = true
}

dependencies{
    implementation ("io.grpc:grpc-netty-shaded:1.27.2")
    implementation ("io.grpc:grpc-protobuf:1.27.2")
    implementation ("io.grpc:grpc-stub:1.27.2")
    implementation("io.github.lognet:grpc-spring-boot-starter:3.5.2")
    //api("com.squareup.wire:wire-runtime:3.1.0")
    //api("com.squareup.wire:wire-schema-multiplatform:3.1.0")
    if (JavaVersion.current().isJava9Compatible) {
        // Workaround for @javax.annotation.Generated
        // see: https://github.com/grpc/grpc-java/issues/3633
        implementation("javax.annotation:javax.annotation-api:1.3.1")
    }
   // protobuf(files("${rootDir}../../../proto/"))
}