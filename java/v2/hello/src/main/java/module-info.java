module hello {
    requires jdk.unsupported;
    requires java.base;
    requires spring.boot;
    requires spring.boot.autoconfigure;
    requires spring.beans;
    requires grpc.api;
    requires lombok;
    requires spring.context;
    requires java.annotation;
    requires grpc.stub;
    requires com.google.common;
    requires grpc.protobuf;
    requires protobuf.java;
}
