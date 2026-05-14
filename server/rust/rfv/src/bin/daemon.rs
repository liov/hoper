use std::net::SocketAddr;

use axum::{routing::get, Router};
use tracing_subscriber::EnvFilter;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt()
        .with_env_filter(EnvFilter::from_default_env())
        .init();
    let relay_addr = rfv::daemon::relay::listen().await.expect("relay listen");
    let signal = rfv::daemon::signal::Hub::new(relay_addr);
    let app = Router::new()
        .route("/rb/signal", get(rfv::daemon::signal::ws_handler))
        .route("/rb/health", get(rfv::daemon::signal::health))
        .with_state(signal);
    let http = std::env::var("RB_HTTP").unwrap_or_else(|_| "0.0.0.0:8080".to_string());
    let addr: SocketAddr = http.parse().expect("RB_HTTP");
    tracing::info!(%addr, relay = %relay_addr, "rfv-daemon started");
    let listener = tokio::net::TcpListener::bind(addr).await.expect("http bind");
    axum::serve(listener, app).await.expect("serve");
}
