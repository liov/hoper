plugins {
    application
    id("org.springframework.boot")
    kotlin("plugin.spring")
    id("com.github.johnrengelman.shadow") version "5.1.0"
    kotlin("kapt")
}


application {
    mainClassName = "xyz.hoper.user.UserApplication"
}

sourceSets {
    main {
        java {
            srcDirs("src/main/java")
        }
    }
}

dependencies {
    implementation(project(":protobuf"))
    implementation("io.vertx:vertx-web-client:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-auth-jwt:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-grpc:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-service-proxy:${rootProject.ext["vertxVersion"]}:processor")
    implementation("io.vertx:vertx-mysql-client:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web-api-contract:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-auth-oauth2:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-redis-client:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-reactive-streams:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-web-graphql:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-rx-java2:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-junit5:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-service-factory:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-pg-client:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-lang-kotlin-coroutines:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-rabbitmq-client:${rootProject.ext["vertxVersion"]}")
    implementation("io.vertx:vertx-lang-kotlin:${rootProject.ext["vertxVersion"]}")
    implementation("org.reflections:reflections:0.9.12")
    implementation("org.springframework.boot:spring-boot-starter")
    compileOnly("org.projectlombok:lombok")
    developmentOnly("org.springframework.boot:spring-boot-devtools")
    annotationProcessor("org.projectlombok:lombok")
    testImplementation("org.springframework.boot:spring-boot-starter-test") {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
    }
    compileOnly("io.vertx:vertx-service-proxy:${rootProject.ext["vertxVersion"]}")
    compileOnly("io.vertx:vertx-codegen:${rootProject.ext["vertxVersion"]}")
    annotationProcessor("io.vertx:vertx-service-proxy:${rootProject.ext["vertxVersion"]}")
    kapt("io.vertx:vertx-codegen:${rootProject.ext["vertxVersion"]}:processor")
    testImplementation("io.vertx:vertx-junit5:${rootProject.ext["vertxVersion"]}")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:${rootProject.ext["junitJupiterEngineVersion"]}")
    testImplementation("org.junit.jupiter:junit-jupiter-api:${rootProject.ext["junitJupiterEngineVersion"]}")
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