val elasticsearchVersion: String by project
val quarkusPlatformGroupId: String by project
val quarkusPlatformArtifactId: String by project
val quarkusPlatformVersion: String by project

plugins {
    application
    java
    kotlin("jvm")
    kotlin("plugin.allopen")
    id("io.quarkus")
}

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}




dependencies {
    implementation(enforcedPlatform("${quarkusPlatformGroupId}:${quarkusPlatformArtifactId}:${quarkusPlatformVersion}"))
    implementation("io.quarkus:quarkus-resteasy-reactive-jackson")
    implementation("io.quarkus:quarkus-mutiny")
    implementation("io.quarkus:quarkus-hibernate-reactive")
    implementation("io.quarkus:quarkus-hibernate-reactive-panache")
    implementation("io.quarkus:quarkus-resteasy-reactive")
    implementation("io.quarkus:quarkus-config-yaml")
    implementation("io.quarkus:quarkus-logging-json")
    implementation("io.quarkus:quarkus-smallrye-context-propagation")
    implementation("io.quarkus:quarkus-rest-client-mutiny")
    implementation("io.quarkus:quarkus-spring-boot-properties")
    implementation("io.quarkus:quarkus-grpc")
    implementation("io.quarkus:quarkus-vertx")
    implementation("io.quarkus:quarkus-reactive-mysql-client")
    implementation("io.quarkus:quarkus-smallrye-reactive-messaging-kafka")
    implementation("io.quarkus:quarkus-vertx-web")
    implementation("io.quarkus:quarkus-kotlin")
    implementation("io.quarkus:quarkus-spring-di")
    implementation("io.quarkus:quarkus-redis-client")
    implementation("io.quarkus:quarkus-reactive-pg-client")
    implementation("io.quarkus:quarkus-reactive-messaging-http")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    implementation("io.quarkus:quarkus-arc")
    testImplementation("io.quarkus:quarkus-junit5")
    testImplementation("io.rest-assured:rest-assured")
}

quarkus {
    setSourceDir("$projectDir/src/main/kotlin")
    setOutputDirectory("$projectDir/build/classes/kotlin/main")
}

allOpen {
    annotation("javax.ws.rs.Path")
    annotation("javax.enterprise.context.ApplicationScoped")
    annotation("io.quarkus.test.junit.QuarkusTest")
}

tasks.withType<org.jetbrains.kotlin.gradle.tasks.KotlinCompile> {
    kotlinOptions.jvmTarget = JavaVersion.VERSION_11.toString()
    kotlinOptions.javaParameters = true
}
