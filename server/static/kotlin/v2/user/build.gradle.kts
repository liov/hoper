import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootJar

plugins {
    kotlin("jvm")
    id("org.springframework.boot")
}

tasks.getByName<BootJar>("bootJar") {
    mainClassName = "xyz.hoper.user.ApplicationKt"
}
