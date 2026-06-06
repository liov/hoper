cmake_minimum_required(VERSION 3.14)
project(runner LANGUAGES CXX)

# Define the application target. To change its name, change BINARY_NAME in the
# top-level CMakeLists.txt, not the value here, or `flutter run` will no longer
# work.
#
# Any new source files that you add to the application should be added here.
add_executable(${BINARY_NAME} WIN32
  "flutter_window.cpp"
  "main.cpp"
  "utils.cpp"
  "win32_window.cpp"
  "${FLUTTER_MANAGED_DIR}/generated_plugin_registrant.cc"
  "Runner.rc"
  "runner.exe.manifest"
)

# Apply the standard set of build settings. This can be removed for applications
# that need different build settings.
apply_standard_settings(${BINARY_NAME})

# Add preprocessor definitions for the build version.
target_compile_definitions(${BINARY_NAME} PRIVATE "FLUTTER_VERSION=\"${FLUTTER_VERSION}\"")
target_compile_definitions(${BINARY_NAME} PRIVATE "FLUTTER_VERSION_MAJOR=${FLUTTER_VERSION_MAJOR}")
target_compile_definitions(${BINARY_NAME} PRIVATE "FLUTTER_VERSION_MINOR=${FLUTTER_VERSION_MINOR}")
target_compile_definitions(${BINARY_NAME} PRIVATE "FLUTTER_VERSION_PATCH=${FLUTTER_VERSION_PATCH}")
target_compile_definitions(${BINARY_NAME} PRIVATE "FLUTTER_VERSION_BUILD=${FLUTTER_VERSION_BUILD}")

# Disable Windows macros that collide with C++ standard library functions.
target_compile_definitions(${BINARY_NAME} PRIVATE "NOMINMAX")

# Add dependency libraries and include directories. Add any application-specific
# dependencies here.
target_link_libraries(${BINARY_NAME} PRIVATE flutter flutter_wrapper_app)
if(CMAKE_SIZEOF_VOID_P EQUAL 8)
  set(RB_WIN_ARCH "x64")
else()
  set(RB_WIN_ARCH "x86")
endif()
set(RB_LIB "${CMAKE_CURRENT_SOURCE_DIR}/../../staticLibs/windows/${RB_WIN_ARCH}/rb.lib")
if(EXISTS "${RB_LIB}")
  target_link_options(${BINARY_NAME} PRIVATE
    /INCLUDE:rb_ice_viewer_new /INCLUDE:rb_ice_viewer_push /INCLUDE:rb_ice_viewer_poll_out
    /INCLUDE:rb_ice_viewer_state /INCLUDE:rb_ice_viewer_write /INCLUDE:rb_ice_viewer_read
    /INCLUDE:rb_ice_viewer_close /INCLUDE:rb_ice_grpc_open /INCLUDE:rb_ice_grpc_close
    /INCLUDE:rb_ice_grpc_list_files /INCLUDE:rb_ice_grpc_thumb_batch /INCLUDE:rb_ice_grpc_read_file
    /INCLUDE:rb_ice_grpc_read_file_range /INCLUDE:rb_ice_grpc_delete_file)
  target_link_libraries(${BINARY_NAME} PRIVATE "${RB_LIB}")
else()
  message(WARNING "缺少 ${RB_LIB}，远程相册 ICE 不可用")
endif()
target_include_directories(${BINARY_NAME} PRIVATE "${CMAKE_SOURCE_DIR}")

# Run the Flutter tool portions of the build. This must not be removed.
add_dependencies(${BINARY_NAME} flutter_assemble)
