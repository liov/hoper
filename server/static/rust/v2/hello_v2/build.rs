extern crate protoc_grpcio;

fn main() {
    let proto_root = "../../../../proto";
    println!("cargo:rerun-if-changed={}", proto_root);
    protoc_grpcio::compile_grpc_protos(
        &["../../../../proto/helloworld.proto"],
        &[proto_root],
        &"src/protobuf",
        None
    ).expect("Failed to compile gRPC definitions!");
}
