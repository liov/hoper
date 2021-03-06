import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.springframework.boot") version "2.3.2.RELEASE" apply false
    kotlin("plugin.spring") version "1.4.10" apply false
    id("io.spring.dependency-management") version "1.0.9.RELEASE"
    kotlin("jvm") version "1.4.10"
    kotlin("plugin.jpa") version "1.4.10" apply false
    kotlin("plugin.serialization") version "1.4.10"
}

ext {
    set("vertxVersion", "3.9.1")
    set("junitJupiterEngineVersion", "5.4.0")
    set("grpc_kotlin_version", "0.1.1")
    set("protobuf_version", "3.11.1")
    set("grpc_version", "1.30.2")
    set("springCloudAlibabaVersion", "2.2.0.RELEASE")
    set("wire_version", "3.2.2")
}

allprojects {
    apply<JavaPlugin>()
    group = "xyz.hoper"
    version = "0.0.1-SNAPSHOT"
    java.sourceCompatibility = JavaVersion.VERSION_11

    repositories {
        maven("https://maven.aliyun.com/repository/public")
        mavenCentral()
        gradlePluginPortal()
        google()
        jcenter()
        mavenLocal()
    }
}

subprojects {
    apply(plugin = "io.spring.dependency-management")
    apply(plugin = "org.jetbrains.kotlin.jvm")


    configurations {
        compileOnly {
            extendsFrom(configurations.annotationProcessor.get())
        }
    }

    dependencies {
        implementation("org.apache.logging.log4j:log4j-core:2.12.1")
        implementation("com.fasterxml.jackson.module:jackson-module-kotlin:2.11.1")
        implementation(kotlin("reflect"))
        implementation(kotlin("stdlib"))
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.3.9")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor:1.3.8")
        implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.0.0-RC2")
        implementation("io.projectreactor.kotlin:reactor-kotlin-extensions:1.1.0-RC1")
        testImplementation("io.projectreactor:reactor-test:3.3.5.RELEASE")
    }

    dependencyManagement {
        imports {
            mavenBom("com.alibaba.cloud:spring-cloud-alibaba-dependencies:${property("springCloudAlibabaVersion")}")
        }
        dependencies {
            dependency("org.apache.logging.log4j:log4j-core:2.13.3")
        }
    }

    tasks.withType<Test> {
        useJUnitPlatform()
    }

    tasks.withType<KotlinCompile> {
        kotlinOptions {
            freeCompilerArgs = listOf("-Xjsr305=strict")
            jvmTarget = "11"
        }
    }

}
