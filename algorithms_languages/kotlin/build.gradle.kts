import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.jetbrains.kotlin.config.KotlinCompilerVersion

plugins {
  java
  kotlin("jvm") version "1.4.10"
  kotlin("plugin.serialization") version "1.4.10"
}

group = "xyz.hoper"
version = "1.0-SNAPSHOT"

repositories {
  maven("https://maven.aliyun.com/repository/public")
  mavenCentral()
  gradlePluginPortal()
  google()
  jcenter()
  mavenLocal()
}

var junitJupiterEngineVersion = "5.4.0"

dependencies {
  implementation(kotlin("stdlib", KotlinCompilerVersion.VERSION))
  implementation(kotlin("reflect", KotlinCompilerVersion.VERSION))
  implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.3.9")
  implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor:1.3.8")
  implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.0.0-RC2")
  implementation("org.objenesis:objenesis:3.0.1")
  implementation("org.apache.commons:commons-lang3:3.8.1")
  implementation("io.netty:netty-all:5.0.0.Alpha2")
  testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:$junitJupiterEngineVersion")
  testImplementation("org.junit.jupiter:junit-jupiter-api:$junitJupiterEngineVersion")
  implementation(kotlin("script-runtime"))
}

configure<JavaPluginConvention> {
  sourceCompatibility = JavaVersion.VERSION_11
}

var moduleName = "kotlin"


tasks {
  compileJava {
    options.encoding = "UTF-8"
    options.compilerArgs = listOf(
      "-Xlint:deprecation",
      "--add-opens=java.base/jdk.internal.misc=jvm",
      "--add-exports=java.base/jdk.internal.misc=jvm"
    )

  }

  compileKotlin {
    kotlinOptions.jvmTarget = "11"
    kotlinOptions.freeCompilerArgs = listOf("-Xinline-classes")
  }
  compileTestKotlin {
    kotlinOptions.jvmTarget = "11"
    kotlinOptions.freeCompilerArgs = listOf("-Xinline-classes")
  }
}

