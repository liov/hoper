filegroup(
    name = "package-srcs",
    srcs = glob(["**"], exclude=["bazel-*/**", ".git/**", ".idea/**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//v2/gateway:all-srcs",
        "//v2/helloworld:all-srcs",
        "//v2/httptpl:all-srcs",
        "//v2/initialize:all-srcs",
        "//v2/protobuf:all-srcs",
        "//v2/tools:all-srcs",
        "//v2/user:all-srcs",
        "//v2/utils:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
