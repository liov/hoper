import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootBuildImage

plugins {
    val kotlinVersion = "1.7.10"
    id("org.springframework.boot") version "3.0.0-SNAPSHOT"
    kotlin("plugin.spring") version kotlinVersion
    id("io.spring.dependency-management") version "1.0.13.RELEASE"
    kotlin("jvm") version kotlinVersion
    kotlin("plugin.allopen") version kotlinVersion
    kotlin("plugin.jpa") version kotlinVersion apply false
    kotlin("plugin.serialization") version kotlinVersion
    java
    idea
    //id("org.springframework.experimental.aot") version "0.11.0-RC1"
}

allprojects {
    apply<JavaPlugin>()
    apply<IdeaPlugin>()
    group = "xyz.hoper"
    version = "0.0.1-SNAPSHOT"
    java.sourceCompatibility = JavaVersion.VERSION_17

    repositories {
        maven { url = uri("https://repo.spring.io/milestone") }
        maven { url = uri("https://repo.spring.io/snapshot") }
        //maven("https://maven.aliyun.com/repository/public")
        mavenCentral()
        gradlePluginPortal()
        google()
        mavenLocal()
    }
}

extra["slf4j.version"] = "1.7.20"

configurations.all {
    resolutionStrategy.eachDependency {
        if (requested.group == "org.slf4j") {
            useVersion("1.7.20")
        }
    }
}

subprojects {
    apply(plugin = "io.spring.dependency-management")
    apply(plugin = "org.jetbrains.kotlin.jvm")
    apply(plugin = "org.jetbrains.kotlin.plugin.spring")
    apply(plugin = "org.jetbrains.kotlin.plugin.serialization")
    apply(plugin = "org.springframework.boot")


    configurations {
        compileOnly {
            extendsFrom(configurations.annotationProcessor.get())
        }
    }

    dependencies {
        implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
        implementation(kotlin("reflect"))
        implementation(kotlin("stdlib-jdk8"))
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
        implementation("io.projectreactor.kotlin:reactor-kotlin-extensions")
        implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.4.0")
        implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
        implementation("org.springframework.boot:spring-boot-starter-data-jdbc")
        implementation("org.springframework.boot:spring-boot-starter-data-r2dbc")
        developmentOnly("org.springframework.boot:spring-boot-devtools")
        compileOnly("org.projectlombok:lombok")
        runtimeOnly("mysql:mysql-connector-java")
        //runtimeOnly("org.postgresql:postgresql")
        //runtimeOnly("org.postgresql:r2dbc-postgresql")
        annotationProcessor("org.springframework.boot:spring-boot-configuration-processor")
        annotationProcessor("org.projectlombok:lombok")
        testImplementation("org.springframework.boot:spring-boot-starter-test")
        testImplementation("io.projectreactor:reactor-test")
        implementation("org.junit.jupiter:junit-jupiter:5.9.0")
        implementation("org.reflections:reflections:0.10.2")
    }

    dependencyManagement {
        val springCloudAlibabaVersion: String by project
        imports {
            mavenBom("com.alibaba.cloud:spring-cloud-alibaba-dependencies:$springCloudAlibabaVersion")
            mavenBom(org.springframework.boot.gradle.plugin.SpringBootPlugin.BOM_COORDINATES)
        }
        dependencies {
        }
    }

    tasks.withType<Test> {
        useJUnitPlatform()
    }

    tasks.withType<KotlinCompile> {
        kotlinOptions {
            javaParameters = true
            freeCompilerArgs = listOf("-Xjsr305=strict")
            jvmTarget = "17"
        }
    }


    tasks.withType<BootBuildImage> {
        builder = "paketobuildpacks/builder:tiny"
        environment = mapOf("BP_NATIVE_IMAGE" to "true")
    }

}
