val elasticsearchVersion: String by project
val quarkusPlatformGroupId: String by project
val quarkusPlatformArtifactId: String by project
val quarkusPlatformVersion: String by project

plugins {
    val kotlinVersion = "1.8.21"
    kotlin("jvm") version kotlinVersion
    kotlin("plugin.allopen") version kotlinVersion
    kotlin("plugin.jpa") version kotlinVersion
    kotlin("plugin.serialization") version kotlinVersion
    id("io.quarkus")
    java
    idea
}

allprojects {
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
}


configurations.all {
    resolutionStrategy.eachDependency {
        if (requested.group == "jakarta.websocket") {
            useVersion("1.1.2")
        }
        if (requested.group == "jakarta.persistence") {
            useVersion("2.2.3")
        }
    }
}

subprojects {

    apply(plugin = "org.jetbrains.kotlin.jvm")
    apply(plugin = "org.jetbrains.kotlin.plugin.serialization")

    sourceSets {
        main {
            java {
                srcDirs("src/main/java")
            }
        }
    }

    dependencies {
        implementation(enforcedPlatform("${quarkusPlatformGroupId}:${quarkusPlatformArtifactId}:${quarkusPlatformVersion}"))
        //implementation(enforcedPlatform("${quarkusPlatformGroupId}:quarkus-camel-bom:${quarkusPlatformVersion}"))
        //2.x版本 SLF4J: No SLF4J providers were found SLF4J: Class path contains SLF4J bindings targeting slf4j-api versions 1.7.x
        implementation("org.slf4j:slf4j-api:1.7.36")
        implementation("ch.qos.logback:logback-core:1.4.5")
        implementation("ch.qos.logback:logback-classic:1.4.5")
        implementation("io.quarkus:quarkus-hibernate-search-orm-elasticsearch")
        implementation("io.quarkus:quarkus-smallrye-reactive-messaging-kafka")
        implementation("io.quarkus:quarkus-kotlin")
        implementation("io.quarkus:quarkus-rest-client-reactive")
        implementation("io.quarkus:quarkus-rest-client-reactive-kotlin-serialization")
        implementation("io.quarkus:quarkus-resteasy-reactive")
        implementation("io.quarkus:quarkus-resteasy-reactive-kotlin-serialization")
        implementation("io.quarkus:quarkus-hibernate-orm-panache-kotlin")
        implementation("io.quarkus:quarkus-jdbc-postgresql")
        //implementation("io.quarkus:quarkus-hibernate-reactive-panache-kotlin")
        //implementation("io.quarkus:quarkus-reactive-pg-client")
        implementation("io.quarkus:quarkus-config-yaml")
        implementation("io.quarkus:quarkus-smallrye-metrics")
        implementation("io.quarkus:quarkus-smallrye-opentracing")
        implementation("io.quarkus:quarkus-vertx")
        implementation("io.quarkus:quarkus-redis-client")
        implementation("io.quarkus:quarkus-smallrye-openapi")
        implementation("io.quarkus:quarkus-arc")
        testImplementation("io.quarkus:quarkus-junit5")
        testImplementation("io.rest-assured:rest-assured")
    }
}
