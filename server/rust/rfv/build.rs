fn main() -> Result<(), Box<dyn std::error::Error>> {
    let proto_root = "../../../proto";
    let hopeio_proto = "../../../thirdparty/protobuf/_proto";
    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(&["remotebrowse/browse.service.proto"], &[proto_root, hopeio_proto])?;
    Ok(())
}
