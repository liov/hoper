plugins {
    kotlin("jvm") version "1.9.0"
    id("org.springframework.boot") version "3.2.4"
    id("io.spring.dependency-management") version "1.1.4"
    //id("org.graalvm.buildtools.native") version "0.9.23" apply false
    id("com.github.johnrengelman.shadow") version "7.0.0" apply false
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
        implementation("org.springframework.boot:spring-boot-starter-data-elasticsearch")
        implementation("org.springframework.boot:spring-boot-starter-data-jpa")
        implementation("org.springframework.boot:spring-boot-starter-data-redis")
        implementation("org.springframework.boot:spring-boot-starter-web")
        developmentOnly("org.springframework.boot:spring-boot-devtools")
        runtimeOnly("com.mysql:mysql-connector-j")
        runtimeOnly("org.postgresql:postgresql")
        annotationProcessor("org.projectlombok:lombok")
        testImplementation("org.springframework.boot:spring-boot-starter-test")
        implementation("io.micrometer:micrometer-tracing-bridge-brave")
        implementation("io.zipkin.reporter2:zipkin-reporter-brave")
        implementation("org.springframework.kafka:spring-kafka")
        runtimeOnly("io.micrometer:micrometer-registry-prometheus")
        implementation("org.reflections:reflections:0.10.2")
        testImplementation("org.springframework.kafka:spring-kafka-test")

    }


    tasks.withType<Test> {
        useJUnitPlatform()
    }

}
