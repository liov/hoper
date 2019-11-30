#cc_library(
#    name = "hello-greet",
#    srcs = ["hello-greet.cc"],
#    hdrs = ["hello-greet.h"],
#)

cc_binary(
    name = "client",
    srcs = ["client.cc"],
    #hdrs = [""],
    visibility = ["//visibility:public"],
)

cc_binary(
    name = "server",
    srcs = ["server.cc"],
#    deps = [
#        ":hello-greet",
#        "//lib:hello-time",
#    ],
)


