rootProject.name = "spring"

include("user")
include("content")
include("protobuf")
project(":protobuf").projectDir = file("../protobuf")

pluginManagement {
    repositories {
        mavenLocal()
        gradlePluginPortal()
        google()
        mavenCentral()
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

