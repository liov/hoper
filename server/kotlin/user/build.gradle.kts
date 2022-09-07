import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.gradle.api.tasks.testing.logging.TestLogEvent.*
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    application
    id("com.github.johnrengelman.shadow") version "7.0.0"
    kotlin("kapt")
}

repositories {
    maven {
        url = uri("https://s01.oss.sonatype.org/content/repositories/snapshots")
        mavenContent {
            snapshotsOnly()
        }
    }
    mavenCentral()
}

val mainVerticleName = "user.UserApplication"
val launcherClassName = "io.vertx.core.Launcher"

val watchForChange = "src/**/*"
val doOnChange = "${projectDir}/gradlew classes"

application {
    mainClass.set("xyz.hoper.user.UserApplication")
}

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}

configurations {
    implementation {
        exclude(group = "org.slf4j", module = "slf4j-log4j12")
    }
}

dependencies {
    val vertxVersion: String by project
    implementation(project(":protobuf"))
    implementation(platform(org.springframework.boot.gradle.plugin.SpringBootPlugin.BOM_COORDINATES))
    implementation(platform("io.vertx:vertx-stack-depchain:$vertxVersion"))
    implementation("io.vertx:vertx-web-client")
    implementation("io.vertx:vertx-service-proxy")
    implementation("io.vertx:vertx-health-check")
    implementation("io.vertx:vertx-web-openapi")
    implementation("io.vertx:vertx-grpc-server")
    implementation("io.vertx:vertx-auth-oauth2")
    implementation("io.vertx:vertx-tcp-eventbus-bridge")
    implementation("io.vertx:vertx-opentracing")
    implementation("io.vertx:vertx-dropwizard-metrics")
    implementation("io.vertx:vertx-reactive-streams")
    implementation("io.vertx:vertx-grpc-client")
    implementation("io.vertx:vertx-jdbc-client")
    implementation("io.vertx:vertx-service-factory")
    implementation("io.vertx:vertx-pg-client")
    implementation("io.vertx:vertx-web-sstore-cookie")
    implementation("io.vertx:vertx-lang-kotlin-coroutines")
    implementation("io.vertx:vertx-web-sstore-redis")
    implementation("io.vertx:vertx-web-validation")
    implementation("io.vertx:vertx-auth-jwt")
    implementation("io.vertx:vertx-web")
    implementation("io.vertx:vertx-zookeeper")
    implementation("io.vertx:vertx-grpc")
    implementation("io.vertx:vertx-mysql-client")
    implementation("io.vertx:vertx-http-service-factory")
    implementation("io.vertx:vertx-micrometer-metrics")
    implementation("io.vertx:vertx-json-schema")
    implementation("io.vertx:vertx-web-api-contract")
    implementation("io.vertx:vertx-uri-template")
    implementation("io.vertx:vertx-rx-java3")
    implementation("io.vertx:vertx-redis-client")
    implementation("io.vertx:vertx-config")
    implementation("io.vertx:vertx-zipkin")
    implementation("io.vertx:vertx-web-graphql")
    implementation("io.vertx:vertx-opentelemetry")
    implementation("io.vertx:vertx-mail-client")
    implementation("io.vertx:vertx-consul-client")
    implementation("io.vertx:vertx-auth-jdbc")
    implementation("io.vertx:vertx-kafka-client")
    implementation("io.vertx:vertx-lang-kotlin")
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    developmentOnly("org.springframework.boot:spring-boot-devtools")
    testImplementation("org.springframework.boot:spring-boot-starter-test") {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
    }
    kapt("io.vertx:vertx-codegen:$vertxVersion:processor")
    testImplementation("io.vertx:vertx-junit5")
}


tasks.withType<Test> {
    useJUnitPlatform()
    testLogging {
        events("PASSED", "FAILED", "SKIPPED")
    }
}

tasks.withType<JavaCompile> {
    options.generatedSourceOutputDirectory.set(file("$projectDir/build/generated"))
    options.compilerArgs = listOf(
        "-Acodegen.output=src/main"
    )
}

tasks.withType<ShadowJar> {
    archiveClassifier.set("fat")
    manifest {
        attributes(mapOf("Main-Verticle" to mainVerticleName))
    }
    mergeServiceFiles()
}

tasks.withType<JavaExec> {
    args = listOf(
        "run",
        mainVerticleName,
        "--redeploy=$watchForChange",
        "--launcher-class=$launcherClassName",
        "--on-redeploy=$doOnChange"
    )
}