import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.gradle.api.tasks.testing.logging.TestLogEvent.*
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    application
    id("com.github.johnrengelman.shadow") version "7.0.0"
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
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    developmentOnly("org.springframework.boot:spring-boot-devtools")
    testImplementation("org.springframework.boot:spring-boot-starter-test") {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
    }
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
