import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
    kotlin("jvm")
    id("org.springframework.boot")
}

tasks.getByName<BootJar>("bootJar") {
    mainClassName = "xyz.hoper.user.ApplicationKt"
}

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}

dependencies {
    api("com.squareup.wire:wire-runtime:3.1.0")
    api("com.squareup.wire:wire-schema-multiplatform:3.1.0")
    api("io.github.lognet:grpc-spring-boot-starter:3.5.2")
    implementation(project(":protobuf"))
}