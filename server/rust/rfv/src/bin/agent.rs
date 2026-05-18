//! 远程相册 Agent：Rust ICE/直连/中继 + 本地 ffmpeg 数据面。
use std::env;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt::init();
    let signal_url = env::var("RB_SIGNAL_URL").unwrap_or_else(|_| "ws://127.0.0.1:19090/rb/signal".into());
    let room = env::args().nth(1).unwrap_or_else(|| env::var("RB_ROOM").unwrap_or_default());
    let root = env::args().nth(2).unwrap_or_else(|| env::var("RB_ROOT").unwrap_or_else(|_| ".".into()));
    if room.is_empty() {
        eprintln!("usage: rfv-agent <room_code> [root_dir]");
        std::process::exit(2);
    }
    let ice_ms: u32 = env::var("RB_ICE_TIMEOUT_MS").ok().and_then(|s| s.parse().ok()).unwrap_or(15000);
    if let Err(e) = rfv::transport::run_agent(signal_url, room, root, ice_ms).await {
        eprintln!("rfv-agent: {e}");
        std::process::exit(1);
    }
}
