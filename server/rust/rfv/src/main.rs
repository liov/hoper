//! 部署在被浏览主机上的 rfv：默认 gRPC/HTTP 媒体服务；`share` 子命令启动 P2P Agent。
mod file;

use std::env;

use axum::{routing::get, Router};
use ffmpeg_next as ffmpeg;
use tower::ServiceBuilder;
use tower_http::trace::TraceLayer;
use tracing_subscriber;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().with_max_level(tracing::Level::DEBUG).init();
    ffmpeg::init().expect("ffmpeg init");
    let args: Vec<String> = env::args().collect();
    if args.len() >= 2 && args[1] == "share" {
        run_share(&args).await;
        return;
    }
    run_media_server().await;
}

#[cfg(feature = "transport")]
async fn run_share(args: &[String]) {
    let room = args.get(2).cloned().or_else(|| env::var("RB_ROOM").ok()).unwrap_or_default();
    let root = args
        .get(3)
        .cloned()
        .or_else(|| env::var("RB_ROOT").ok())
        .unwrap_or_else(|| ".".into());
    if room.is_empty() {
        eprintln!("usage: rfv share <room_code> [root_dir]");
        std::process::exit(2);
    }
    let signal_url = env::var("RB_SIGNAL_URL").unwrap_or_else(|_| "ws://127.0.0.1:8080/rb/signal".into());
    let ice_ms: u32 = env::var("RB_ICE_TIMEOUT_MS").ok().and_then(|s| s.parse().ok()).unwrap_or(15000);
    if let Err(e) = rfv::transport::run_agent(signal_url, room, root, ice_ms).await {
        eprintln!("rfv share: {e}");
        std::process::exit(1);
    }
}

#[cfg(not(feature = "transport"))]
async fn run_share(_args: &[String]) {
    eprintln!("rfv share 需要编译 feature transport（cargo build --features host）");
    std::process::exit(2);
}

async fn run_media_server() {
    let grpc_addr = env::var("RFV_GRPC_ADDR").unwrap_or_else(|_| "0.0.0.0:50051".into());
    tokio::spawn(async move {
        if let Err(e) = rfv::grpc_server::serve(grpc_addr).await {
            tracing::error!("rfv grpc: {e}");
        }
    });
    let app = Router::new()
        .route("/", get(|| async { "rfv media" }))
        .route("/list_files", get(file::list_files_handler))
        .route("/thumbnail/*path", get(file::file_thumbnail_handler))
        .layer(ServiceBuilder::new().layer(TraceLayer::new_for_http()));
    let http = env::var("RFV_HTTP").unwrap_or_else(|_| "0.0.0.0:3000".into());
    let listener = tokio::net::TcpListener::bind(&http).await.expect("http bind");
    tracing::info!(%http, "rfv media listening");
    axum::serve(listener, app.into_make_service()).await.expect("serve");
}
