use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let proto_root = root.join("../../../proto");
    let hopeio_proto = root.join("../../../thirdparty/protobuf/_proto");
    let rb_grpc = std::env::var("CARGO_FEATURE_MEDIA").is_ok()
        || std::env::var("CARGO_FEATURE_RB_CORE").is_ok()
        || std::env::var("CARGO_FEATURE_VIEWER_FFI").is_ok();
    let signal = std::env::var("CARGO_FEATURE_CLIENT").is_ok()
        || std::env::var("CARGO_FEATURE_DAEMON").is_ok()
        || std::env::var("CARGO_FEATURE_TRANSPORT").is_ok()
        || std::env::var("CARGO_FEATURE_VIEWER_FFI").is_ok()
        || std::env::var("CARGO_FEATURE_AGENT_SERVE").is_ok();
    if rb_grpc {
        let mut protos = vec!["remotebrowse/browse.service.proto"];
        if signal {
            protos.insert(0, "remotebrowse/signal.proto");
        }
        let includes = [
            proto_root.to_string_lossy().into_owned(),
            hopeio_proto.to_string_lossy().into_owned(),
        ];
        prost_build::Config::new()
            .compile_protos(&protos, &[includes[0].as_str(), includes[1].as_str()])?;
    } else if signal {
        prost_build::Config::new()
            .compile_protos(&["remotebrowse/signal.proto"], &[&proto_root])?;
    }
    Ok(())
}
