import java.io.FileInputStream
import java.util.Properties
import org.jetbrains.kotlin.gradle.dsl.JvmTarget

plugins {
    id("com.android.application")
    // Built-in Kotlin（Flutter 3.44+）；勿再 apply kotlin-android
    id("dev.flutter.flutter-gradle-plugin")
}

// 读取 key.properties（release）
val keystoreProperties = Properties()
val keystorePropertiesFile = rootProject.file("key.properties")
if (keystorePropertiesFile.exists()) {
    keystoreProperties.load(FileInputStream(keystorePropertiesFile))
}

// 团队共用 debug 签名（可选）：存在 debug-key.properties 时覆盖默认 ~/.android/debug.keystore，
// 避免 Mac/Windows 各自 debug 证书不同导致 Flutter 覆盖安装前 Uninstalling old version...
// 模板见 debug-key.properties.example；生成 keystore 后可将 android/team-debug.keystore 提交仓库。
val debugKeystoreProperties = Properties()
val debugKeystorePropertiesFile = rootProject.file("debug-key.properties")
if (debugKeystorePropertiesFile.exists()) {
    debugKeystoreProperties.load(FileInputStream(debugKeystorePropertiesFile))
}

android {
    namespace = "xyz.hoper.app"
    compileSdk = flutter.compileSdkVersion
    ndkVersion = "28.2.13676358"

    compileOptions {
        isCoreLibraryDesugaringEnabled = true
        sourceCompatibility = JavaVersion.VERSION_17
        targetCompatibility = JavaVersion.VERSION_17
    }

    // Rust cdylib → dynLibs/android/<abi>/librb.so（build_flutter_lib.sh --android）
    sourceSets {
        getByName("main") {
            jniLibs.srcDirs("../../dynLibs/android")
        }
    }
    packaging {
        jniLibs {
            pickFirsts += listOf("lib/**/librb.so")
        }
    }

    defaultConfig {
        multiDexEnabled = true
        // TODO: Specify your own unique Application ID (https://developer.android.com/studio/build/application-id.html).
        applicationId = "xyz.hoper.app"
        // You can update the following values to match your application needs.
        // For more information, see: https://flutter.dev/to/review-gradle-config.
        minSdk = flutter.minSdkVersion
        targetSdk = flutter.targetSdkVersion
        versionCode = flutter.versionCode
        versionName = flutter.versionName

        ndk {
            // 仅 arm64；勿用 +=，且需配合 gradle.properties 的 disable-abi-filtering
            abiFilters.clear()
            abiFilters.add("arm64-v8a")
        }
    }

    signingConfigs {
        if (debugKeystorePropertiesFile.exists()) {
            getByName("debug") {
                keyAlias = debugKeystoreProperties["keyAlias"] as String
                keyPassword = debugKeystoreProperties["keyPassword"] as String
                storeFile = debugKeystoreProperties["storeFile"]?.let { rootProject.file(it) }
                storePassword = debugKeystoreProperties["storePassword"] as String
            }
        }
        if (keystorePropertiesFile.exists()) {
            create("release") {
                keyAlias = keystoreProperties["keyAlias"] as String
                keyPassword = keystoreProperties["keyPassword"] as String
                storeFile = keystoreProperties["storeFile"]?.let { file(it) }
                storePassword = keystoreProperties["storePassword"] as String
                enableV2Signing = true
            }
        }
    }

    buildTypes {
        release {
            signingConfig = if (keystorePropertiesFile.exists()) {
                signingConfigs.getByName("release")
            } else {
                signingConfigs.getByName("debug")
            }
        }
    }

}

kotlin {
    compilerOptions {
        jvmTarget.set(JvmTarget.JVM_17)
    }
}

dependencies {
    coreLibraryDesugaring("com.android.tools:desugar_jdk_libs:2.1.5")
}

flutter {
    source = "../.."
}
