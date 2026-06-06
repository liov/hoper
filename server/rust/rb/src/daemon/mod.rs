pub mod relay;
pub mod signal;

use std::net::SocketAddr;
use std::time::Duration;

use axum::{routing::get, Router};

/// 永不返回：HTTP/信令或 relay 异常时退避后整体重启监听。
pub async fn run() {
    let mut backoff = Duration::from_secs(1);
    loop {
        match run_once().await {
            Err(e) => {
                tracing::error!(%e, wait_secs = backoff.as_secs(), "rb-daemon restart");
                tokio::time::sleep(backoff).await;
                backoff = (backoff * 2).min(Duration::from_secs(30));
            }
            Ok(()) => {
                tracing::warn!("rb-daemon stopped unexpectedly, restart");
                tokio::time::sleep(Duration::from_secs(1)).await;
                backoff = Duration::from_secs(1);
            }
        }
    }
}

async fn run_once() -> Result<(), String> {
    let relay_advertise = relay::listen().await.map_err(|e| format!("relay listen: {e}"))?;
    let signal = signal::Hub::new(relay_advertise.clone());
    let app = Router::new()
        .route("/rb/signal", get(signal::ws_handler))
        .route("/rb/health", get(signal::health))
        .with_state(signal);
    let http = std::env::var("RB_HTTP").unwrap_or_else(|_| "0.0.0.0:8080".to_string());
    let addr: SocketAddr = http.parse().map_err(|e| format!("RB_HTTP parse: {e}"))?;
    tracing::info!(%addr, relay = %relay_advertise, "rb-daemon starting");
    let listener = bind_http(addr).await;
    let http_addr = listener.local_addr().map(|a| a.to_string()).unwrap_or_else(|_| "?".into());
    tracing::info!(%http_addr, relay = %relay_advertise, "rb-daemon listening");
    axum::serve(listener, app).await.map_err(|e| format!("http serve: {e}"))
}

async fn bind_http(addr: SocketAddr) -> tokio::net::TcpListener {
    loop {
        match tokio::net::TcpListener::bind(addr).await {
            Ok(ln) => return ln,
            Err(e) => {
                tracing::warn!(%addr, %e, "http bind failed, retry in 1s");
                tokio::time::sleep(Duration::from_secs(1)).await;
            }
        }
    }
}
