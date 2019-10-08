fn main() {
    tonic_build::compile_protos("../../../protobuf/helloworld.proto").unwrap();
}