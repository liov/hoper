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
    implementation("io.quarkus:quarkus-kubernetes")
    implementation("io.quarkus:quarkus-smallrye-context-propagation")
    implementation("io.quarkus:quarkus-resteasy-reactive-jsonb")
    implementation("io.quarkus:quarkus-smallrye-opentracing")
    implementation("io.quarkus:quarkus-micrometer-registry-prometheus")
    implementation("io.quarkus:quarkus-opentelemetry-exporter-jaeger")
    implementation("io.quarkus:quarkus-mailer")
    implementation("io.quarkus:quarkus-kotlin")
    implementation("io.quarkus:quarkus-vertx-web")
    implementation("io.quarkus:quarkus-smallrye-openapi")
    implementation("io.quarkus:quarkus-vertx")
    implementation("io.quarkus:quarkus-grpc")
    implementation("io.quarkus:quarkus-smallrye-reactive-streams-operators")
    implementation("io.quarkus:quarkus-smallrye-metrics")
    implementation("io.quarkus:quarkus-smallrye-context-propagation")
    implementation("io.quarkus:quarkus-smallrye-graphql")
    implementation("io.quarkus:quarkus-logging-json")
    implementation("io.quarkus:quarkus-redis-client")
    implementation("io.quarkus:quarkus-reactive-pg-client")
    implementation("io.quarkus:quarkus-smallrye-jwt-build")
    implementation("io.quarkus:quarkus-reactive-messaging-http")
    implementation("io.quarkus:quarkus-arc")
    implementation(enforcedPlatform("${quarkusPlatformGroupId}:${quarkusPlatformArtifactId}:${quarkusPlatformVersion}"))
    implementation("io.quarkus:quarkus-resteasy")

    testImplementation("io.quarkus:quarkus-junit5")
    testImplementation("io.rest-assured:kotlin-extensions")
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

tasks.withType<Test> {
    useJUnitPlatform()
    systemProperty("java.util.logging.manager", "org.jboss.logmanager.LogManager")
}