import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootBuildImage

plugins {
    id("org.springframework.boot") version "2.6.0"
    kotlin("plugin.spring") version "1.6.0"
    id("io.spring.dependency-management") version "1.0.11.RELEASE"
    kotlin("jvm") version "1.6.0"
    kotlin("plugin.allopen") version "1.6.0"
    id("io.quarkus")
    kotlin("plugin.jpa") version "1.6.0" apply false
    kotlin("plugin.serialization") version "1.6.0"
    //id("org.springframework.experimental.aot") version "0.11.0-RC1"
}

ext {
    set("vertxVersion", "4.2.1")
    set("junitJupiterEngineVersion", "5.4.0")
    set("grpc_kotlin_version", "1.2.0")
    set("protobuf_version", "3.19.1")
    set("grpc_version", "1.42.1")
    set("springCloudAlibabaVersion", "2.2.0.RELEASE")
    set("wire_version", "4.0.0-alpha.20")
}

allprojects {
    apply<JavaPlugin>()
    group = "xyz.hoper"
    version = "0.0.1-SNAPSHOT"
    java.sourceCompatibility = JavaVersion.VERSION_11

    repositories {
        maven { url = uri("https://repo.spring.io/milestone") }
        //maven("https://maven.aliyun.com/repository/public")
        mavenCentral()
        gradlePluginPortal()
        google()
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
        implementation("org.apache.logging.log4j:log4j-core:2.15.0")
        implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
        implementation(kotlin("reflect"))
        implementation(kotlin("stdlib"))
        implementation("org.jetbrains.kotlin:kotlin-reflect")
        implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core")
        implementation("org.jetbrains.kotlinx:kotlinx-serialization-json")
        implementation("io.projectreactor.kotlin:reactor-kotlin-extensions")
        testImplementation("io.projectreactor:reactor-test")
    }

    dependencyManagement {
        imports {
            mavenBom("com.alibaba.cloud:spring-cloud-alibaba-dependencies:${property("springCloudAlibabaVersion")}")
        }
        dependencies {
            dependency("org.jetbrains.kotlinx:kotlinx-coroutines-reactor:1.5.2")
            dependency("org.jetbrains.kotlinx:kotlinx-serialization-json:1.3.0")
            dependency("io.projectreactor.kotlin:reactor-kotlin-extensions:1.1.3")
            dependency("io.projectreactor:reactor-test:3.3.5.RELEASE")
            dependency("org.apache.logging.log4j:log4j-core:2.15.0")
            dependency("com.fasterxml.jackson.module:jackson-module-kotlin:2.13.0")
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

    tasks.withType<BootBuildImage> {
        builder = "paketobuildpacks/builder:tiny"
        environment = mapOf("BP_NATIVE_IMAGE" to "true")
    }

}
