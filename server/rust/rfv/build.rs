use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let proto_root = root.join("../../../proto");
    let hopeio_proto = root.join("../../../thirdparty/protobuf/_proto");
    let media = std::env::var("CARGO_FEATURE_MEDIA").is_ok();
    let signal = std::env::var("CARGO_FEATURE_CLIENT").is_ok() || std::env::var("CARGO_FEATURE_DAEMON").is_ok();
    if media {
        let mut protos = vec!["remotebrowse/browse.service.proto"];
        if signal {
            protos.insert(0, "remotebrowse/signal.proto");
        }
        tonic_build::configure()
            .build_server(true)
            .build_client(false)
            .compile_protos(&protos, &[&proto_root, &hopeio_proto])?;
    } else if signal {
        prost_build::Config::new()
            .compile_protos(&[proto_root.join("remotebrowse/signal.proto")], &[&proto_root])?;
    }
    Ok(())
}
