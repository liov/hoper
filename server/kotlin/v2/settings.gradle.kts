rootProject.name = "v2"

include("user")
include("protobuf")
include("web")
include("quarkus")

pluginManagement {
    repositories {
        mavenLocal()
        gradlePluginPortal()
        google()
        mavenCentral()
    }
    plugins {
        val quarkusPluginVersion: String by settings
        id("io.quarkus") version quarkusPluginVersion
    }
    resolutionStrategy {
        eachPlugin {
            val plugin = requested.id.id
            val module = when {
                plugin.startsWith("com.squareup.wire") -> "com.squareup.wire:wire-gradle-plugin:${requested.version}"
                else -> return@eachPlugin
            }
            println("resolutionStrategy for plugin=$plugin : $module")
            useModule(module)
            if (requested.id.id == "com.google.protobuf") {
                useModule("com.google.protobuf:protobuf-gradle-plugin:${requested.version}")
            }
        }
    }
}
