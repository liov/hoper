import org.gradle.api.tasks.testing.logging.TestLogEvent.*
import org.springframework.boot.gradle.tasks.bundling.BootBuildImage

plugins {
    application
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


dependencies {
    implementation(project(":protobuf"))
}


tasks.withType<Test> {
    useJUnitPlatform()
    testLogging {
        events("PASSED", "FAILED", "SKIPPED")
    }
}



/*tasks.withType<JavaCompile> {
    options.generatedSourceOutputDirectory.set(file("$projectDir/build/generated"))
    options.compilerArgs = listOf(
        "-Acodegen.output=src/main"
    )
}*/


tasks.withType<BootBuildImage> {
    //builder = "paketobuildpacks/builder:tiny"
    //environment = mapOf("BP_NATIVE_IMAGE" to "true")
}