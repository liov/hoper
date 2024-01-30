import com.google.protobuf.gradle.*
import java.io.ByteArrayOutputStream

plugins {
    val kotlinVersion = "1.9.0"
    kotlin("jvm") version kotlinVersion
    //id("com.squareup.wire") version "3.1.0"
    id("com.google.protobuf") version "0.9.1"
    //kotlin("kapt")
    java
    idea
}

apply<JavaPlugin>()
apply<IdeaPlugin>()

group = "xyz.hoper"
version = "0.0.1-SNAPSHOT"

repositories {
    //maven("https://maven.aliyun.com/repository/public")
    mavenCentral()
    gradlePluginPortal()
    google()
    mavenLocal()
}

java {
    sourceCompatibility = JavaVersion.VERSION_17
    targetCompatibility = JavaVersion.VERSION_17
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

val grpcKotlinVersion:String by project

val protopath: String = file("${rootDir}/../../../proto").absolutePath
val projectpath: String = file("${rootDir}/../../../../tiga").absolutePath


sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
        proto {
            srcDirs(protopath,)
            println(srcDirs)
            includes+="$projectpath/protobuf/_proto"
        }
    }
}

val grpcVersion:String by project
val protocVersion:String by project

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:$protocVersion"
    }

    plugins {
        id("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:$grpcVersion"
        }
        id("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:$grpcKotlinVersion"
        }
//        id("reactor") {
//            artifact = "com.salesforce.servicelibs:reactor-grpc:1.0.0"
//        }
    }

    generateProtoTasks {
        tasks.forEach{task->
            if (task.name == "extractProto") {
                task.actions.add{
                    val extractDir = "$buildDir/extracted-protos/main"
                    delete( "$extractDir/patch")
                }
            }
        }
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
    implementation("com.google.protobuf:protobuf-java:$protocVersion")
    implementation("com.google.protobuf:protobuf-java-util:$protocVersion")
    //implementation("com.google.protobuf:protobuf-kotlin-lite:$protocVersion")
    api("io.grpc:grpc-netty-shaded:$grpcVersion")
    api("io.grpc:grpc-protobuf:$grpcVersion")
    api("io.grpc:grpc-protobuf-lite:$grpcVersion")
    api("io.grpc:grpc-stub:$grpcVersion")
    //api("io.grpc:grpc-kotlin-stub:$grpcKotlinVersion")
    protobuf(files("$projectpath/protobuf/_proto"))
    //api("com.squareup.wire:wire-runtime:${rootProject.ext["wire_version"]}")
    //api("com.squareup.wire:wire-schema-multiplatform:${rootProject.ext["wire_version"]}")
    if (JavaVersion.current().isJava9Compatible) {
        // Workaround for @javax.annotation.Generated
        // see: https://github.com/grpc/grpc-java/issues/3633
        api("org.apache.tomcat:annotations-api:6.0.53")
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
    var outputStr = protolib("github.com/grpc-ecosystem/grpc-gateway/v2")
    includes += outputStr
    stdout.reset()
    outputStr = protolib("github.com/googleapis/googleapis")
    stdout.reset()
    includes += outputStr
    includes += "$projectpath/protobuf/_proto"

    return includes
}

fun protolib(lib:String): String {
    val stdout = ByteArrayOutputStream()
    val args = mutableListOf<String>("go", "list", "-m", "-f", "{{.Dir}}", lib)
    exec {
        workingDir = File(projectpath)
        commandLine(args)
        standardOutput = stdout
    }
    return stdout.toString("utf-8").trim()
}
