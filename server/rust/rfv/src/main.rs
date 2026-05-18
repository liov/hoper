//! 被浏览主机：仅 `rfv` / `rfv <房间码>`，无子命令。浏览目录由客户端 wire 指定。
use std::env;

use ffmpeg_next as ffmpeg;
use tracing_subscriber;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().with_max_level(tracing::Level::DEBUG).init();
    ffmpeg::init().expect("ffmpeg init");
    run_host(room_from_argv()).await;
}

/// 房间码：`RB_ROOM` 或第一个参数 `rfv <room>`（不是子命令）。
fn room_from_argv() -> Option<String> {
    if let Ok(s) = env::var("RB_ROOM") {
        let s = s.trim().to_string();
        if !s.is_empty() {
            return Some(s);
        }
    }
    let arg = env::args().nth(1)?;
    if arg.starts_with('-') {
        return None;
    }
    Some(arg)
}

async fn run_host(room: Option<String>) {
    let listen = env::var("RFV_LISTEN")
        .or_else(|_| env::var("RFV_GRPC_ADDR"))
        .unwrap_or_else(|_| "0.0.0.0:50051".into());
    if let Some(room) = room {
        spawn_agent(room);
    }
    if let Err(e) = rfv::grpc_server::serve(listen).await {
        tracing::error!("rfv serve: {e}");
        std::process::exit(1);
    }
}

#[cfg(feature = "transport")]
fn spawn_agent(room: String) {
    let sandbox = env::var("RB_AGENT_SANDBOX").ok().filter(|s| !s.is_empty());
    let signal_url = env::var("RB_SIGNAL_URL").unwrap_or_else(|_| "ws://127.0.0.1:8080/rb/signal".into());
    let ice_ms: u32 = env::var("RB_ICE_TIMEOUT_MS").ok().and_then(|s| s.parse().ok()).unwrap_or(15000);
    tokio::spawn(async move {
        tracing::info!(%room, "rfv agent: path from viewer wire");
        if let Err(e) = rfv::transport::run_agent(signal_url, room, sandbox, ice_ms).await {
            tracing::error!("rfv agent: {e}");
        }
    });
}

#[cfg(not(feature = "transport"))]
fn spawn_agent(_room: String) {}
