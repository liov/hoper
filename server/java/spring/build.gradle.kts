plugins {
    val kotlinVersion = "1.9.0"
    id("org.springframework.boot") version "3.1.1"
    id("io.spring.dependency-management") version "1.1.0"
    id("org.graalvm.buildtools.native") version "0.9.23" apply false
    id("com.github.johnrengelman.shadow") version "7.0.0" apply false
    kotlin("jvm") version kotlinVersion
    kotlin("plugin.allopen") version kotlinVersion
    kotlin("plugin.jpa") version kotlinVersion
    kotlin("plugin.serialization") version kotlinVersion
    kotlin("plugin.spring") version kotlinVersion apply false
    java
    idea
}

allprojects {
    apply<JavaPlugin>()
    apply<IdeaPlugin>()
    group = "xyz.hoper"
    version = "0.0.1-SNAPSHOT"

    repositories {
        maven { url = uri("https://repo.spring.io/milestone") }
        maven { url = uri("https://repo.spring.io/snapshot") }
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
}

configurations.all {
    resolutionStrategy.eachDependency {
        if (requested.group == "org.slf4j") {
            useVersion("1.8.0")
        }
    }
}

subprojects {
    apply(plugin = "io.spring.dependency-management")
    apply(plugin = "org.springframework.boot")
    apply(plugin = "org.jetbrains.kotlin.jvm")
    apply(plugin = "org.jetbrains.kotlin.plugin.spring")
    apply(plugin = "org.jetbrains.kotlin.plugin.serialization")
    apply(plugin = "org.jetbrains.kotlin.plugin.jpa")



    configurations {

        compileOnly {
            extendsFrom(configurations.annotationProcessor.get())
        }
    }

    sourceSets {
        main {
            java {
                srcDirs("src/main/java")
            }
        }
    }

    dependencies {
        implementation("org.slf4j:slf4j-api:2.0.4")
        implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.5.1")
        implementation(kotlin("reflect"))
        implementation(kotlin("stdlib"))
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
        implementation("io.projectreactor.kotlin:reactor-kotlin-extensions")
        implementation("com.fasterxml.jackson.module:jackson-module-kotlin")

        implementation("org.springframework.boot:spring-boot-starter-actuator")
        implementation("org.springframework.boot:spring-boot-starter-graphql")
        implementation("org.springframework.boot:spring-boot-starter-webflux")
        implementation("io.micrometer:micrometer-tracing-bridge-brave")
        implementation("io.zipkin.reporter2:zipkin-reporter-brave")
        implementation("org.springframework.kafka:spring-kafka")
        developmentOnly("org.springframework.boot:spring-boot-devtools")
        runtimeOnly("io.micrometer:micrometer-registry-prometheus")
        implementation("org.reflections:reflections:0.10.2")
        //annotationProcessor("org.projectlombok:lombok")
        annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")
        testImplementation("org.springframework.boot:spring-boot-starter-test")
        developmentOnly("org.jetbrains.kotlinx:kotlinx-coroutines-debug")
        testImplementation("io.projectreactor:reactor-test")
        testImplementation("org.springframework.graphql:spring-graphql-test")
        testImplementation("org.springframework.kafka:spring-kafka-test")

    }


    tasks.withType<Test> {
        useJUnitPlatform()
    }

    tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile> {
        kotlinOptions {
            javaParameters = true
            freeCompilerArgs = listOf("-Xjsr305=strict")
            jvmTarget = JavaVersion.VERSION_17.toString()
        }
    }


}
