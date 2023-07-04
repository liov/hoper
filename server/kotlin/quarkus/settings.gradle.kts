rootProject.name = "quarkus"


include("user")
include("protobuf")
project(":protobuf").projectDir = file("../protobuf")

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

