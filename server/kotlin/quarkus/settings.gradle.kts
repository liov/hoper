rootProject.name = "quarkus"


include("protobuf")
include("user")


pluginManagement {
    val quarkusPluginVersion: String by settings
    val quarkusPluginId: String by settings
    repositories {
        mavenLocal()
        gradlePluginPortal()
        google()
        mavenCentral()
    }
    plugins {
        id(quarkusPluginId) version quarkusPluginVersion
    }
}

