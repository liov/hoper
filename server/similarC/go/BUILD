package(default_visibility = ["//visibility:public"])
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/liov/hoper/go/v2
gazelle(name = "gazelle")
#bazel run //:gazelle -- update-repos -from_file=v2/go.mod

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
        "//v2:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
