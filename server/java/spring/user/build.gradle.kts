import org.springframework.boot.gradle.tasks.bundling.BootBuildImage

plugins {
    application
    id("org.graalvm.buildtools.native") apply true
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
    runtimeOnly("org.postgresql:postgresql")
    runtimeOnly("org.postgresql:r2dbc-postgresql")
    implementation("org.springframework.boot:spring-boot-starter-data-r2dbc")
    implementation("org.springframework.boot:spring-boot-starter-data-redis-reactive")
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