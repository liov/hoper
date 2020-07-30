import java.lang.System.getenv

import com.google.protobuf.gradle.generateProtoTasks
import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.ofSourceSet
import com.google.protobuf.gradle.plugins
import com.google.protobuf.gradle.protobuf
import com.google.protobuf.gradle.protoc

plugins {
    kotlin("jvm")
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
        java {
            srcDirs("src/main/java")
        }
        proto {
            srcDir("${rootDir}/../../../std_proto")
            println(srcDirs)
        }
    }
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:${rootProject.ext["protobuf_version"]}"
    }

    plugins {
        id("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:${rootProject.ext["grpc_version"]}"
        }
        id("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:${rootProject.ext["grpc_kotlin_version"]}"
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
                        id("grpckt")
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
    implementation("io.grpc:grpc-kotlin-stub:${rootProject.ext["grpc_kotlin_version"]}")
    implementation("com.google.protobuf:protobuf-java:${rootProject.ext["protobuf_version"]}")
    implementation("com.google.protobuf:protobuf-java-util:3.11.1")
    api("io.grpc:grpc-netty-shaded:${rootProject.ext["grpc_version"]}")
    api("io.grpc:grpc-protobuf:${rootProject.ext["grpc_version"]}")
    api("io.grpc:grpc-stub:${rootProject.ext["grpc_version"]}")
    implementation("com.google.guava:guava:28.2-jre")
    //api("com.squareup.wire:wire-runtime:${rootProject.ext["wire_version"]}")
    //api("com.squareup.wire:wire-schema-multiplatform:${rootProject.ext["wire_version"]}")
    if (JavaVersion.current().isJava9Compatible) {
        // Workaround for @javax.annotation.Generated
        // see: https://github.com/grpc/grpc-java/issues/3633
        implementation("org.apache.tomcat:annotations-api:6.0.53")
    }
   // protobuf(files("${rootDir}../../../proto/"))
}