import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.springframework.boot.gradle.tasks.bundling.BootJar
import org.springframework.boot.gradle.tasks.run.BootRun

plugins {
    application
    java
    id("com.github.johnrengelman.shadow") version "5.1.0"
    kotlin("kapt")
}

ext {
    set("vertxVersion", "3.9.0")
    set("junitJupiterEngineVersion", "5.4.0")
}


application {
    mainClassName = "xyz.hoper.web.Application"
}

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}

dependencies {
    implementation("io.vertx:vertx-web-client:${project.project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-auth-jwt:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-grpc:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-service-proxy:${project.ext["vertxVersion"]}:processor")
    implementation("io.vertx:vertx-mysql-client:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web-api-contract:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-auth-oauth2:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-redis-client:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-reactive-streams:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web-graphql:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-rx-java2:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-junit5:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-service-factory:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-pg-client:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-lang-kotlin-coroutines:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-rabbitmq-client:${project.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-lang-kotlin:${project.ext["vertxVersion"]}")
    implementation("org.reflections:reflections:0.9.12")

    compileOnly("io.vertx:vertx-service-proxy:${project.ext["vertxVersion"]}")
    compileOnly("io.vertx:vertx-codegen:${project.ext["vertxVersion"]}")
    annotationProcessor("io.vertx:vertx-service-proxy:${project.ext["vertxVersion"]}")
    kapt("io.vertx:vertx-codegen:${project.ext["vertxVersion"]}:processor")
    testImplementation("io.vertx:vertx-junit5:${project.ext["vertxVersion"]}")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:${project.ext["junitJupiterEngineVersion"]}")
    testImplementation("org.junit.jupiter:junit-jupiter-api:${project.ext["junitJupiterEngineVersion"]}")
}


tasks.withType<Test> {
    useJUnitPlatform()
    testLogging {
        events("PASSED", "FAILED", "SKIPPED")
    }
}

tasks.withType<JavaCompile> {
    options.annotationProcessorGeneratedSourcesDirectory = file("$projectDir/build/generated")
    options.compilerArgs = listOf(
            "-Acodegen.output=src/main"
    )
}