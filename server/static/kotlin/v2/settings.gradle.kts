rootProject.name = "v2"
include('../hello')
project(':hello').projectDir = new File(settingsDir, '../hello')