//! 远程相册 Viewer CLI：Rust ICE/直连/中继 + P2P 列表。
use std::env;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let signal_url = env::var("RB_SIGNAL_URL").unwrap_or_else(|_| "ws://127.0.0.1:19090/rb/signal".into());
    let room = env::args().nth(1).unwrap_or_else(|| env::var("RB_ROOM").unwrap_or_default());
    let list_root = env::args().nth(2).unwrap_or_else(|| env::var("RB_LIST_ROOT").unwrap_or_default());
    if room.is_empty() {
        eprintln!("usage: rfv-viewer <room_code> [list_root]");
        std::process::exit(2);
    }
    let ice_ms: u32 = env::var("RB_ICE_TIMEOUT_MS").ok().and_then(|s| s.parse().ok()).unwrap_or(15000);
    if let Err(e) = rfv::transport::run_viewer(signal_url, room, list_root, ice_ms).await {
        eprintln!("rfv-viewer: {e}");
        std::process::exit(1);
    }
}
