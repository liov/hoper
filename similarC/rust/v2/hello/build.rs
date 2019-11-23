fn main() {
    protoc_rust_grpc::run(protoc_rust_grpc::Args {
        out_dir: "src",
        includes: &["../../../protobuf"],
        input: &[
            "../../../protobuf/helloworld.proto",
        ],
        rust_protobuf: true,
        ..Default::default()
    }).expect("protoc-rust-grpc");
}