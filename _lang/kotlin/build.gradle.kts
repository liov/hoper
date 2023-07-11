import org.jetbrains.kotlin.gradle.tasks.KotlinCompile
import org.jetbrains.kotlin.config.KotlinCompilerVersion

plugins {
  val kotlinVersion = "1.9.0"
  kotlin("jvm") version kotlinVersion
  kotlin("plugin.serialization") version kotlinVersion
  java
  idea
}

allprojects {
  apply<JavaPlugin>()
  apply<IdeaPlugin>()
  group = "xyz.hoper"
  version = "0.0.1-SNAPSHOT"
  java.sourceCompatibility = JavaVersion.VERSION_11

  repositories {
    maven("https://maven.aliyun.com/repository/public")
    mavenCentral()
    gradlePluginPortal()
    google()
    mavenLocal()
  }
}




subprojects {
  apply(plugin = "org.jetbrains.kotlin.jvm")

  val junitJupiterEngineVersion = "5.9.2"
  dependencies {
    implementation(kotlin("stdlib", KotlinCompilerVersion.VERSION))
    implementation(kotlin("reflect", KotlinCompilerVersion.VERSION))
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.6.4")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor:1.6.4")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:1.5.0")
    implementation("org.objenesis:objenesis:3.2")
    implementation("org.apache.commons:commons-lang3:3.12.0")
    implementation("io.netty:netty-all:4.1.90.Final")
    testRuntimeOnly("org.junit.jupiter:junit-jupiter-engine:$junitJupiterEngineVersion")
    testImplementation("org.junit.jupiter:junit-jupiter-api:$junitJupiterEngineVersion")
    implementation(kotlin("script-runtime", KotlinCompilerVersion.VERSION))
  }

  configure<JavaPluginExtension> {
    sourceCompatibility = JavaVersion.VERSION_11
  }

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
}
