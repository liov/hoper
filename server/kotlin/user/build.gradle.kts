plugins {
    application
    id("org.springframework.boot")
    kotlin("plugin.spring")
    id("com.github.johnrengelman.shadow") version "5.1.0"
    kotlin("kapt")
}


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

dependencies {
    val vertxVersion = "4.2.1"
    implementation(project(":protobuf"))
    implementation("io.vertx:vertx-web-client:$vertxVersion")
    implementation("io.vertx:vertx-auth-jwt:$vertxVersion")
    implementation("io.vertx:vertx-web:$vertxVersion")
    implementation("io.vertx:vertx-grpc:$vertxVersion")
    implementation("io.vertx:vertx-service-proxy:$vertxVersion")
    implementation("io.vertx:vertx-mysql-client:$vertxVersion")
    implementation("io.vertx:vertx-web-api-contract:$vertxVersion")
    implementation("io.vertx:vertx-auth-oauth2:$vertxVersion")
    implementation("io.vertx:vertx-redis-client:$vertxVersion")
    implementation("io.vertx:vertx-reactive-streams:$vertxVersion")
    implementation("io.vertx:vertx-web-graphql:$vertxVersion")
    implementation("io.vertx:vertx-rx-java2:$vertxVersion")
    implementation("io.vertx:vertx-junit5:$vertxVersion")
    implementation("io.vertx:vertx-service-factory:$vertxVersion")
    implementation("io.vertx:vertx-pg-client:$vertxVersion")
    implementation("io.vertx:vertx-lang-kotlin-coroutines:$vertxVersion")
    implementation("io.vertx:vertx-rabbitmq-client:$vertxVersion")
    implementation("io.vertx:vertx-lang-kotlin:$vertxVersion")
    implementation("org.reflections:reflections:0.10.2")
    implementation("org.springframework.boot:spring-boot-starter")
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    compileOnly("org.projectlombok:lombok")
    developmentOnly("org.springframework.boot:spring-boot-devtools")
    annotationProcessor("org.projectlombok:lombok")
    testImplementation("org.springframework.boot:spring-boot-starter-test") {
        exclude(group = "org.junit.vintage", module = "junit-vintage-engine")
    }
    compileOnly("io.vertx:vertx-service-proxy:$vertxVersion")
    compileOnly("io.vertx:vertx-codegen:$vertxVersion")
    annotationProcessor("io.vertx:vertx-service-proxy:$vertxVersion")
    kapt("io.vertx:vertx-codegen:$vertxVersion:processor")
    testImplementation("io.vertx:vertx-junit5:$vertxVersion")
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