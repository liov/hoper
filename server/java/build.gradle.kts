plugins {
    kotlin("jvm") version "1.9.25"
    kotlin("plugin.spring") version "1.9.25"
    id("org.springframework.boot") version "3.4.2"
    id("io.spring.dependency-management") version "1.1.7"
    //id("org.graalvm.buildtools.native") version "0.9.23" apply false
    id("com.gradleup.shadow") version "8.3.6" apply false
    java
    idea
}

allprojects {
    apply<JavaPlugin>()
    apply<IdeaPlugin>()
    apply(plugin = "kotlin")
    group = "xyz.hoper"
    version = "0.0.1-SNAPSHOT"

    repositories {
        maven { url = uri("https://repo.spring.io/milestone") }
        maven { url = uri("https://repo.spring.io/snapshot") }
        //maven { url = uri("file://${rootProject.projectDir}/protobuf") }
        //maven("https://maven.aliyun.com/repository/public")
        mavenCentral()
        google()
        mavenLocal()
    }

    java {
        toolchain {
            languageVersion = JavaLanguageVersion.of(21)
        }
    }

    kotlin {
        jvmToolchain(21)
        compilerOptions {
            freeCompilerArgs.addAll("-Xjsr305=strict")
        }
    }

}

subprojects {
    if (name == "protobuf" || name == "spring") {
        return@subprojects
    }
    apply(plugin = "io.spring.dependency-management")
    apply(plugin = "org.springframework.boot")

    /*
configurations.all {
    resolutionStrategy.eachDependency {
        if (requested.group == "org.slf4j") {
            useVersion("1.8.0")
        }
    }
}*/

    configurations {
        compileOnly {
            extendsFrom(configurations.annotationProcessor.get())
        }
    }.all{
        exclude(group = "org.springframework.boot",module = "spring-boot-starter-logging")
    }

    sourceSets {
        main {
            java {
                srcDirs("src/main/java")
            }
            kotlin {
                srcDirs("src/main/kotlin")
            }
        }
    }

    dependencies {
        // 添加 Log4j2 依赖
        implementation("org.springframework.boot:spring-boot-starter-log4j2")
        // 排除 Logback 依赖
        implementation("org.jetbrains.kotlin:kotlin-reflect")
        implementation("org.springframework.boot:spring-boot-starter-data-elasticsearch")
        implementation("org.springframework.boot:spring-boot-starter-data-jpa")
        implementation("org.springframework.boot:spring-boot-starter-data-redis")
        implementation("org.springframework.boot:spring-boot-starter-web")
        developmentOnly("org.springframework.boot:spring-boot-devtools")
        //runtimeOnly("com.mysql:mysql-connector-j")
        runtimeOnly("org.postgresql:postgresql")
        annotationProcessor("org.projectlombok:lombok")
        testImplementation("org.springframework.boot:spring-boot-starter-test")
        implementation("io.micrometer:micrometer-tracing-bridge-brave")
        implementation("io.zipkin.reporter2:zipkin-reporter-brave")
        implementation("org.springframework.kafka:spring-kafka")
        runtimeOnly("io.micrometer:micrometer-registry-prometheus")
        implementation("org.reflections:reflections:0.10.2")
        testImplementation("org.jetbrains.kotlin:kotlin-test-junit5")
        testRuntimeOnly("org.junit.platform:junit-platform-launcher")
    }


    tasks.withType<Test> {
        useJUnitPlatform()
    }

}
