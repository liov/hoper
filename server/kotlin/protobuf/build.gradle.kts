import com.google.protobuf.gradle.generateProtoTasks
import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.ofSourceSet
import com.google.protobuf.gradle.plugins
import com.google.protobuf.gradle.protobuf
import com.google.protobuf.gradle.protoc
import java.io.ByteArrayOutputStream

plugins {
    //id("com.squareup.wire") version "3.1.0"
    id("com.google.protobuf") version "0.8.18"
    java
    idea
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

val protopath = file("${rootDir}/../../proto").absolutePath
val projectpath = file("${rootDir}/../go/lib").absolutePath


sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
        proto {
            srcDirs(protopath)
            println(srcDirs)
            println(includes)
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
                    //id("grpckt")
                }
        }
    }
}

idea {

}
//卧槽
//https://blog.csdn.net/qq_15807167/article/details/89737226
//
tasks.getByName<Jar>("jar") {
    enabled = true
}

dependencies {
    implementation("io.grpc:grpc-kotlin-stub:${rootProject.ext["grpc_kotlin_version"]}")
    implementation("com.google.protobuf:protobuf-java:${rootProject.ext["protobuf_version"]}")
    implementation("com.google.protobuf:protobuf-java-util:${rootProject.ext["protobuf_version"]}")
    api("io.grpc:grpc-netty-shaded:${rootProject.ext["grpc_version"]}")
    api("io.grpc:grpc-protobuf:${rootProject.ext["grpc_version"]}")
    api("io.grpc:grpc-stub:${rootProject.ext["grpc_version"]}")
    implementation("com.google.guava:guava:31.0.1-jre")
    protobuf(files("$projectpath/protobuf").filter { file -> file.name.contains("third")||file.name.endsWith(".gen.proto") })
    protobuf(files("$projectpath/protobuf/third"))
    //protobuf(files(protolib("github.com/grpc-ecosystem/grpc-gateway/v2")))
    //protobuf(files(protolib("google.golang.org/protobuf")))
    //protobuf(files(protolib("github.com/googleapis/googleapis")))
    //api("com.squareup.wire:wire-runtime:${rootProject.ext["wire_version"]}")
    //api("com.squareup.wire:wire-schema-multiplatform:${rootProject.ext["wire_version"]}")
    if (JavaVersion.current().isJava9Compatible) {
        // Workaround for @javax.annotation.Generated
        // see: https://github.com/grpc/grpc-java/issues/3633
        implementation("org.apache.tomcat:annotations-api:6.0.53")
    }
    // protobuf(files("${rootDir}../../../proto/"))
}

/*
task<Exec>("googeapis"){
    print(projectpath)
    description = "获取path"
    workingDir = File(projectpath)
    commandLine = listOf("go", "list","-m","-f","{{.Dir}}","github.com/grpc-ecosystem/grpc-gateway/v2")
    doLast {
        print("进来了")
        val outputStr = standardOutput.toString()
        println(outputStr)
    }
}
*/

fun allProtolib(): List<String> {
    val stdout = ByteArrayOutputStream()
    val includes = mutableListOf(protopath)
    val args = mutableListOf<String>("go", "list", "-m", "-f", "{{.Dir}}", "")
    exec {
        workingDir = File(projectpath)
        args[5] = "github.com/grpc-ecosystem/grpc-gateway/v2"
        commandLine(args)
        standardOutput = stdout
    }
    var outputStr = stdout.toString("utf-8").trim()
    includes += outputStr
    stdout.reset()
    exec {
        workingDir = File(projectpath)
        args[5] = "google.golang.org/protobuf"
        commandLine(args)
        standardOutput = stdout
    }
    outputStr = stdout.toString("utf-8").trim()
    includes += outputStr
    stdout.reset()
    exec {
        workingDir = File(projectpath)
        args[5] = "github.com/googleapis/googleapis"
        commandLine(args)
        standardOutput = stdout
    }
    outputStr = stdout.toString("utf-8").trim()
    stdout.reset()
    includes += outputStr
    includes += "$projectpath/protobuf"
    includes += "$projectpath/protobuf/third"

    return includes
}

fun protolib(lib:String): String {
    val stdout = ByteArrayOutputStream()
    val includes = mutableListOf(protopath)
    val args = mutableListOf<String>("go", "list", "-m", "-f", "{{.Dir}}", lib)
    exec {
        workingDir = File(projectpath)
        commandLine(args)
        standardOutput = stdout
    }
    return stdout.toString("utf-8").trim()
}