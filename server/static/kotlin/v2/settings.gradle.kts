rootProject.name = "v2"
include("hello")
project(":hello").projectDir = File(settingsDir, "../hello")