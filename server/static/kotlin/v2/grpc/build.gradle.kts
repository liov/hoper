import com.google.protobuf.gradle.generateProtoTasks
import com.google.protobuf.gradle.id
import com.google.protobuf.gradle.ofSourceSet
import com.google.protobuf.gradle.plugins
import com.google.protobuf.gradle.protobuf
import com.google.protobuf.gradle.protoc

plugins{
    id("com.google.protobuf") version "0.8.12"
}

apply(plugin ="kotlin")
apply(plugin ="com.google.protobuf")

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}


dependencies {
    implementation("io.grpc:grpc-kotlin-stub:${rootProject.ext["grpc_kotlin_version"]}")
    implementation("com.google.protobuf:protobuf-java:${rootProject.ext["protobuf_version"]}")
    implementation("com.google.protobuf:protobuf-java-util:3.11.1")
    implementation("io.grpc:grpc-netty-shaded:${rootProject.ext["grpc_version"]}")
    implementation("io.grpc:grpc-protobuf:${rootProject.ext["grpc_version"]}")
    implementation("io.grpc:grpc-stub:${rootProject.ext["grpc_version"]}")
    compileOnly("javax.annotation:javax.annotation-api:1.2")
    implementation("com.google.guava:guava:28.2-jre")
}

protobuf {
    protoc { artifact = "com.google.protobuf:protoc:${rootProject.ext["protobuf_version"]}" }
    plugins {
        // Specify protoc to generate using kotlin protobuf plugin
        id("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:${rootProject.ext["grpc_version"]}"
        }
        // Specify protoc to generate using our grpc kotlin plugin
        id("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:${rootProject.ext["grpc_kotlin_version"]}"
        }
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