cmake_minimum_required(VERSION 3.14)
project(c C)

set(CMAKE_C_STANDARD 11)

foreach(_target
        union pointer)
    add_executable(${_target} "${_target}.c"
            ${hw_proto_srcs}
            ${hw_grpc_srcs})
    target_link_libraries(${_target}
            ${_GRPC_GRPCPP_UNSECURE}
            ${_PROTOBUF_LIBPROTOBUF})
endforeach()