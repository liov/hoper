import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.springframework.boot.gradle.tasks.bundling.BootJar
import org.springframework.boot.gradle.tasks.run.BootRun

plugins {
    application
    id("com.github.johnrengelman.shadow") version "5.1.0"
}

ext {
    set("vertxVersion", "3.9.0")
    set("junitJupiterEngineVersion", "5.4.0")
}

var mainClassName = "io.vertx.core.Launcher"
application {
    mainClassName = mainClassName
}

var mainVerticleName = "xyz.hoper.MainVerticle"
var watchForChange = "src/**/*"
var doOnChange = "../gradlew classes"

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

tasks.withType<ShadowJar> {
    manifest {
        attributes (mapOf("Main-Verticle" to mainVerticleName))
    }
    mergeServiceFiles {
        include("META-INF/services/io.vertx.core.spi.VerticleFactory")
    }
}

tasks{
    getByName<BootRun>("bootRun") {
        main = "xyz.hoper.web.ApplicationKt"
        args("run, $mainVerticleName, --redeploy=$watchForChange, --launcher-class=$mainClassName, --on-redeploy=$doOnChange")
    }
}
