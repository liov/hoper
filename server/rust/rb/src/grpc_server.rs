//! 本机调试：HTTP/2 + protobuf 与 Range 媒体同端口（axum Router）。
use std::time::Duration;

pub use crate::remotebrowse_svc::{proto, MediaSvc};

/// 永不返回：监听或 serve 失败时退避重试。
pub async fn serve_forever(addr: String) {
    let bind: std::net::SocketAddr = match addr.parse() {
        Ok(a) => a,
        Err(e) => {
            tracing::error!(%addr, %e, "rb http2 addr parse failed, retry in 5s");
            loop {
                tokio::time::sleep(Duration::from_secs(5)).await;
            }
        }
    };
    let mut backoff = Duration::from_secs(1);
    loop {
        match serve_once(bind).await {
            Err(e) => {
                tracing::error!(%bind, %e, wait_secs = backoff.as_secs(), "rb http2 restart");
                tokio::time::sleep(backoff).await;
                backoff = (backoff * 2).min(Duration::from_secs(30));
            }
            Ok(()) => {
                tracing::warn!(%bind, "rb http2 stopped unexpectedly, restart");
                tokio::time::sleep(Duration::from_secs(1)).await;
                backoff = Duration::from_secs(1);
            }
        }
    }
}

/// 同一 `TcpListener` 上提供远程浏览 HTTP/2 API 与 `/rb/v1/media`。
async fn serve_once(bind: std::net::SocketAddr) -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let media = MediaSvc::default();
    let app = crate::remotebrowse::http_api::router(media.clone())
        .merge(crate::remotebrowse::media_http::router(media))
        .merge(crate::file::router())
        .layer(crate::remotebrowse::http_log::layer());
    tracing::info!(%bind, "rb http2 bind");
    let listener = bind_http(bind).await;
    tracing::info!(addr = %listener.local_addr()?, "rb http2 listening");
    axum::serve(listener, app).await?;
    Ok(())
}

async fn bind_http(addr: std::net::SocketAddr) -> tokio::net::TcpListener {
    loop {
        match tokio::net::TcpListener::bind(addr).await {
            Ok(ln) => return ln,
            Err(e) => {
                tracing::warn!(%addr, %e, "rb http2 bind failed, retry in 1s");
                tokio::time::sleep(Duration::from_secs(1)).await;
            }
        }
    }
}
