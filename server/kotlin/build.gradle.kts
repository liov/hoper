import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.springframework.boot.gradle.tasks.bundling.BootBuildImage

plugins {
    val kotlinVersion = "1.7.21"
    id("org.springframework.boot") version "3.0.0"
    id("io.spring.dependency-management") version "1.1.0"
    id("org.graalvm.buildtools.native") version "0.9.17"
    kotlin("jvm") version kotlinVersion
    kotlin("plugin.allopen") version kotlinVersion
    kotlin("plugin.jpa") version kotlinVersion
    kotlin("plugin.serialization") version kotlinVersion
    kotlin("plugin.spring") version kotlinVersion
    java
    idea
    //id("org.springframework.experimental.aot") version "0.11.0-RC1"
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
    apply(plugin = "org.jetbrains.kotlin.jvm")
    apply(plugin = "org.jetbrains.kotlin.plugin.spring")
    apply(plugin = "org.jetbrains.kotlin.plugin.serialization")



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
        implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.4.1")
        implementation(kotlin("reflect"))
        implementation(kotlin("stdlib-jdk8"))
        implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
        implementation("io.projectreactor.kotlin:reactor-kotlin-extensions")
        //implementation("org.springframework.boot:spring-boot-starter-data-jdbc")
        //runtimeOnly("mysql:mysql-connector-java")
        implementation("org.reflections:reflections:0.10.2")
        runtimeOnly("org.postgresql:postgresql")

        //annotationProcessor("org.projectlombok:lombok")
        testImplementation("org.junit.jupiter:junit-jupiter:5.9.1")
        testImplementation("io.projectreactor:reactor-test")

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
            jvmTarget = JavaVersion.VERSION_17.toString()
        }
    }


    tasks.withType<BootBuildImage> {
        //builder = "paketobuildpacks/builder:tiny"
        //environment = mapOf("BP_NATIVE_IMAGE" to "true")
    }

}
