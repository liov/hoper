use std::path::PathBuf;

fn main() {
    let root = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let proto = root.join("../../../proto/remotebrowse/signal.proto");
    let include = root.join("../../../proto");
    prost_build::Config::new()
        .compile_protos(&[proto], &[include])
        .expect("compile signal.proto");
}
